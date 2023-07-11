package jgoson_test

import (
	"testing"

	"github.com/knightpp/oleksii-brainfuck/jgoson"
	"github.com/stretchr/testify/assert"
)

func TestRecursion(t *testing.T) {
	t.Run("map with field", func(t *testing.T) {
		res := jgoson.Parse(map[string]any{
			"struct_name": 10,
		})

		assert.Equal(t, "Generated", res.Name)
		assert.Len(t, res.Fields, 1)
		assert.Equal(t, "struct_name", res.Fields[0].Name)
		assert.Equal(t, "int", res.Fields[0].Type.Name)
	})

	t.Run("map inside map with field", func(t *testing.T) {
		res := jgoson.Parse(map[string]any{
			"inner": map[string]any{
				"struct_name_inside": []any{32},
			},
		})

		assert.Equal(t, &jgoson.Type{
			Name: "Generated",
			Fields: []jgoson.Field{
				{
					Name: "inner",
					Type: &jgoson.Type{
						Name: "inner",
						Fields: []jgoson.Field{
							{
								Name: "struct_name_inside",
								Type: &jgoson.Type{
									Name:    "int",
									IsSlice: true,
								},
							},
						},
					},
				},
			},
		}, res)
	})

	t.Run("slice", func(t *testing.T) {
		res := jgoson.Parse(map[string]any{"tags": []any{"string"}})

		assert.Equal(t, &jgoson.Type{
			Name: "Generated",
			Fields: []jgoson.Field{
				{
					Name: "tags",
					Type: &jgoson.Type{
						Name:    "string",
						IsSlice: true,
					},
				},
			},
		}, res)
	})
}
