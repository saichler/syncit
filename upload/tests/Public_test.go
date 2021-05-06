package tests

import (
	"fmt"
	"github.com/saichler/upload/lib"
	"strconv"
	"testing"
)

func Test_padd1(t *testing.T) {
	paddTest("James Bond 7", 3, "James Bond 007", t, false)
}

func Test_padd2(t *testing.T) {
	paddTest("PI=3.14", 2, "PI=03.14", t, false)
}
func Test_padd3(t *testing.T) {
	paddTest("It's 3:13pm", 2, "It's 03:13pm", t, false)
}
func Test_padd4(t *testing.T) {
	paddTest("99UR1337", 6, "000099UR001337", t, false)
}
func Test_padd5(t *testing.T) {
	paddTest("099UR01337", 6, "000099UR001337", t, false)
}
func Test_padd6(t *testing.T) {
	paddTest("099UR01337", 0, "099UR01337", t, false)
}
func Test_padd7(t *testing.T) {
	paddTest("099UR01337", -1, "099UR01337", t, true)
}
func Test_padd8(t *testing.T) {
	paddTest("NoNumbers", 6, "NoNumbers", t, false)
}
func Test_padd9(t *testing.T) {
	paddTest("Extra123456Extra", 5, "Extra123456Extra", t, false)
}
func Test_padd10(t *testing.T) {
	paddTest("Extra123..456Extra", 4, "Extra0123..0456Extra", t, false)
}

func paddTest(str string, size int, expected string, t *testing.T, errorExpected bool) {
	out, err := lib.Padd(str, size)
	if err != nil && !errorExpected {
		fmt.Println("Error in input:", err)
		t.Fail()
		return
	} else if err != nil {
		return
	}

	if out != expected {
		fmt.Println("Failed with Input:", str, "+", strconv.Itoa(size), "expected:'", expected, "' but got '", out, "'")
		t.Fail()
	}
}
