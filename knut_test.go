package knut

import (
	"reflect"
	"testing"
)

type Config struct {
	ValString         string
	ValString2        string
	ValTrailingString string
	ValBoolean        bool

	ValSliceString []string
	ValSliceInt    []int
	ValSliceUInt   []uint
	ValSliceF32    []float32
	ValSliceF64    []float64
	ValSliceBool   []bool

	F32 float32
	F64 float64

	ValInt8  int8
	ValInt16 int16
	ValInt32 int32
	ValInt64 int64
	ValInt   int

	ValUInt8  uint8
	ValUInt16 uint16
	ValUInt32 uint32
	ValUInt64 uint64
	ValUInt   uint
}

func validate(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Errorf("Expected %v (%s), got %v (%s)", expected, reflect.TypeOf(expected), got, reflect.TypeOf(got))
	}
}

func TestKnut(t *testing.T) {

	var config Config
	err := Unmarshal("test/testData.txt", &config)

	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	validate(t, config.ValString, "Knut")
	validate(t, config.ValString2, "   Knut       ")
	validate(t, config.ValTrailingString, "noTrail")
	validate(t, config.ValBoolean, true)

	validate(t, config.ValSliceString[0], "1")
	validate(t, config.ValSliceString[1], "2")
	validate(t, config.ValSliceString[2], "3")
	validate(t, config.ValSliceString[3], "4")
	validate(t, config.ValSliceString[4], " extra white spaces     ")

	validate(t, config.ValSliceInt[0], 1)
	validate(t, config.ValSliceInt[1], 2)
	validate(t, config.ValSliceInt[2], 3)
	validate(t, config.ValSliceInt[3], -4)

	validate(t, config.ValSliceUInt[0], uint(1))
	validate(t, config.ValSliceUInt[1], uint(2))
	validate(t, config.ValSliceUInt[2], uint(3))
	validate(t, config.ValSliceUInt[3], uint(503))

	validate(t, config.ValSliceF32[0], float32(35.56))
	validate(t, config.ValSliceF32[1], float32(1005.99))
	validate(t, config.ValSliceF32[2], float32(9999.99))

	validate(t, config.ValSliceF64[0], float64(1.56))
	validate(t, config.ValSliceF64[1], float64(89765.9))
	validate(t, config.ValSliceF64[2], float64(12345.67))

	validate(t, config.ValSliceBool[0], true)
	validate(t, config.ValSliceBool[1], true)
	validate(t, config.ValSliceBool[2], false)

	validate(t, config.F32, float32(99.99))
	validate(t, config.F64, float64(234.49))

	validate(t, config.ValInt8, int8(-128))
	validate(t, config.ValInt16, int16(32767))
	validate(t, config.ValInt32, int32(2147483647))
	validate(t, config.ValInt64, int64(9223372036854775807))
	validate(t, config.ValInt, int(-9223372036854775808))

	validate(t, config.ValUInt8, uint8(255))
	validate(t, config.ValUInt16, uint16(65535))
	validate(t, config.ValUInt32, uint32(4294967295))
	validate(t, config.ValUInt64, uint64(18446744073709551615))
	validate(t, config.ValUInt, uint(18446744073709551615))
}

func TestInvalidPath(t *testing.T) {
	var config Config
	err := Unmarshal("test/not a real file", &config)
	if err == nil {
		t.Error("Expected an error")
	}
}
