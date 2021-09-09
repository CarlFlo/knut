package knut

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/CarlFlo/malm"
)

// Unmarshal will load the file from the provided filepath and unmarshal it into the struct
func Unmarshal(path string, v interface{}) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	elem := reflect.ValueOf(v).Elem()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// remove trailing whitespace
		line = strings.TrimSpace(line)

		if len(line) == 0 || line[0] == ';' {
			// No data or a comment line
			continue
		}

		// process the line
		// split line on '=', left is variable name, right is value
		parts := strings.Split(line, "=")

		if len(parts) != 2 {
			return fmt.Errorf("invalid line: %s", line)
		}

		setFieldInStruct(parts[0], parts[1], elem)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func setFieldInStruct(fieldName, value string, elem reflect.Value) error {

	// cast interface to struct
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("invalid type '%s'", elem.Kind())
	}

	// exported field
	field := elem.FieldByName(fieldName)

	// Check if field is valid and can be set
	if !field.IsValid() && !field.CanSet() {
		return fmt.Errorf("invalid field '%s'", fieldName)
	}

	switch field.Kind() {
	case reflect.Int:

		// cast string to int
		inputToInt, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("%s", err)
		}

		val := int64(inputToInt)
		if !field.OverflowInt(val) {
			field.SetInt(val)
		}

	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		inputToBool, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(inputToBool)

	default:
		malm.Info("'%s' is currently unsupported", field.Kind())
	}

	return nil
}
