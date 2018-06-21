package vet_bug

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func AssertEqual(tb testing.TB, expected, actual int) {
	tb.Helper()

	// Pointless, but gives us an error in `go vet`
	config := aws.Config{}
	fmt.Println(config)

	if expected != actual {
		tb.Error("Expected:", expected, "; Got:", actual)
	}
}
