package main

import (
	"fmt"

	"github.com/MartinSimango/dstruct"
)

type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Phone Phone  `json:"phone"`
}

type Phone struct {
	number string
	model  string
}

func main() {
	builder := dstruct.NewBuilder().
		AddField("Person", Person{Name: "Martin", Age: 25, Phone: Phone{number: "123456789", model: "iPhone 16"}}, `json:"person"`).
		AddField("Job", "Software Developer", "")

	// Builder.Instance() will return the struct instance of type any

	fmt.Printf("Struct before removing fields:\n%+v\n\n", builder.Build().Instance())

	// Can remove nested fields
	builder.RemoveField("Person.Age").
		RemoveField("Job").
		RemoveField("UnknownField"). // this will simply be ignored
		RemoveField("Person.Phone.number")

	fmt.Printf("Struct after removing fields:\n%+v\n", builder.Build().Instance())
}
