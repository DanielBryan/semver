package semver

import (
	"testing"
)

// Test string representation of versions
func TestString(t *testing.T) {
	tests := []struct {
		v Version
		s string
	}{
		{Version{1, 2, 3, ""}, "v1.2.3"},
		{Version{1, 2, 3, "beta"}, "v1.2.3-beta"},
		{Version{0, 2, 3, "beta"}, "v0.2.3-beta"},
	}

	for _, tc := range tests {
		if tc.v.String() != tc.s {
			t.Errorf("String representation of %#v did not match what was expected: '%s'", tc.v, tc.s)
			t.Fail()
		}
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		a     Version
		b     Version
		equal bool
	}{
		{Version{1, 2, 3, ""}, Version{1, 2, 3, ""}, true},
		{Version{1, 2, 3, "beta"}, Version{1, 2, 3, "beta"}, true},
		{Version{0, 0, 1, "foo"}, Version{0, 0, 1, "foo"}, true},

		{Version{1, 2, 3, "beta"}, Version{0, 2, 3, "beta"}, false},
		{Version{1, 2, 3, "beta"}, Version{1, 0, 3, "beta"}, false},
		{Version{1, 2, 3, "beta"}, Version{1, 2, 0, "beta"}, false},
		{Version{1, 2, 3, "beta"}, Version{1, 2, 3, ""}, false},
	}

	for _, tc := range tests {
		eq := tc.a.Equals(tc.b)
		if eq && !tc.equal {
			t.Errorf("Expected these not be equal: %v, %v", tc.a, tc.b)
			t.Fail()
		} else if !eq && tc.equal {
			t.Errorf("Expected these to be equal: %v, %v", tc.a, tc.b)
			t.Fail()
		}
	}
}

func TestGreaterThan(t *testing.T) {
	tests := []struct {
		a       Version
		b       Version
		greater bool
	}{

		// equal
		{Version{1, 2, 3, ""}, Version{1, 2, 3, ""}, false},
		{Version{1, 2, 3, "beta"}, Version{1, 2, 3, "beta"}, false},

		// greater
		{Version{2, 2, 3, "beta"}, Version{1, 2, 3, "beta"}, true},
		{Version{1, 3, 3, "beta"}, Version{1, 2, 3, "beta"}, true},
		{Version{1, 2, 4, "beta"}, Version{1, 2, 3, "beta"}, true},
		{Version{1, 2, 3, "c"}, Version{1, 2, 3, "beta"}, true},
		{Version{1, 2, 3, "z"}, Version{1, 2, 3, "foo"}, true},
		{Version{Major: 4}, Version{3, 6, 1, ""}, true},

		// lesser
		{Version{0, 0, 1, "foo"}, Version{0, 0, 1, "z"}, false},
		{Version{0, 0, 1, "foo"}, Version{0, 0, 2, "foo"}, false},
		{Version{3, 6, 1, ""}, Version{Major: 4}, false},
	}

	for _, tc := range tests {
		eq := tc.a.GreaterThan(tc.b)
		if eq && !tc.greater {
			t.Errorf("Expected %v to be greater than %v", tc.a, tc.b)
			t.Fail()
		} else if !eq && tc.greater {
			t.Errorf("Expected %v NOT to be greater than %v", tc.a, tc.b)
			t.Fail()
		}
	}
}

func TestLessThan(t *testing.T) {
	tests := []struct {
		a    Version
		b    Version
		less bool
	}{

		// equal
		{Version{1, 2, 3, ""}, Version{1, 2, 3, ""}, false},
		{Version{1, 2, 3, "beta"}, Version{1, 2, 3, "beta"}, false},

		// greater
		{Version{2, 2, 3, "beta"}, Version{1, 2, 3, "beta"}, false},
		{Version{1, 3, 3, "beta"}, Version{1, 2, 3, "beta"}, false},
		{Version{1, 2, 4, "beta"}, Version{1, 2, 3, "beta"}, false},
		{Version{1, 2, 3, "c"}, Version{1, 2, 3, "beta"}, false},
		{Version{1, 2, 3, "z"}, Version{1, 2, 3, "foo"}, false},
		{Version{Major: 4}, Version{3, 6, 1, ""}, false},

		// lesser
		{Version{0, 0, 1, "foo"}, Version{0, 0, 1, "z"}, true},
		{Version{0, 0, 1, "foo"}, Version{0, 0, 2, "foo"}, true},
		{Version{3, 6, 1, ""}, Version{Major: 4}, true},
	}

	for _, tc := range tests {
		eq := tc.a.LessThan(tc.b)
		if eq && !tc.less {
			t.Errorf("Expected %v to be less than %v", tc.a, tc.b)
			t.Fail()
		} else if !eq && tc.less {
			t.Errorf("Expected %v NOT to be less than %v", tc.a, tc.b)
			t.Fail()
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"v0", true},
		{"v1", true},
		{"v1.2", true},
		{"v1.2.3", true},
		{"v1.2.3-rc.2", true},
		{"v2.0.0-beta", true},

		{"", false},
		{"v", false},
		{"v.", false},
		{"v..", false},

		{"1", false},
		{"1.2", false},

		// TODO(dbryan) - should these be valid or invalid?
		//{"v1.", true},
		//{"v1.2.", true},

		{"v1.lol", false},
		{"v1.2.haha", false},

		{"1.2.3", false},
		{"1.2.3haha", false},
		{"v-", false},
	}

	for _, tc := range tests {
		_, err := Parse(tc.input)
		if err != nil && tc.valid {
			t.Errorf("Unexpected error for %s: %v", tc.input, err)
			t.Fail()
		} else if err == nil && !tc.valid {
			t.Errorf("This string should cause an error: %s", tc.input)
			t.Fail()
		}
	}
}
