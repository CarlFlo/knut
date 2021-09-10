package knut

import "testing"

type Config struct {
	Name       string
	Number     int
	IsTrue     bool
	F32        float32
	F64        float64
	ValInt8    int8
	ValInt16   int16
	ValInt32   int32
	ValInt64   int64
	Uint8Value uint8
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

	validate(t, config.Name, "Knut")
	validate(t, config.Number, int(532))
	validate(t, config.IsTrue, true)
	validate(t, config.F32, float32(99.99))
	validate(t, config.F64, float64(234.49))
	validate(t, config.ValInt8, int8(-128))
	validate(t, config.ValInt16, int16(32767))
	validate(t, config.ValInt32, int32(2147483647))
	validate(t, config.ValInt64, int64(9223372036854775807))
}

func TestInvalidPath(t *testing.T) {
	var config Config
	err := Unmarshal("test/not a real file", &config)
	if err == nil {
		t.Error("Expected an error")
	}
}
