package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInSlice(t *testing.T) {
	testSlice := []string{"a", "ba", "cba"}

	b, _ := InSlice("a", testSlice)
	require.True(t, b)

	b, _ = InSlice("ba", testSlice)
	require.True(t, b)

	b, _ = InSlice("cba", testSlice)
	require.True(t, b)

	b, _ = InSlice("b", testSlice)
	require.False(t, b)
}

func TestBlockSlideSlice(t *testing.T) {
	testSlice := []int{34, 56, 4, 13, 7, 4, 98, 3, 85, 0, 3, 4, 636, 8, 5}

	BlockSlideSlice(testSlice, 3, func(subSlice interface{}) {
		fmt.Printf("%+v\n", subSlice.([]int))
	})

	testSlice = []int{34, 56, 4, 13, 7, 4, 98, 3, 85, 0, 3, 4, 636, 8, 5}

	BlockSlideSlice(testSlice, 2, func(subSlice interface{}) {
		fmt.Printf("%+v\n", subSlice.([]int))
	})
}
