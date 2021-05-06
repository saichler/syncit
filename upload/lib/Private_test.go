package lib

import (
	"fmt"
	"testing"
)

func Test_isNumeric(t *testing.T) {
	str := "0123456789"
	for _, char := range str {
		if !isNumeric(char) {
			t.Fail()
			fmt.Println("Failed! Expected char ", string(char), " to be numeric")
		}
	}
	str = "abcdAb s$#@^#$^%"
	for _, char := range str {
		if isNumeric(char) {
			t.Fail()
			fmt.Println("Failed! Expected char ", string(char), " to be none numeric")
		}
	}
}

func Test_maxZeros(t *testing.T) {
	str := "000000"
	maxZeros := maxZeros(6)
	if str != maxZeros {
		t.Fail()
		fmt.Println("Failed!")
	}
}

func Test_isDot(t *testing.T) {
	str := "."
	for _, char := range str {
		if !isDot(char) {
			t.Fail()
			fmt.Println("Expected only dots")
			return
		}
	}
}
