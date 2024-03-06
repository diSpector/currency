package client

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func ValidateCurrencyArgs(args []string) error {
	if len(args) == 0 {
		return errors.New("emptry currencies slice")
	}

	r := regexp.MustCompile(`[A-Z]{3}`)

	for i := range args {
		if len(args[i]) != 3 || !r.MatchString(strings.ToUpper(args[i])) {
			return errors.New("expect currency contains 3 latin letters, passed: " + args[i])
		}
	}

	return nil
}
