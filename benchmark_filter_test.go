package fgo_test

import (
	"fmt"
	"strconv"
	"testing"

	f "github.com/slavsan/fgo"
	"github.com/slavsan/fgo/internal/assert"
)

func imperativeFilter(input []string) []string {
	output := make([]string, 0, len(input))
	for _, v := range input {
		if len(v) > 10 {
			output = append(output, v)
		}
	}
	return output
}

func funFilter(t *testing.T, input []string) []string {
	output, _, err :=
		f.Pipe[string, string, any](
			input,
			f.Filter(func(s string) bool { return len(s) > 10 }),
		)
	if t != nil {
		assert.Nil(t, err)
	}

	return output
}

func TestFilterImplementations(t *testing.T) {
	input := []string{"foo bar baz", "foo", "baz spam eggs"}
	expected := []string{"foo bar baz", "baz spam eggs"}
	assert.Equal(t, expected, imperativeFilter(input))
	assert.Equal(t, expected, funFilter(t, input))
}

func BenchmarkImperativeFilter(b *testing.B) {
	testCases := []int{10, 100, 1000, 10000}
	for _, tc := range testCases {
		tc := tc
		b.Run(fmt.Sprintf("with %d items", tc), func(b *testing.B) {
			input := generateStrings(tc)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = imperativeFilter(input)
			}
		})
	}
}

func BenchmarkFunctionalFilter(b *testing.B) {
	testCases := []int{10, 100, 1000, 10000}
	for _, tc := range testCases {
		tc := tc
		b.Run(fmt.Sprintf("with %d items", tc), func(b *testing.B) {
			input := generateStrings(tc)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = funFilter(nil, input)
			}
		})
	}
}

func generateStrings(number int) []string {
	output := make([]string, 0, number)
	for i := 0; i < number; i++ {
		output = append(output, strconv.Itoa(i))
	}
	return output
}
