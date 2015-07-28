package template

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/jehiah/go-strftime"
)

var funcs = template.FuncMap{
	"last":       Last,
	"join":       Join,
	"split":      Split,
	"repeat":     Repeat,
	"keys":       Keys,
	"type":       Type,
	"map":        Map,
	"upper":      strings.ToUpper,
	"lower":      strings.ToLower,
	"contains":   strings.Contains,
	"replace":    Replace,
	"trim":       Trim,
	"ltrim":      TrimLeft,
	"rtrim":      TrimRight,
	"default":    Default,
	"center":     Center,
	"random":     Random,
	"capitalize": Capitalize,
	"add":        Add,
	"sub":        Sub,
	"div":        Div,
	"mul":        Mul,
	"lalign":     AlignLeft,
	"ralign":     AlignRight,
	"odd":        Odd,
	"even":       Even,
	"date":       Date,
}

// Parse template.
func Compile(s string, d map[string]interface{}) (*bytes.Buffer, error) {
	t := template.Must(template.New("template").Funcs(funcs).Parse(s))

	b := new(bytes.Buffer)
	err := t.Execute(b, d)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Convert []interface{} to []string{}.
func arrIntfToStr(inp []interface{}) []string {
	outp := make([]string, len(inp))
	for i, val := range inp {
		outp[i] = fmt.Sprintf("%v", val)
	}
	return outp
}

// Convert []string{} to []interface{}.
func arrStrToIntf(inp []string) []interface{} {
	outp := make([]interface{}, len(inp))
	for i, val := range inp {
		outp[i] = val
	}
	return outp
}

// Last element in array.
func Last(i int, inp interface{}) (bool, error) {
	if reflect.ValueOf(inp).Kind() != reflect.Slice && reflect.ValueOf(inp).Kind() != reflect.Array {
		return false, fmt.Errorf("Incorrect type: %s, needs to be: slice or array", reflect.ValueOf(inp).Kind())
	}

	return i == reflect.ValueOf(inp).Len()-1, nil
}

// Join elements in an array to a string.
func Join(sep string, inp []interface{}) string {
	return strings.Join(arrIntfToStr(inp), sep)
}

// Split string into an array.
func Split(sep string, inp string) []interface{} {
	return arrStrToIntf(strings.Split(inp, sep))
}

// Repeat string x number of times.
func Repeat(rep int, inp string) string {
	return strings.Repeat(inp, rep)
}

// Keys from interface{}.
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
	for i := range k {
		k[i] = vk[i].Interface()
	}

	return k, nil
}

// Type of variable (usefull for debugging templates).
func Type(inp interface{}) string {
	return fmt.Sprintf("%v", reflect.TypeOf(inp))
}

// Map returns true if type is a map.
func Map(inp interface{}) bool {
	return reflect.TypeOf(inp).Kind() == reflect.Map
}

// Replace string.
func Replace(oldStr string, newStr string, str string) string {
	return strings.Replace(str, oldStr, newStr, -1)
}

// Trim trims characters on the left and right side of the string.
func Trim(trim string, str string) string {
	return strings.Trim(str, trim)
}

// TrimLeft trims characters on the left side of the string.
func TrimLeft(trim string, str string) string {
	return strings.TrimLeft(str, trim)
}

// TrimRight trims characters on the right side of string.
func TrimRight(trim string, str string) string {
	return strings.TrimRight(str, trim)
}

// Default returns default if argument is empty.
func Default(def interface{}, inpOpt ...interface{}) interface{} {
	if len(inpOpt) > 0 {
		def = inpOpt[0]
	}
	return def
}

// Center centers string.
func Center(size int, str string) string {
	if size < len(str) {
		return str
	}

	pad := (size - len(str)) / 2
	lpad := pad
	rpad := size - len(str) - lpad

	return fmt.Sprintf("%s%s%s", strings.Repeat(" ", lpad), str, strings.Repeat(" ", rpad))
}

// Random returns random number.
func Random(size int) int {
	return rand.Intn(size)
}

// Capitalize capitalizes first character in string.
func Capitalize(str string) string {
	for i, v := range str {
		return strings.ToUpper(string(v)) + str[i+1:]
	}
	return ""
}

// Add do additions to number.
func Add(y int, x int) int {
	return x + y
}

// Sub do subtractions to number.
func Sub(y int, x int) int {
	return x - y
}

// Div do divisions to number.
func Div(y int, x int) int {
	return x / y
}

// Mul do multiplication to number.
func Mul(y int, x int) int {
	return x * y
}

// AlignLeft aligns text to the left.
func AlignLeft(size int, str string) string {
	if size < len(str) {
		return str
	}

	pad := (size - len(str))

	return fmt.Sprintf("%s%s", str, strings.Repeat(" ", pad))
}

// AlignRight aligns text to the right.
func AlignRight(size int, str string) string {
	if size < len(str) {
		return str
	}

	pad := (size - len(str))

	return fmt.Sprintf("%s%s", strings.Repeat(" ", pad), str)
}

// Odd return true if it's an odd number.
func Odd(x int) bool {
	if (x%2 - 1) == 0 {
		return true
	}
	return false
}

// Even return true if it's an even number.
func Even(x int) bool {
	if (x % 2) == 0 {
		return true
	}
	return false
}

// Date return date as formated string.
func Date(args ...interface{}) string {
	if len(args) == 1 {
		return strftime.Format(args[0].(string), time.Now())
	}
	return strftime.Format("%Y-%m-%d %H:%M:%S", time.Now())
}
