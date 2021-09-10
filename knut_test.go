package knut

import "testing"

type Config struct {
	ValString  string
	ValBoolean bool

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
		t.Error("Expected ", expected, ", got ", got)
	}
}

func TestKnut(t *testing.T) {

	var config Config
	err := Unmarshal("test/testData.txt", &config)

	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	validate(t, config.ValString, "Knut")
	validate(t, config.ValBoolean, true)

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
