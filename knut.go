package knut

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type fieldHandler func(value *string, field *reflect.Value) error

var fieldHandlers = map[reflect.Kind]fieldHandler{
	reflect.Int8:    handleInt,
	reflect.Int16:   handleInt,
	reflect.Int32:   handleInt,
	reflect.Int:     handleInt,
	reflect.Int64:   handleInt,
	reflect.Uint8:   handleUInt,
	reflect.Uint16:  handleUInt,
	reflect.Uint32:  handleUInt,
	reflect.Uint:    handleUInt,
	reflect.Uint64:  handleUInt,
	reflect.String:  handleString,
	reflect.Bool:    handleBool,
	reflect.Float32: handleFloat,
	reflect.Float64: handleFloat,
	reflect.Slice:   handleSlice,
}

// Unmarshal will load the file from the provided filepath and unmarshal it into the struct
func Unmarshal(path string, v interface{}) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	elem := reflect.ValueOf(v).Elem()

	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("invalid type '%s'", elem.Kind())
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// remove trailing whitespace on the left side only
		line = strings.TrimLeft(line, " ")

		if len(line) == 0 || line[0] == '#' {
			// No data or a comment line
			continue
		}

		// process the line
		// split line on '=', left is variable name, right is value
		parts := strings.Split(line, "=")

		if len(parts) != 2 {
			return fmt.Errorf("invalid line: %s", line)
		}

		// Remove trailing whitespaces from the variable name
		parts[0] = strings.TrimRight(parts[0], " ")

		// Trims trailing whitespaces from the value
		parts[1] = strings.TrimSpace(parts[1])

		err := setFieldInStruct(parts[0], parts[1], elem)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func setFieldInStruct(fieldName, value string, elem reflect.Value) error {

	// exported field
	field := elem.FieldByName(fieldName)

	// Check if field is valid and can be set
	if !field.IsValid() && !field.CanSet() {
		return fmt.Errorf("invalid field '%s'", fieldName)
	}

	if handler, ok := fieldHandlers[field.Kind()]; ok {
		return handler(&value, &field)
	}

	return fmt.Errorf("'%s' is currently unsupported", field.Kind())
}

func handleSlice(value *string, field *reflect.Value) error {

	// This function parses a string in this format [value1, value2, value3, ...] and returns a slice
	// of the parsed values

	if (*value)[0] != '[' && len(*value)-1 != ']' {
		return fmt.Errorf("invalid slice format: %s", *value)
	}

	*value = (*value)[1 : len(*value)-1]

	// Split the string into individual values
	values := strings.Split(*value, ",")

	// Create a new slice of the appropriate type
	slice := reflect.MakeSlice(field.Type(), len(values), len(values))

	// Parse each value and add it to the slice
	for i, value := range values {
		value = strings.TrimSpace(value)

		preProcessString(&value)

		parsedValue, err := parseSliceValue(value, field.Type().Elem())
		if err != nil {
			return err
		}
		slice.Index(i).Set(parsedValue)
	}

	// Set the value of the field to the new slice
	field.Set(slice)

	return nil
}

// parseSliceValue parses a string into a value of the given type
func parseSliceValue(value string, t reflect.Type) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.Bool:

		inputToBool, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid input for bool: %s", value)
		}
		return reflect.ValueOf(inputToBool), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsedValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(parsedValue).Convert(t), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsedValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(parsedValue).Convert(t), nil

	case reflect.Float32, reflect.Float64:
		parsedValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(parsedValue).Convert(t), nil

	case reflect.String:
		return reflect.ValueOf(value), nil

	default:
		return reflect.Value{}, fmt.Errorf("unsupported type: %v", t)
	}
}

func handleString(value *string, field *reflect.Value) error {
	preProcessString(value)
	field.SetString(*value)
	return nil
}

func handleInt(value *string, field *reflect.Value) error {
	parsedInt, err := strconv.ParseInt(*value, 10, 64)
	if err != nil {
		return err
	}

	field.SetInt(parsedInt)
	return nil
}

func handleUInt(value *string, field *reflect.Value) error {
	parsedUInt, err := strconv.ParseUint(*value, 10, 64)
	if err != nil {
		return err
	}

	field.SetUint(parsedUInt)
	return nil
}

func handleFloat(value *string, field *reflect.Value) error {
	inputToFloat, err := strconv.ParseFloat(*value, 64)
	if err != nil {
		return err
	}
	field.SetFloat(inputToFloat)
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

// remove leading and trailing ' and " characters
func preProcessString(value *string) {

	if len(*value) >= 2 {
		if (*value)[0] == '\'' && (*value)[len((*value))-1] == '\'' {
			*value = (*value)[1 : len((*value))-1]
		} else if (*value)[0] == '"' && (*value)[len(*value)-1] == '"' {
			*value = (*value)[1 : len((*value))-1]
		}
	}
}
