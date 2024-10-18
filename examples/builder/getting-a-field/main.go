package main

import (
	"fmt"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/examples"
)

type Address struct {
	Street string `json:"street"`
}
type Phone struct {
	Number string `json:"number"`
}
type Person struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Address Address `json:"address"`
}

func main() {
	builder := dstruct.NewBuilder().
		AddField("Person", &Person{Name: "Martin", Age: 25, Address: Address{Street: "street"}}, `json:"person"`)

	addressBuilder := builder.GetField("Person.Address")
	addressBuilder.AddField("City", "New York", "")

	fmt.Printf("Address field after adding City:\n%+v\n\n", addressBuilder.Build())
	fmt.Printf("Struct after modifying address field:\n%+v\n\n", builder.Build())

	// Get a copy of the address struct - this does not modify the original struct
	addressBuilderCopy := builder.NewBuilderFromField("Person.Address")
	addressBuilderCopy.AddField("PostalCode", "12345", "")

	fmt.Printf(
		"Copy of Address field after adding PostalCode:\n%+v\n\n",
		addressBuilderCopy.Build(),
	)
	fmt.Printf("Struct after modifying address field (No PostalCode):\n%+v\n\n", builder.Build())

	// Error Examples:
	errorExample_1()
	errorExample_2()
}

func errorExample_1() {
	defer examples.RecoverFromPanic()
	builder := dstruct.NewBuilder().
		AddField("Person", Person{Name: "Martin", Age: 25, Address: Address{Street: "Jackson Street"}}, `json:"person"`)

	// This will panic because the field "People" does not exist
	builder.GetField("People")
}

func errorExample_2() {
	defer examples.RecoverFromPanic()
	builder := dstruct.NewBuilder().
		AddField("Person", Person{Name: "Martin", Age: 25, Address: Address{Street: "Jackson Street"}}, `json:"person"`)

	// This will panic beecause teh field "Street" is not a struct or doesn't fully dereference to a struct value
	builder.GetField("Person.Address.Street")
}
