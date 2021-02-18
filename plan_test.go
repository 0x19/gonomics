package gonomics

import (
	"flag"
	"os"
	"testing"
)

var plan string

// TestMain used to give option to test only free plan Nomics API's.
// Enter "go test -plan=free -v" to run only free plan API's.
func TestMain(m *testing.M) {
	flag.StringVar(&plan, "plan", "paid", "if plan=free, only run free plan Nomics API's, otherwise run all")

	code := m.Run()
	os.Exit(code)
}
