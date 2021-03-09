package gomodule

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type PlugResolver struct {
	errors []string
}

// Checks if slice contains element
func contains(slice []string, el string) bool {
	for _, a := range slice {
		if a == el {
			return true
		}
	}
	return false
}

// Resolves patterns
func (pg *PlugResolver) GlobWithDeps(src string, exclude []string) ([]string, error) {
	// if this patterns should produce error, do it
	if contains(pg.errors, src) {
		return []string{}, fmt.Errorf("Wrong")
	}
	// if this patterns is excluded, return nothing
	if contains(exclude, src) {
		return []string{}, nil
	}
	// Otherwise return this pattern
	return []string{src}, nil
}

func TestResolve(t *testing.T) {
	// Patterns we want to resolve
	patters := []string{"a", "b", "c", "a"}
	// Patterns we don't want to resolve
	exclude := []string{"a"}
	// Patterns that will produce errors
	plug := PlugResolver{errors: []string{"b"}}
	resolved, unresolved := resolvePatterns(&plug, patters, exclude)
	assert.True(t, reflect.DeepEqual(resolved, []string{"c"}))
	assert.True(t, reflect.DeepEqual(unresolved, []string{"b"}))
}
