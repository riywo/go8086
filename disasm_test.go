package go8086

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
)

type disasmTest struct {
	address uint16
	bytes   Bytes
	out     string
	line    string
}

func disasmTests() []disasmTest {
	file, err := os.Open("test/data.s")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tests := []disasmTest{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Fields(scanner.Text())
		address, _ := strconv.ParseUint(s[0], 16, 16)
		in := hexStrToBytes(s[1])
		t := disasmTest{uint16(address), in, strings.Join(s[2:], " "), line}
		tests = append(tests, t)
	}
	return tests
}

func hexStrToBytes(s string) Bytes {
	if len(s) == 0 {
		return Bytes{}
	} else {
		i, _ := strconv.ParseInt(s[0:2], 16, 16)
		b := byte(i)
		return append(Bytes{b}, hexStrToBytes(s[2:])...)
	}
}

func TestDisasm(t *testing.T) {
	for _, test := range disasmTests() {
		op := getOpcode(nil, test.address, test.bytes)
		assert.Equal(t, test.out, op.Disasm())
		assert.Equal(t, test.bytes, op.bytes)
	}
}
