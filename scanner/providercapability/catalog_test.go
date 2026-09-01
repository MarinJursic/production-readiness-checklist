package providercapability

import (
	"slices"
	"testing"
)

func TestEmbeddedCapabilitiesAreValidAndStable(t *testing.T) {
	capabilities, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(capabilities) != 3 {
		t.Fatalf("capability count = %d", len(capabilities))
	}
	ids, err := IDs()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"prc.collect.prc.36.002.c1@0.1",
		"prc.collect.prc.36.004.c1@0.1",
		"prc.collect.prc.36.005.c1@0.1",
	}
	if !slices.Equal(ids, want) {
		t.Fatalf("unexpected capability IDs: %v", ids)
	}
}
