package main

import (
    "strings"
    "strconv"
    "reflect"
)

// Convert interface{} to string
func intfToStr(inp interface{}) string {
    switch inp.(type) {
        case string:
            return inp.(string)
        case int, int32, int64:
            return strconv.Itoa(inp.(int))
    }

    // Generate some error
    return ""
}

// Convert []interface{} to []string{}
func arrIntfToStr(inp []interface{}) []string {
    outp := make([]string, len(inp))
    for i, val := range inp {
        outp[i] = intfToStr(val)
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

// Compare index with slice length to determine if it's the last element
func arrLast(i int, inp interface{}) bool {
    return i == reflect.ValueOf(inp).Len() - 1
}

// Join elements in an array to a string
func arrJoin(inp []interface{}, sep string) string {
    return strings.Join(arrIntfToStr(inp), sep)
}

// Split string into an array
func strSplit(inp string, sep string) []interface{} {
    return arrStrToIntf(strings.Split(inp, sep))
}

// Repeat string x number of times
func strRepeat(a string, x int) string {
        return strings.Repeat(a, x)
}