package engine

import (
	"errors"
	"fmt"
)

type FibonacciEngine struct{}

func NewFibonacciEngine() *FibonacciEngine {
	return &FibonacciEngine{}
}

// ErrInvalidCount identifies a count outside the supported range.
var ErrInvalidCount = errors.New("invalid Fibonacci count")

const (
	// MaxFibonacciCount is 94 terms, from F(0) through F(93), that fit in uint64.
	MaxFibonacciCount = 94
)

// CalculateFibonacciSequence generates a sequence of the first n Fibonacci numbers.
func (f *FibonacciEngine) CalculateFibonacciSequence(n int) ([]uint64, error) {
	if n < 0 {
		return []uint64{}, fmt.Errorf("%w: input must be a non-negative integer", ErrInvalidCount)
	}
	if n > MaxFibonacciCount {
		return []uint64{}, fmt.Errorf("%w: input exceeds maximum allowed count of %d", ErrInvalidCount, MaxFibonacciCount)
	}

	sequence := make([]uint64, n)

	if n == 0 {
		return sequence, nil
	}

	sequence[0] = 0
	if n == 1 {
		return sequence, nil
	}
	sequence[1] = 1

	// Standard linear progression: F(n) = F(n-1) + F(n-2)
	for i := 2; i < n; i++ {
		sequence[i] = sequence[i-1] + sequence[i-2]
	}

	return sequence, nil
}
