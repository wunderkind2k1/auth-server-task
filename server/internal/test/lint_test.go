package test

import (
	"fmt"
	"testing"
)

// TestLintIssues is a test with intentional linting issues
func TestLintIssues(t *testing.T) {
	fmt.Println("This is a test with linting issues")

	// should trigger error handling warning
	_ = fmt.Errorf("error without handling")
}
