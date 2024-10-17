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
		AddField("Person", &Person{Name: "Martin", Age: 25, Address: Address{}}, `json:"person"`)

	builder.GetField("Person").AddField("Height", 175, `json:"height"`)

	addressBuilder := builder.GetField("Person.Address")

	fmt.Printf("Address field:\n%+v\n\n", addressBuilder.Build())

	streetBuilder := addressBuilder.GetField(
		"Street",
	) // this will be an empty struct because Street is a primitive type

	fmt.Printf("Street field:\n%+v\n\n", streetBuilder.Build())

	fmt.Println(addressBuilder.AddField("Cool", 2, "").Build())

	fmt.Println(builder.Build())

	// fmt.Printf("Struct after adding fields:\n%+v\n\n", builder.Build())
	//
	// // Add fields to the nested struct
	//
	// builder.GetField("Person").
	// 	AddField("Height", 175, `json:"height"`)
	//
	// builder.GetField("Person.Address").AddField("City", "New York", `json:"city"`)
	//
	// fmt.Printf("Struct after adding fields to the nested struct:\n%+v\n", builder.Build())
}
