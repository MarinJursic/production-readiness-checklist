package providercapability

import "testing"

func TestEmbeddedCapabilitiesAreValidAndStable(t *testing.T) {
	capabilities, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(capabilities) != 1 {
		t.Fatalf("capability count = %d", len(capabilities))
	}
	ids, err := IDs()
	if err != nil {
		t.Fatal(err)
	}
	if ids[0] != "prc.collect.prc.36.004.c1@0.1" {
		t.Fatalf("unexpected capability IDs: %v", ids)
	}
}
