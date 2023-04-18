# Knut

![Tests](https://github.com/CarlFlo/knut/actions/workflows/go.yml/badge.svg)

Knut is a simple module for loading data from a file into a struct.
The data in the file can be commented to provide clarity to the user when the field name is not enough or ambiguous.

Test converage: **84.6%**

## Why
Typically, JSON is used to load configuration files. However, JSON lacks support for comments, which I require in my projects to enhance clarity and provide a better user experience. 

Thus, I am creating this module as a personal exercise and to fulfill my own needs.

## Usage

Comments and variables can be defined this way in a normal txt file.
> Note: A comment and variable cannot share the same line.
```
# This is a comment for the variable that is below it.
# It should provide a description or an example that helps the user
# properly understand the purpose of the variable and how to use it.
Money=10520

Currency=Krona

# Trailing whitespaces are otherwise removed
Text=' this string will keep its spaces   '

Sunny=true
# There can be a space between the variable, =, and value
Values = [1,2,3,4,5]
HelloWorld = [Hello, World]

# In celsius
Temperature=23.48
Temperatures=[21.4, 18.8, 15.43]
```

Empty lines and lines starting with '#' will be ignored.

The variables will be loaded into a struct by the program using reflection.

```go
type Config struct {
    Currency        string
    Text            string
    Money           int
    Sunny           bool
    HelloWorld      []string
    Temperature     float32
    Temperatures    []float32
}

var config Config
err := knut.Unmarshal("config.txt", &config)
```

### Supported types
- int8, int16, int32 & int
- uint8, uint16, uint32 & uint
- string
- bool
- float32 & float64

Slices for all listed types above are also supported

### Format

Whitespaces will be trimmed, unless the text is encapsulated with single or double quotes ' "

```
# These are valid
      Name=   Knut   
Code    =   42   
ExtraSpaces = '    this string will have extra trailing spaces      '
```


## Roadmap
- [X] Basic functionality
- [X] Support for basic types: string, int and bool
- [X] Error handling
- [X] Testing
- [X] Support for all integer types and float
- [X] Support for slices
- [ ] Support for maps
- [ ] Support for slices with custom types
- [ ] Greater test coverage
