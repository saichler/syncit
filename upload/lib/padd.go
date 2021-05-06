package lib

import (
	"bytes"
	"errors"
)

//1.My objective was to keep the cyclomatic complexity low, we can discuss why this is important.
//2.Probably forgot some negative testing so code coverage will not be perfect, hope you have static analysis in place
//3.Coding standards are per GoLang

func Padd(str string, size int) (string, error) {
	//check for 0, return same string
	if size == 0 {
		return str, nil
	} else if size < 0 {
		//less than zero, return string + error
		return str, errors.New("Illegal size less than 0, returning same string")
	}
	//Hold final result buffer
	result := &bytes.Buffer{}
	//hold sequence of numeric chars
	var buff *bytes.Buffer
	//max number of zeros string to prefix
	maxZeros := maxZeros(size)
	dot := false

	//iterate over the chars
	for _, char := range str {
		if isDot(char) {
			//in case of a dot, always flush
			buff = flush(buff, result, maxZeros)
			dot = !dot
			result.WriteString(".")
		} else if !isNumeric(char) {
			//if not numeric, always flush
			buff = flush(buff, result, maxZeros)
			result.WriteByte(byte(char))
			dot = false
		} else {
			if dot {
				//if post the dot, just add to the final result
				result.WriteByte(byte(char))
			} else {
				if buff == nil {
					buff = &bytes.Buffer{}
				}
				//add to numeric buffer
				buff.WriteByte(byte(char))
			}
		}
	}

	//flush if ends with numbers
	flush(buff, result, maxZeros)

	return result.String(), nil
}

func flush(buff *bytes.Buffer, result *bytes.Buffer, maxZeros string) *bytes.Buffer {
	if buff != nil {
		//In case there is more numbers than the max prefix zeros
		if len(buff.Bytes()) < len(maxZeros) {
			result.WriteString(maxZeros[len(buff.Bytes()):])
		}
		result.WriteString(buff.String())
	}
	return nil
}

//returns true if the char is a dot
func isDot(char int32) bool {
	if char == 46 {
		return true
	}
	return false
}

//Generate a string with the max zeros available to pad by the size
func maxZeros(size int) string {
	buff := bytes.Buffer{}
	for i := 0; i < size; i++ {
		buff.WriteString("0")
	}
	return buff.String()
}

//Returns true if the char is numeric
func isNumeric(char int32) bool {
	if char >= 48 && char <= 57 {
		return true
	}
	return false
}
