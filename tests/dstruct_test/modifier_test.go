package dstruct_test

import (
	"testing"

	"github.com/MartinSimango/dstruct"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	assert := assert.New(t)

	type Person struct {
		Age int
	}

	b := dstruct.ExtendStruct(Person{Age: 2}).Build()
	b.Set("Age", 20)

	val, _ := b.Get("Age")

	assert.EqualValues(20, val)

}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	type Person struct {
		Age int
	}

	val, _ := dstruct.ExtendStruct(Person{Age: 2}).Build().Get("Age")

	assert.EqualValues(2, val)
}
