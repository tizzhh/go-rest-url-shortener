package random

import "testing"

func TestNewRandomString(t *testing.T) {
	testCases := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 50",
			size: 50,
		},
		{
			name: "size = 100",
			size: 100,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			str1 := NewRandomString(tt.size)
			str2 := NewRandomString(tt.size)

			if len(str1) != tt.size {
				t.Errorf("got %d, want %d", len(str1), tt.size)
			}
			if len(str2) != tt.size {
				t.Errorf("got %d, want %d", len(str2), tt.size)
			}

			if str1 == str2 {
				t.Errorf("got same strings: %s %s", str1, str2)
			}
		})
	}
}
