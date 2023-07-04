package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecursion(t *testing.T) {
	t.Run("one field", func(t *testing.T) {
		res := recursion("")

		assert.Equal(t, "string", res.Name)
		assert.Nil(t, res.Fields)
	})

	t.Run("map with field", func(t *testing.T) {
		res := recursion(map[string]any{
			"struct_name": 10,
		})

		assert.Equal(t, "Generated", res.Name)
		assert.Len(t, res.Fields, 1)
		assert.Equal(t, "struct_name", res.Fields[0].Name)
		assert.Equal(t, "int", res.Fields[0].Type.Name)
	})

	t.Run("map inside map with field", func(t *testing.T) {
		res := recursion(map[string]any{
			"inner": map[string]any{
				"struct_name_inside": []int{},
			},
		})

		assert.Equal(t, "Generated", res.Name)
		assert.Len(t, res.Fields, 1)

		innerField := res.Fields[0]
		assert.Equal(t, "inner", innerField.Name)

		innerType := innerField.Type
		assert.Equal(t, "inner", innerType.Name)
		assert.Len(t, innerType.Fields, 1)

		finalField := innerType.Fields[0]
		assert.Equal(t, "struct_name_inside", finalField.Name)
		assert.Equal(t, "[]int", finalField.Type.Name)
		assert.Len(t, finalField.Type.Fields, 0)
	})
}
