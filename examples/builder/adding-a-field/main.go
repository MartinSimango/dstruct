package main

import (
	"fmt"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/examples"
)

type Address struct {
	Street string `json:"street"`
}

type Person struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Address Address `json:"address"`
}

type Phone struct {
	Number string `json:"number"`
}

func main() {
	builder := dstruct.NewBuilder().
		AddField("Person", Person{Name: "Martin", Age: 25, Address: Address{Street: "Jackson Street"}}, `json:"person"`).
		AddField("Job", "Software Developer", "")

	fmt.Printf("Struct after adding fields:\n%+v\n\n", builder.Build())

	// Add fields to the nested struct

	builder.GetField("Person").
		AddField("Height", 175, `json:"height"`).
		AddEmbeddedField(Phone{Number: "123456789"}, "") // Add an embedded field

	builder.GetField("Person.Address").AddField("City", "New York", `json:"city"`)

	fmt.Printf("Struct after adding fields to the nested struct:\n%+v\n", builder.Build())

	// Uncomment the lines below to see the panic
	errorExample_1()
	errorExample_2()
}

func errorExample_1() {
	defer examples.RecoverFromPanic()
	builder := dstruct.NewBuilder().
		AddField("Person", Person{Name: "Martin", Age: 25, Address: Address{Street: "Jackson Street"}}, `json:"person"`).
		AddField("job", "Software Developer", "")

	// This will panic because the field "job" is unexported
	builder.Build()
}

func errorExample_2() {
	defer examples.RecoverFromPanic()
	builder := dstruct.NewBuilder().
		AddField("Person", Person{Name: "Martin", Age: 25, Address: Address{Street: "Jackson Street"}}, `json:"person"`).
		AddField("Job Titile", "Software Developer", "")

	builder.Build() // This will panic because the field "Job Titile" is an invalid struct field name
}
