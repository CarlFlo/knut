package knut

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
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
	case reflect.Int8:
		err := handleInt(&value, 8, &field)
		if err != nil {
			return err
		}
	case reflect.Int16:
		err := handleInt(&value, 16, &field)
		if err != nil {
			return err
		}
	case reflect.Int32:
		err := handleInt(&value, 32, &field)
		if err != nil {
			return err
		}
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		err := handleInt(&value, 64, &field)
		if err != nil {
			return err
		}

	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		err := handleBool(&value, &field)
		if err != nil {
			return err
		}
	case reflect.Float32:
		handlefloat(&value, 32, &field)

	case reflect.Float64:
		handlefloat(&value, 64, &field)

	default:
		return fmt.Errorf("'%s' is currently unsupported", field.Kind())
	}

	return nil
}

func handleInt(value *string, bitsize int, field *reflect.Value) error {
	parsedInt, err := strconv.ParseInt(*value, 10, bitsize)
	if err != nil {
		return err
	}

	if field.OverflowInt(parsedInt) {
		return fmt.Errorf("value '%s' is too large for %d bit integer", *value, bitsize)
	}

	field.SetInt(parsedInt)
	return nil
}

func handleBool(value *string, field *reflect.Value) error {
	inputToBool, err := strconv.ParseBool(*value)
	if err != nil {
		return err
	}
	field.SetBool(inputToBool)
	return nil
}

func handlefloat(value *string, bitsize int, field *reflect.Value) error {
	// bitsize 32 for float32, 64 for float64
	inputToFloat, err := strconv.ParseFloat(*value, bitsize)
	if err != nil {
		return err
	}
	field.SetFloat(inputToFloat)
	return nil
}
