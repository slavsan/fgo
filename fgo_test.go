package fgo_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	f "github.com/slavsan/fgo"
	"github.com/slavsan/fgo/internal/assert"
)

func TestFilterStrings(t *testing.T) {
	testCases := []struct {
		title     string
		input     []string
		expected  []string
		predicate func(string) bool
	}{
		{
			title:     "can filter out elements with length bigger than 3",
			input:     []string{"foo", "bar", "spam", "eggs", "baz"},
			predicate: func(s string) bool { return len(s) > 3 },
			expected:  []string{"spam", "eggs"},
		},
		{
			title:     "can filter out elements with length bigger smaller than 4",
			input:     []string{"foo", "bar", "spam", "eggs", "baz"},
			predicate: func(s string) bool { return len(s) < 4 },
			expected:  []string{"foo", "bar", "baz"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[string, string, *strings.Builder](
					tc.input,
					f.Filter[string](tc.predicate),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFilterIntegers(t *testing.T) {
	testCases := []struct {
		title     string
		input     []int
		expected  []int
		predicate func(int) bool
	}{
		{
			title:     "can filter out elements with bigger than 10",
			input:     []int{8, 11, 7, 12, 1},
			predicate: func(i int) bool { return i > 10 },
			expected:  []int{11, 12},
		},
		{
			title:     "can filter out elements with length smaller than 10",
			input:     []int{8, 11, 7, 12, 1},
			predicate: func(i int) bool { return i < 10 },
			expected:  []int{8, 7, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[int, int, *strings.Builder](
					tc.input,
					f.Filter[int](tc.predicate),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapIntegerToInteger(t *testing.T) {
	testCases := []struct {
		title     string
		input     []int
		expected  []int
		transform func(i int) int
	}{
		{
			title:     "can increment integers",
			input:     []int{8, 11, 7, 12, 1},
			transform: func(i int) int { return i + 1 },
			expected:  []int{9, 12, 8, 13, 2},
		},
		{
			title:     "can decrement integers",
			input:     []int{8, 11, 7, 12, 1},
			transform: func(i int) int { return i - 1 },
			expected:  []int{7, 10, 6, 11, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[int, int, *strings.Builder](
					tc.input,
					f.Map[int, int](tc.transform),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapStringToString(t *testing.T) {
	testCases := []struct {
		title     string
		input     []string
		expected  []string
		transform func(i string) string
	}{
		{
			title:     "can increment integers",
			input:     []string{"foo", "bar", "baz"},
			transform: func(s string) string { return s + " x" },
			expected:  []string{"foo x", "bar x", "baz x"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[string, string, any](
					tc.input,
					f.Map[string, string](tc.transform),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapStringToInteger(t *testing.T) {
	testCases := []struct {
		title     string
		input     []string
		expected  []int
		transform func(i string) int
	}{
		{
			title:     "can increment integers",
			input:     []string{"foo", "bar", "baz", "spam", "eggs"},
			transform: func(s string) int { return len(s) },
			expected:  []int{3, 3, 3, 4, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[string, int, any](
					tc.input,
					f.Map[string, int](tc.transform),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapIntegerToString(t *testing.T) {
	testCases := []struct {
		title     string
		input     []int
		expected  []string
		transform func(i int) string
	}{
		{
			title:     "can increment integers",
			input:     []int{22, 52, 11, 2, 98},
			transform: func(i int) string { return strconv.Itoa(i) },
			expected:  []string{"22", "52", "11", "2", "98"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[int, string, any](
					tc.input,
					f.Map(tc.transform),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFilterAndMap(t *testing.T) {
	testCases := []struct {
		title     string
		input     []int
		expected  []string
		predicate func(i int) bool
		transform func(i int) string
	}{
		{
			title:     "can filter and map",
			input:     []int{223, 52, 111, 2, 98},
			predicate: func(i int) bool { return len(strconv.Itoa(i)) > 2 },
			transform: func(i int) string { return strconv.Itoa(i + 100) },
			expected:  []string{"323", "211"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[int, string, any](
					tc.input,
					f.Filter(tc.predicate),
					f.Map(tc.transform),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapAndFilter(t *testing.T) {
	testCases := []struct {
		title     string
		input     []int
		expected  []string
		transform func(i int) string
		predicate func(i string) bool
	}{
		{
			title:     "can filter and map",
			input:     []int{223, 52, 111, 2, 98},
			transform: func(i int) string { return strconv.Itoa(i) },
			predicate: func(s string) bool { return len(s) <= 2 },
			expected:  []string{"52", "2", "98"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[int, string, any](
					tc.input,
					f.Map(tc.transform),
					f.Filter(tc.predicate),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFilterMapReduce(t *testing.T) {
	testCases := []struct {
		title     string
		input     []string
		expected  []string
		reducedTo string
		transform func(i int) string
		predicate func(i string) bool
	}{
		{
			title:     "with strings builder as reducer",
			input:     []string{"foo", "bar", "spam", "eggs", "baz"},
			expected:  []string{"spam 4", "eggs 4"},
			reducedTo: "spam 4 eggs 4 ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, reduced, err :=
				f.Pipe[string, string, *strings.Builder](
					tc.input,
					f.Filter[string](func(s string) bool { return len(s) > 3 }),
					f.Map[string, string](func(s string) string { return fmt.Sprintf("%s %d", s, len(s)) }),
					f.Reduce[string, *strings.Builder](&strings.Builder{}, func(sb *strings.Builder, s string) *strings.Builder {
						sb.WriteString(s)
						sb.WriteString(" ")
						return sb
					}),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.reducedTo, reduced.String())
		})
	}
}

func TestMapFilterReduce(t *testing.T) {
	testCases := []struct {
		title     string
		input     []string
		expected  []string
		reducedTo string
		transform func(i int) string
		predicate func(i string) bool
	}{
		{
			title:     "with strings builder as reducer",
			input:     []string{"foo", "bar", "spam", "eggs", "baz"},
			expected:  []string{"foo 3", "bar 3", "baz 3"},
			reducedTo: "foo 3 bar 3 baz 3 ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, reduced, err :=
				f.Pipe[string, string, *strings.Builder](
					tc.input,
					f.Map[string, string](func(s string) string { return fmt.Sprintf("%s %d", s, len(s)) }),
					f.Filter[string](func(s string) bool { return strings.HasSuffix(s, "3") }),
					f.Reduce[string, *strings.Builder](&strings.Builder{}, func(sb *strings.Builder, s string) *strings.Builder {
						sb.WriteString(s)
						sb.WriteString(" ")
						return sb
					}),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.reducedTo, reduced.String())
		})
	}
}

func TestFilterStringMapToIntegerThenReduce(t *testing.T) {
	type sum struct {
		value int
	}

	testCases := []struct {
		title     string
		input     []string
		expected  []int
		reducedTo sum
		transform func(s string) int
		predicate func(i string) bool
	}{
		{
			title:     "with strings builder as reducer",
			input:     []string{"foo", "bar", "spam", "eggs", "baz"},
			expected:  []int{4, 4},
			reducedTo: sum{8},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, reducedTo, err :=
				f.Pipe[string, int, sum](
					tc.input,
					f.Filter[string](func(s string) bool { return len(s) == 4 }),
					f.Map[string, int](func(s string) int { return len(s) }),
					f.Reduce[int, sum](sum{}, func(acc sum, i int) sum {
						acc.value += i
						return acc
					}),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.reducedTo, reducedTo)
		})
	}
}

func TestFilterTakeStrings(t *testing.T) {
	testCases := []struct {
		title     string
		input     []string
		expected  []string
		predicate func(string) bool
		only      int
	}{
		{
			title:     "can filter out elements with length bigger than 3",
			input:     []string{"foo", "bar", "spam", "eggs", "baz"},
			predicate: func(s string) bool { return len(s) > 3 },
			only:      1,
			expected:  []string{"spam"},
		},
		{
			title:     "can filter out elements with length bigger smaller than 4",
			input:     []string{"foo", "bar", "spam", "eggs", "baz"},
			predicate: func(s string) bool { return len(s) < 4 },
			only:      4,
			expected:  []string{"foo", "bar", "baz"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[string, string, *strings.Builder](
					tc.input,
					f.Filter[string](tc.predicate),
					f.Take(tc.only),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFilterUniqueStrings(t *testing.T) {
	testCases := []struct {
		title     string
		input     []string
		expected  []string
		predicate func(string) bool
		unique    func(string) string
	}{
		{
			title:     "can filter out elements with length bigger than 3",
			input:     []string{"eggs", "foo", "spam", "bar", "spam", "eggs", "spam", "baz"},
			predicate: func(s string) bool { return len(s) > 3 },
			unique:    func(s string) string { return s },
			expected:  []string{"eggs", "spam"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[string, string, *strings.Builder](
					tc.input,
					f.Filter[string](tc.predicate),
					f.Unique(tc.unique),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFilterMapUnique(t *testing.T) {
	testCases := []struct {
		title     string
		input     []string
		expected  []int
		transform func(i int) string
		unique    func(int) string
		predicate func(i string) bool
	}{
		{
			title:    "with strings builder as reducer",
			input:    []string{"foo", "bar", "yo", "spam", "eggs", "baz"},
			unique:   func(i int) string { return strconv.Itoa(i) },
			expected: []int{3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result, _, err :=
				f.Pipe[string, int, *strings.Builder](
					tc.input,
					f.Filter[string](func(s string) bool { return len(s) > 2 }),
					f.Map[string, int](func(s string) int { return len(s) }),
					f.Unique[int](tc.unique),
				)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
