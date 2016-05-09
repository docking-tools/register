package template

import (
    "os"
)

// env returns the value of the environment variable set
func env(s string) (string, error) {
	return os.Getenv(s), nil
}