package helper

import (
	"bufio"
	"fmt"
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