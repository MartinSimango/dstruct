package util

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr
}
