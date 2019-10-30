package testflags

import (
	"flag"
	"testing"
)

var unitTest = flag.Bool("unit", true, "Run unit tests")
var integrationTest = flag.Bool("integration", true, "Run integration tests")

func IntegrationTest(t *testing.T) {
	if !*integrationTest {
		t.SkipNow()
	}

	t.Parallel()
}

func UnitTest(t *testing.T) {
	if !*unitTest && !testing.Short() {
		t.SkipNow()
	}

	t.Parallel()
}
