package main

import (
    "fmt"
    "strings"
    "reflect"
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
func arrLast(i int, inp interface{}) bool {
    return i == reflect.ValueOf(inp).Len() - 1
}

// Join elements in an array to a string
func arrJoin(sep string, inp []interface{}) string {
    return strings.Join(arrIntfToStr(inp), sep)
}

// Split string into an array
func strSplit(sep string, inp string) []interface{} {
    return arrStrToIntf(strings.Split(inp, sep))
}

// Repeat string x number of times
func strRepeat(rep int, inp string) string {
        return strings.Repeat(inp, rep)
}
