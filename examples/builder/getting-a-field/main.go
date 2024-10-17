package main

import (
	"fmt"

	"github.com/MartinSimango/dstruct"
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
}
