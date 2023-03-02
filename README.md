# dstruct
A golang package that allows one to create, modify and generate structs dynamically. 

Features:
* Building structs at runtime
* Extending existing struct at runtime
* Merging multiple structs at runtime
* Adding new fields into struct
* Removing existing fields from struct
* Modifying field values in structs
* Reading field values in strucst
* Generating struct values


Limitations:
* You cannot extend structs with unexported embedded fields
* If a struct pointer is nil the struct's subfields won't be added to the dynamic struct
For example
```go

type NilStructPointer struct {
	Field  int
}

type Struct struct {
	Field string
	NilField *NilStructPointer
}

func main() {
	// This will create a struct with fields 
	// * Field
	// * NilField 
	dstruct.ExtendStruct(Struct{}).Build()

	// However this will create a struct with fields 
	// * Field
	// * NilField 
	// * NilField.Field
	dstruct.ExtendStruct(Struct{NilField: &NilStructPointer{}}).Build()

}

```
* Dynamic structs with struct fields of type `any (interface {})` cannot be created. If you try
extend or merge structs which have struct fields of type `any` their value must be set to a concrete type. 

For example

```go
type A struct {
	Field any
}

func main() {
	// This will be fine
	dstruct.ExtendStruct(A{Field: 2})

	// This will panic as Field's type cannot be determined
	dstruct.ExtendStruct(A{})

	// This will also panic as Field's type cannot be determined
	dstruct.NewBuilder().AddField("Field", nil, `json:"Field"`)

}

```


## Sections
* [How it works?](https://github.com/MartinSimango/dstruct#how-it-works)
* [Using the builder](https://github.com/MartinSimango/dstruct#using-the-builder)
* [Using the modifier](https://github.com/MartinSimango/dstruct#using-the-modifier)
* [Using the struct generator](https://github.com/MartinSimango/dstruct#using-a-struct-generator)
* [Extending a struct](https://github.com/MartinSimango/dstruct#extending-a-struct)
* [Merging structs](https://github.com/MartinSimango/dstruct#merging-structs)






## How it works?

Dstruct uses a tree to represent dynamic structs which allows these structs to easily to be manipulated. Nodes and their children represent struct fields and their subfields respectively. Once the tree structure of the struct is created the tree is converted into a dynamic struct using the reflect package. 


Dstruct has 3 main interfaces that are implemented in order to allow these features: 

1. ```dstruct.Builder``` is responsible for adding and removing fields from a struct.

  ```go
  type Builder interface {
      AddField(name string, value interface{}, tag string) Builder
      Build() DynamicStructModifier
      GetField(name string) Builder
      GetFieldCopy(field string) Builder
      RemoveField(name string) Builder
  }
  ```

2. ```dstruct.DynamicStructModifier``` is responsible for reading and editing fields with the struct as well as storing the actual struct.

  ``` go

  type DynamicStructModifier interface {
      Instance() any
      New() any
      Get(field string) (any, error)
      Set(field string, value any) error
      GetFields() map[string]field
  }

  ```

3. ```dstruct.GeneratedStruct``` is responsible for generating struct values and is an extension of the DynamicStructModifier. A generated struct values
are randomly generation based on Generation functions.

```go
type GeneratedStruct interface {
    DynamicStructModifier
    Generate()
    GetFieldGenerationConfig(field string) *generator.GenerationConfig
    SetFieldGenerationConfig(field string, generationConfig *generator.GenerationConfig) error
}

```


## Using the Builder


```go

type Person struct {
	Name string
	Age  int
}

func main() {
	structBuilder := dstruct.NewBuilder().
		AddField("Person", Person{Name: "Martin", Age: 25}, `json:"person"`).
		AddField("Job", "Software Developer", "").
		RemoveField("Person.Age")

	fmt.Printf("Struct: %+v\n", structBuilder.Build().Instance())
}

```
Output
```sh
$ Struct: {Person:{Name:Martin} Job:Software Developer}
```

## Using the Modifier
```go

type Person struct {
	Name string
	Age  int
}

func main() {
	structBuilder := dstruct.NewBuilder().
		AddField("Person", Person{Name: "Martin", Age: 25}, `json:"person"`).
		AddField("Job", "Software Developer", "")

	structModifier := structBuilder.Build()
	structModifier.Set("Person.Name", "Martin Simango")
	structModifier.Set("Job", "Software Engineer")

	name, _ := structModifier.Get("Person.Name")

	fmt.Printf("New name: %s\n", name.(string))
	fmt.Printf("Struct: %+v\n", structModifier.Instance())
}

```
Output
```sh
$ New name: Martin Simango
$ Struct: {Person:{Name:Martin Simango Age:25} Job:Software Engineer}
```


## Using the Struct Generator

```go

type Person struct {
	Name string
	Age  int
}

func main() {
	structBuilder := dstruct.NewBuilder().
		AddField("Person", Person{Name: "Martin", Age: 25}, `json:"person"`).
		AddField("Job", "Software Developer", "")

	strct := structBuilder.Build().Instance()

	generatedStruct := dstruct.NewGeneratedStruct(strct)
	// change the age to be between 50 and 60
	generatedStruct.GetFieldGenerationConfig("Person.Age").SetIntMin(50).SetIntMax(60)
	generatedStruct.Generate()
	fmt.Printf("Struct with age between 50 and 60: %+v\n", generatedStruct.Instance())
	// change the age to be between 10 and 20
	generatedStruct.GetFieldGenerationConfig("Person.Age").SetIntMin(10).SetIntMax(20)
	generatedStruct.Generate()
	fmt.Printf("Struct with age between 10 and 20: %+v\n", generatedStruct.Instance())

}

```

Output:
```sh
$ Struct with age between 50 and 60: {Person:{Name:string Age:59} Job:string}
$ Struct with age between 10 and 20: {Person:{Name:string Age:16} Job:string}
```

## Extending a struct

```go
type Address struct {
	Street string
}

type Person struct {
	Name    string
	Age     int
	Address Address
}

func main() {
	strct := dstruct.ExtendStruct(Person{
		Name: "Martin",
		Age:  25,
		Address: Address{
			Street: "Alice Street",
		},
	})
	strct.GetField("Address").AddField("StreetNumber", 1, "")

	fmt.Printf("Extended struct: %+v\n", strct.Build().Instance())

}

```
Output:
```sh
$ Extended struct: {Name:Martin Age:25 Address:{Street:Alice Street StreetNumber:1}}
```

## Merging structs

```go
type PersonDetails struct {
	Age    int
	Height float64
}

type PersonName struct {
	Name    string
	Surname string
}

func main() {
	strct, _ := dstruct.MergeStructs(PersonDetails{Age: 0, Height: 190.4}, PersonName{Name: "Martin", Surname: "Simango"})
	fmt.Printf("Merged structs: %+v\n", strct)
}

```

Output:
```sh
$ Merged structs: {Age:0 Height:190.4 Name:Martin Surname:Simango}
```