package example_bin

import "fmt"

func getName(str1, str2 string) string {
	var name = str1 + " " + str2 + " " + "Team"
	return name
}

func main() {
	var name = getName("Red", "Stone")
	fmt.Println(name)
}
