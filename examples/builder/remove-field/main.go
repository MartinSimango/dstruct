package main

import (
	"fmt"

	"github.com/MartinSimango/dstruct"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	var builder dstruct.Builder = dstruct.NewBuilder().
		AddField("Person", Person{Name: "Martin", Age: 25}, `json:"person"`).
		AddField("Job", "Software Developer", "")

	fmt.Printf("Struct before removing fields:\n%+v\n", builder.Build().Instance())

	builder.RemoveField("Person.Age")
	builder.RemoveField("Job")

	fmt.Printf("Struct after removing fields:\n%+v\n", builder.Build().Instance())
}