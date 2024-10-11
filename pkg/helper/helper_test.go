package helper

import (
	"fmt"
	"testing"
)

func TestHelperOffsets(t *testing.T) {

	h, _ := NewHelper("./test.rpp")

	output := h.Parse()

	fmt.Printf("%+v", output)

	if len(output.Sections) < 1 {
		t.Fatalf("testing faild: expected more then one section. got %d", len(output.Sections))
	}

}
