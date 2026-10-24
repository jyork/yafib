package engine

import (
	"errors"
	"reflect"
	"testing"
)

func TestCalculateFibonacciSequence(t *testing.T) {
	tests := []struct {
		name      string
		input     int
		wantError bool
		want      []uint64
	}{
		{
			name:      "Negative value input",
			input:     -5,
			wantError: true,
			want:      []uint64{},
		},
		{
			name:      "Zero count input",
			input:     0,
			wantError: false,
			want:      []uint64{},
		},
		{
			name:  "Single term base case",
			input: 1,
			want:  []uint64{0},
		},
		{
			name:  "Two terms base case",
			input: 2,
			want:  []uint64{0, 1},
		},
		{
			name:  "Standard mini-sequence (n=5)",
			input: 5,
			want:  []uint64{0, 1, 1, 2, 3},
		},
		{
			name:  "Ten terms sequence",
			input: 10,
			want:  []uint64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34},
		},
		{
			name:      "Exceeds maximum allowed count",
			input:     95,
			wantError: true,
			want:      []uint64{},
		},
	}

	for _, tt := range tests {
		f := NewFibonacciEngine()
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.CalculateFibonacciSequence(tt.input)
			if tt.wantError && !errors.Is(err, ErrInvalidCount) {
				t.Fatalf("error = %v, want ErrInvalidCount", err)
			}
			if (err != nil) != tt.wantError {
				t.Fatalf("CalculateFibonacciSequence(%d) error = %v, wantError = %v",
					tt.input, err, tt.wantError)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateFibonacciSequence(%d) =\nGot : %v\nWant: %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestCalculateFibonacciSequenceMaxCount validates the largest supported sequence.
func TestCalculateFibonacciSequenceMaxCount(t *testing.T) {
	n := 94
	f := NewFibonacciEngine()
	seq, err := f.CalculateFibonacciSequence(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(seq) != n {
		t.Fatalf("Expected slice length of %d, got %d", n, len(seq))
	}

	// The 94th element in our 0-indexed array is F(93).
	// F(93) is the absolute highest Fibonacci number that fits into a uint64.
	var expectedLastElement uint64 = 12200160415121876738
	gotLastElement := seq[n-1]

	if gotLastElement != expectedLastElement {
		t.Errorf("Unexpected final value for maximum count.\nGot : %d\nWant: %d", gotLastElement, expectedLastElement)
	}
}
