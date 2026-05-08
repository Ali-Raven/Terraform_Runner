package helper

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"github.com/TwiN/go-color"
)

// ========================================================================= ReadRequired ==================================================================
func ReadRequired(reader *bufio.Reader, label string) string {
	for {
		fmt.Print(label)
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)
		fmt.Println("--------")

		if value != "" {
			return value
		}

		fmt.Println(color.Red + "This field is required." + color.Reset)
	}
}

// ========================================================================= Atoi ==================================================================
// for simple string to int conversion
func Atoi(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

// ========================================================================= ReadLine ==================================================================
func ReadLine(reader *bufio.Reader) string {
	s, _ := reader.ReadString('\n')
	return strings.TrimSpace(s)
}

// ========================================================================= ReadInt ==================================================================
func ReadInt(reader *bufio.Reader) (int, error) {
	var input []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		if b == '\n' || b == ' ' || b == '\r' {
			break
		}
		input = append(input, b)
	}

	// converting bytes to strings => then to integer
	numStr := string(input)
	num, err := strconv.Atoi(strings.TrimSpace(numStr))
	if err != nil {
		return 0, fmt.Errorf("%sinvalid Integer (You enter => %s) %s", color.Red, numStr, color.Reset)
	}

	return num, nil
}