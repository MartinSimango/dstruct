package core

import (
	"time"

	"github.com/MartinSimango/dstruct/generator"
)

const ISO8601 string = "2018-03-20T09:12:28Z"

func GenerateDateTimeFunc() generator.GenerationFunction {

	// TODO have a proper implementation
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			return time.Now().UTC().Format(time.RFC3339)
		},
	}

}

func GenerateDateTimeBetweenDatesFunc(startDate, endDate time.Time) generator.GenerationFunction {
	// TODO implement
	// f := generateDateTime
	// return f
	return nil
}
