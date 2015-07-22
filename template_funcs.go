package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/mickep76/tf/vendor/github.com/jehiah/go-strftime"
)

// Convert []interface{} to []string{}
func arrIntfToStr(inp []interface{}) []string {
	outp := make([]string, len(inp))
	for i, val := range inp {
		outp[i] = fmt.Sprintf("%v", val)
	}
	return outp
}

// Convert []string{} to []interface{}
func arrStrToIntf(inp []string) []interface{} {
	outp := make([]interface{}, len(inp))
	for i, val := range inp {
		outp[i] = val
	}
	return outp
}

// Determine if index is the last element in the array
func IsLast(i int, inp interface{}) bool {
	return i == reflect.ValueOf(inp).Len()-1
}

// Join elements in an array to a string
func Join(sep string, inp []interface{}) string {
	return strings.Join(arrIntfToStr(inp), sep)
}

// Split string into an array
func Split(sep string, inp string) []interface{} {
	return arrStrToIntf(strings.Split(inp, sep))
}

// Repeat string x number of times
func Repeat(rep int, inp string) string {
	return strings.Repeat(inp, rep)
}

// Get keys from interface{}
func Keys(inp interface{}) (interface{}, error) {
	if inp == nil {
		return nil, nil
	}

	val := reflect.ValueOf(inp)
	if val.Kind() != reflect.Map {
		return nil, fmt.Errorf("Cannot call keys on a non-map value: %v", inp)
	}

	vk := val.MapKeys()
	k := make([]interface{}, val.Len())
	for i, _ := range k {
		k[i] = vk[i].Interface()
	}

	return k, nil
}

// Get type (usefull for debugging templates)
func Type(inp interface{}) string {
	return fmt.Sprintf("%v", reflect.TypeOf(inp))
}

// Test if type is a map i.e. not printable
func IsMap(inp interface{}) bool {
	return reflect.TypeOf(inp).Kind() == reflect.Map
}

// String replace
func Replace(oldStr string, newStr string, str string) string {
	return strings.Replace(str, oldStr, newStr, -1)
}

// String trim
func Trim(trim string, str string) string {
	return strings.Trim(str, trim)
}

// String trim left
func TrimLeft(trim string, str string) string {
	return strings.TrimLeft(str, trim)
}

// String trim right
func TrimRight(trim string, str string) string {
	return strings.TrimRight(str, trim)
}

// If no value is passed for the second arg. it returns the default
func Default(def interface{}, inp_opt ...interface{}) interface{} {
	if len(inp_opt) > 0 {
		def = inp_opt[0]
	}
	return def
}

func Center(size int, str string) string {
	if size < len(str) {
		return str
	}

	pad := (size - len(str)) / 2
	lpad := pad
	rpad := size - len(str) - lpad

	return fmt.Sprintf("%s%s%s", strings.Repeat(" ", lpad), str, strings.Repeat(" ", rpad))
}

func Random(size int) int {
	return rand.Intn(size)
}

// Capitalize first character in string
func Capitalize(str string) string {
	for i, v := range str {
		return strings.ToUpper(string(v)) + str[i+1:]
	}
	return ""
}

func Add(y int, x int) int {
	return x + y
}

func Sub(y int, x int) int {
	return x - y
}

func Div(y int, x int) int {
	return x / y
}

func Mul(y int, x int) int {
	return x * y
}

func AlignLeft(size int, str string) string {
	if size < len(str) {
		return str
	}

	pad := (size - len(str))

	return fmt.Sprintf("%s%s", str, strings.Repeat(" ", pad))
}

func AlignRight(size int, str string) string {
	if size < len(str) {
		return str
	}

	pad := (size - len(str))

	return fmt.Sprintf("%s%s", strings.Repeat(" ", pad), str)
}

func Odd(x int) bool {
	if (x%2 - 1) == 0 {
		return true
	}
	return false
}

func Even(x int) bool {
	if (x % 2) == 0 {
		return true
	}
	return false
}

func Date(args ...interface{}) string {
	if len(args) == 1 {
		return strftime.Format(args[0].(string), time.Now())
	}
	return strftime.Format("%Y-%m-%d %H:%M:%S", time.Now())
}
