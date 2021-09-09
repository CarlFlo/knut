package knut

import "testing"

type Config struct {
	Name       string
	Number     int
	IsTrue     bool
	Uint8Value uint8
}

func TestKnut(t *testing.T) {

	var config Config
	Unmarshal("test/testData.txt", &config)

	if config.Name != "Knut" {
		t.Error("Expected 'Knut', got ", config.Name)
	} else if config.Number != 532 {
		t.Error("Expected '532', got ", config.Number)
	} else if config.IsTrue != true {
		t.Error("Expected 'true', got ", config.IsTrue)
	}
}

func TestInvalidPath(t *testing.T) {
	var config Config
	err := Unmarshal("test/not a real file", &config)
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestInvalidType(t *testing.T) {
	var config Config
	err := Unmarshal("test/testData_invalid.txt", &config)
	if err != nil {
		t.Error("Expected an error about an invalid type")
	}
}
