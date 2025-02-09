package commands

import (
	"reflect"
	"testing"
)

func TestTrimSides(t *testing.T) {
	cases := []struct {
		lines          []string
		pixelsToRemove int
		expected       []string
	}{
		{
			lines:          []string{"# : # @ : ", "# % # # # ", "# : @ : : "},
			pixelsToRemove: 3,
			expected:       []string{": # ", "% # ", ": @ "},
		},
		{
			lines:          []string{"# : # @ : ", "# % # $ # ", "# : @ : * "},
			pixelsToRemove: 2,
			expected:       []string{": # @ ", "% # $ ", ": @ : "},
		},
	}

	for _, c := range cases {
		result := trimSides(c.lines, c.pixelsToRemove)
		if !reflect.DeepEqual(result, c.expected) {
			t.Errorf("expected %v, got %v", c.expected, result)
		}
	}
}
