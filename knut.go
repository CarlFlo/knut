package knut

import (
	"bufio"
	"fmt"
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

	switch field.Kind() {
	case reflect.Int8:
		return handleInt(&value, 8, &field)
	case reflect.Int16:
		return handleInt(&value, 16, &field)
	case reflect.Int32:
		return handleInt(&value, 32, &field)
	case reflect.Int, reflect.Int64:
		return handleInt(&value, 64, &field)
	case reflect.Uint8:
		return handleUInt(&value, 8, &field)
	case reflect.Uint16:
		return handleUInt(&value, 16, &field)
	case reflect.Uint32:
		return handleUInt(&value, 32, &field)
	case reflect.Uint, reflect.Uint64:
		return handleUInt(&value, 64, &field)
	case reflect.String:

		// remove leading and trailing ' and " characters
		if len(value) >= 2 {
			if value[0] == '\'' && value[len(value)-1] == '\'' {
				value = value[1 : len(value)-1]
			} else if value[0] == '"' && value[len(value)-1] == '"' {
				value = value[1 : len(value)-1]
			}
		}

		field.SetString(value)
	case reflect.Bool:
		return handleBool(&value, &field)
	case reflect.Float32:
		return handlefloat(&value, 32, &field)
	case reflect.Float64:
		return handlefloat(&value, 64, &field)
	case reflect.Slice:
		return handleSlice(&value, &field)
	default:
		return fmt.Errorf("'%s' is currently unsupported", field.Kind())
	}

	return nil
}

func handleSlice(value *string, field *reflect.Value) error {

	// This function parses a string in this format [value1, value2, value3, ...] and returns a slice
	// of the parsed values

	*value = (*value)[1 : len(*value)-1]

	// Split the string into individual values
	values := strings.Split(*value, ",")

	// Create a new slice of the appropriate type
	slice := reflect.MakeSlice(field.Type(), len(values), len(values))

	// Parse each value and add it to the slice
	for i, v := range values {
		v = strings.TrimSpace(v)

		// remove leading and trailing ' and " characters
		if len(v) >= 2 {
			if v[0] == '\'' && v[len(v)-1] == '\'' {
				v = v[1 : len(v)-1]
			} else if v[0] == '"' && v[len(v)-1] == '"' {
				v = v[1 : len(v)-1]
			}
		}

		parsedValue, err := parseValue(v, field.Type().Elem())
		if err != nil {
			return err
		}
		slice.Index(i).Set(parsedValue)
	}

	// Set the value of the field to the new slice
	field.Set(slice)

	return nil
}

// parseValue parses a string into a value of the given type
func parseValue(value string, t reflect.Type) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(strings.ToLower(value) == "true"), nil
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

func handleInt(value *string, bitsize int, field *reflect.Value) error {
	parsedInt, err := strconv.ParseInt(*value, 10, bitsize)
	if err != nil {
		return err
	}

	field.SetInt(parsedInt)
	return nil
}

func handleUInt(value *string, bitsize int, field *reflect.Value) error {
	parsedUInt, err := strconv.ParseUint(*value, 10, bitsize)
	if err != nil {
		return err
	}

	field.SetUint(parsedUInt)
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

func handleBool(value *string, field *reflect.Value) error {
	inputToBool, err := strconv.ParseBool(*value)
	if err != nil {
		return err
	}
	field.SetBool(inputToBool)
	return nil
}
