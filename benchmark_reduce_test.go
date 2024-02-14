package fgo_test

import (
	"fmt"
	"testing"

	f "github.com/slavsan/fgo"
	"github.com/slavsan/fgo/internal/assert"
)

func imperativeReduce(input []string) int {
	var sum int
	for _, v := range input {
		sum += len(v)
	}
	return sum
}

func funReduce(t *testing.T, input []string) int {
	_, sum, err :=
		f.Pipe[string, int, int](
			input,
			f.Reduce[string, int](0, func(acc int, s string) int { return acc + len(s) }),
		)
	if t != nil {
		assert.Nil(t, err)
	}

	return sum
}

func TestReduceImplementations(t *testing.T) {
	input := []string{"foo bar baz", "foo", "baz spam eggs"}
	expected := 27
	assert.Equal(t, expected, imperativeReduce(input))
	assert.Equal(t, expected, funReduce(t, input))
}

func BenchmarkImperativeReduce(b *testing.B) {
	testCases := []int{10, 100, 1000, 10000}
	for _, tc := range testCases {
		tc := tc
		b.Run(fmt.Sprintf("with %d items", tc), func(b *testing.B) {
			input := generateStrings(tc)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = imperativeReduce(input)
			}
		})
	}
}

func BenchmarkFunctionalReduce(b *testing.B) {
	testCases := []int{10, 100, 1000, 10000}
	for _, tc := range testCases {
		tc := tc
		b.Run(fmt.Sprintf("with %d items", tc), func(b *testing.B) {
			input := generateStrings(tc)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = funReduce(nil, input)
			}
		})
	}
}
