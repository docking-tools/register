package template

import (
    "os"
    "testing"
)

func TestEnv(t *testing.T) {
	if err := os.Setenv("foo", "bar"); err != nil {
		t.Fatal(err)
	}

	result, err := env("foo")
	if err != nil {
		t.Fatal(err)
	}

	if result != "bar" {
		t.Errorf("expected %#v to be %#v", result, "bar")
	}
}