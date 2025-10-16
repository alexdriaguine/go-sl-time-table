package utils_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/alexdriaguine/go-sl-time-table/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {

	t.Run("can map ints", func(t *testing.T) {
		list := []int{1, 2, 3, 4, 5}
		want := []int{2, 4, 6, 8, 10}
		got := utils.Map(list, func(i int) int {
			return i * 2
		})

		assert.Equal(t, want, got)
	})

	t.Run("can map strings to int", func(t *testing.T) {
		list := []string{"1", "2", "3", "4", "5"}
		want := []int{1, 2, 3, 4, 5}

		mapFunc := func(s string) int {
			v, _ := strconv.Atoi(s)
			return v
		}

		got := utils.Map(list, mapFunc)

		assert.Equal(t, want, got)
	})

	t.Run("maps structs", func(t *testing.T) {
		type A struct {
			FirstName string
			LastName  string
			Age       int
		}

		type B struct {
			FullName string
		}

		list := []A{{"Alex", "Driaguine", 33}, {"Oscar", "Tjena", 19}}
		want := []B{{"Alex Driaguine"}, {"Oscar Tjena"}}

		mapFunc := func(a A) B {
			return B{fmt.Sprintf("%s %s", a.FirstName, a.LastName)}
		}

		got := utils.Map(list, mapFunc)
		assert.Equal(t, want, got)
	})
}

func TestFilter(t *testing.T) {
	t.Run("can filter ints", func(t *testing.T) {
		list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		want := []int{2, 4, 6, 8, 10}

		filterFunc := func(i int) bool {
			return i%2 == 0
		}
		got := utils.Filter(list, filterFunc)

		assert.Equal(t, got, want)
	})

	t.Run("can filter strings", func(t *testing.T) {
		list := []string{"alex", "alexandre", "papacito", "papi chulo", "papi"}
		want := []string{"alexandre", "papacito", "papi chulo"}

		filterFunc := func(s string) bool {
			return len(s) > 4
		}
		got := utils.Filter(list, filterFunc)

		assert.Equal(t, got, want)
	})

	t.Run("can filter structs", func(t *testing.T) {
		type A struct {
			Name string
			Age  int
		}

		list := []A{{"Alex", 33}, {"Oscar", 19}, {"Robin", 25}}

		mapFunc := func(a A) bool {
			return a.Age > 20
		}

		got := utils.Filter(list, mapFunc)
		want := []A{{"Alex", 33}, {"Robin", 25}}
		assert.Equal(t, want, got)
	})
}
