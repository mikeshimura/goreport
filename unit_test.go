package goreport

import (
"testing"
	"fmt"
)

func TestUnit(t *testing.T) {
	fmt.Println(AddComma("12.66"))
	fmt.Println(AddComma("123.66"))
	fmt.Println(AddComma("1234.66"))
	fmt.Println(AddComma("12345.66"))
	fmt.Println(AddComma("123456.66"))
	fmt.Println(AddComma("1234567.66"))
	fmt.Println(AddComma("12345678.66"))
	fmt.Println(AddComma("123456789.66"))
	fmt.Println(AddComma("1234567890.66"))
}