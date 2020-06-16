package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm displays a prompt `s` to the user and returns a bool indicating yes / no answer.
// Confirm limits the maximum amount of unrecognized inputs to 3, then exits with false.
func Confirm(s string) (bool, error) {
	r := bufio.NewReader(os.Stdin)

	for tries := 3; tries > 0; tries-- {
		fmt.Printf("%s [Y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			return false, err
		}

		// Empty input (i.e. "\n")
		res = strings.TrimSpace(res)
		if res == "" {
			return true, nil
		}

		if res = strings.ToLower(res); res == "yes" || res == "y" {
			return true, nil
		} else if res == "no" || res == "n" {
			return false, nil
		}
	}

	return false, nil
}
