# Knut

This is a simple module for loading data from a file into a struct.
The data in the file can be commented to provide clarity to the user when the field name is not enough or ambiguous.

Test converage: **83.3%**

## Why
Normally, when loading a configuration file so would I just use JSON. However, JSON does not support comments, which I want for my projects as it provides added clarity and user-friendliness.


I am making this module an exercise for myself and as something that I personally need.

## Usage

Comments and variables can be defined this way in a normal txt file.
```
; This is a comment for the variable that is below it.
; It should provide a description or an example that helps the user
; properly understand the purpose of the variable and how to use it.
Money=10520

Currency=Krona

Sunny=true

; In celsius
Temperature=23.48
```

The field followed by a '=' and then directly the value without and spaces. Empty lines and lines starting with ';' will be ignored.

The variables will be loaded into a struct by the program using reflection.

```go
type config struct {
    Currency        string
    Money           int
    Sunny           bool
    Temperature     float32
}

var config Config
err := knut.Unmarshal("config.txt", &config)
```

## Roadmap
- [X] Basic functionality
- [X] Support for basic types: string, int and bool
- [X] Error handling
- [X] Testing
- [X] Support for all integer types
- [ ] Support for arrays
- [ ] Greater test coverage
