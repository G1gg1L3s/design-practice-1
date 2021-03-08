package gomodule

import (
	"fmt"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/roman-mazur/bood/gomodule")

	// Ninja rule to execute go build.
	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	// Ninja rule to execute go mod vendor.
	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	goTest = pctx.StaticRule("gotest", blueprint.RuleParams{
		Command:     "cd $workDir && go test -v $pkg > $outputPath",
		Description: "Build and test $pkg",
	}, "workDir", "pkg", "outputPath")
)

// goBinaryModuleType implements the simplest Go binary build without running tests for the target Go package.
type testedBinaryModule struct {
	blueprint.SimpleName

	properties struct {
		// Go package name to build as a command with "go build".
		Pkg string
		// Go package name to test as a command with "go test".
		TestPkg string
		// List of source files.
		Srcs []string
		// List of test source files.
		TestSrcs []string
		// Exclude patterns.
		SrcsExclude []string
		// Test Exclude patterns.
		TestSrcsExclude []string
		// If to call vendor command.
		VendorFirst bool
		// Example of how to specify dependencies.
		Deps []string
	}
}

func (tb *testedBinaryModule) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return tb.properties.Deps
}

func (tb *testedBinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	// path to main binary
	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	// resolve sources
	inputs, unresolved := resolvePatterns(ctx, tb.properties.Srcs, tb.properties.SrcsExclude)
	// if we have bad patterns, return
	if len(unresolved) != 0 {
		reportUnresolved(ctx, unresolved)
		return
	}

	// resolve test sources
	testInputs, unresolved := resolvePatterns(ctx, tb.properties.TestSrcs, tb.properties.TestSrcsExclude)
	// if we have bad patterns, return
	if len(unresolved) != 0 {
		reportUnresolved(ctx, unresolved)
		return
	}

	if tb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		inputs = append(inputs, vendorDirPath)
	}

	// Generate rule for the main build
	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go binary", name),
		Rule:        goBuild,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        tb.properties.Pkg,
		},
	})

	// Append our main binary to test input, so you ninja will rerun the tests
	// if we change one of the sources
	testInputs = append(testInputs, outputPath)

	// test artifact
	outTestPath := fmt.Sprintf(".%s.test.out", name)
	outTestPath = path.Join(config.BaseOutputDir, outTestPath)

	// Generate our rule for the tests
	// It will produce test artifact which won't run the tests again if nothing changes
	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Testing %s", name),
		Rule:        goTest,
		Outputs:     []string{outTestPath},
		Implicits:   testInputs,
		Args: map[string]string{
			"outputPath": outTestPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        tb.properties.TestPkg,
		},
	})

}

// reportUnresolved iterates over all patterns and prints error message
func reportUnresolved(ctx blueprint.ModuleContext, unresolved []string) {
	for _, pattern := range unresolved {
		ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", pattern)
	}
}

// resolvePatterns returns resolved files and poorly constructed patterns
func resolvePatterns(ctx blueprint.ModuleContext, patterns []string, exclude []string) ([]string, []string) {
	var result = []string{}
	var unresolved = []string{}

	for _, src := range patterns {
		if matches, err := ctx.GlobWithDeps(src, exclude); err == nil {
			result = append(result, matches...)
		} else {
			unresolved = append(unresolved, src)
		}
	}

	return result, unresolved
}

// BinFactory is a factory for go binary module type which supports Go command packages without running tests.
func BinFactory() (blueprint.Module, []interface{}) {
	mType := &testedBinaryModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
