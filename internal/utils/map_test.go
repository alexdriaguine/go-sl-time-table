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
		got := utils.Map(list, func(s string) int {
			v, _ := strconv.Atoi(s)
			return v
		})

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

		got := utils.Map(list, func(a A) B { return B{fmt.Sprintf("%s %s", a.FirstName, a.LastName)} })
		assert.Equal(t, want, got)
	})
}
