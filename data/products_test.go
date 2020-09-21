package data

import "testing"

// TestChecksValdiation unit test function
func TestChecksValdiation(t *testing.T) {
	p := &Product{
		Name:  "test",
		Price: 2.22,
		SKU:   "abasdf-cvcv-werr",
	}
	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}

}
