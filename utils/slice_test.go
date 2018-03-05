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

	BlockSlideSlice(testSlice, 3, func(subSlice interface{}) bool {
		fmt.Printf("%+v\n", subSlice.([]int))
		return true
	})

	testSlice = []int{34, 56, 4, 13, 7, 4, 98, 3, 85, 0, 3, 4, 636, 8, 5}

	BlockSlideSlice(testSlice, 2, func(subSlice interface{}) bool {
		fmt.Printf("%+v\n", subSlice.([]int))
		return true
	})
}

func TestToInt64Slice(t *testing.T) {
	conv, err := ToInt64Slice([]string{"1234", "56325", "678", "2357", "3566", "-1"})
	require.NoError(t, err)

	require.EqualValues(t, []int64{1234, 56325, 678, 2357, 3566, -1}, conv)

}
