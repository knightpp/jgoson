package jgoson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"strings"
	"testing"

	"github.com/knightpp/jgoson"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]any
		expected *jgoson.Type
	}{
		{
			name:  "Nil",
			input: nil,
			expected: &jgoson.Type{
				Name:    "Generated",
				Fields:  nil,
				IsSlice: false,
			},
		},
		{
			name:  "Empty",
			input: map[string]any{},
			expected: &jgoson.Type{
				Name:    "Generated",
				Fields:  nil,
				IsSlice: false,
			},
		},
		{
			name: "Struct",
			input: map[string]any{
				"array": []any{
					map[string]any{
						"arrayField1": "Field1",
						"arrayField2": map[string]any{
							"innerField1": "Type1",
							"innerField2": nil,
							"innerField3": false,
						},
					},
				},
				"dictionary": map[string]any{
					"field1": 10,
					"field2": map[string]any{
						"field2A": 3.14,
						"field2B": []any{
							map[string]any{
								"field2BA": "NestedField",
								"field2BB": map[string]any{
									"field2BBA": nil,
									"field2BBB": false,
								},
							},
						},
						"field2C": false,
					},
				},
			},
			expected: &jgoson.Type{
				Name: "Generated",
				Fields: []jgoson.Field{
					{
						Name: "array",
						Type: &jgoson.Type{
							Name: "array",
							Fields: []jgoson.Field{
								{
									Name: "arrayField1", Type: &jgoson.Type{Name: "string"},
								},
								{
									Name: "arrayField2", Type: &jgoson.Type{
										Name: "arrayField2",
										Fields: []jgoson.Field{
											{
												Name: "innerField1", Type: &jgoson.Type{Name: "string"},
											},
											{
												Name: "innerField2", Type: &jgoson.Type{Name: "any"},
											},
											{
												Name: "innerField3", Type: &jgoson.Type{Name: "bool"},
											},
										},
									},
								},
							},
							IsSlice: true,
						},
					},
					{
						Name: "dictionary",
						Type: &jgoson.Type{
							Name: "dictionary",
							Fields: []jgoson.Field{
								{
									Name: "field1", Type: &jgoson.Type{Name: "int"},
								},
								{
									Name: "field2", Type: &jgoson.Type{
										Name: "field2",
										Fields: []jgoson.Field{
											{
												Name: "field2A", Type: &jgoson.Type{Name: "float64"},
											},
											{
												Name: "field2B", Type: &jgoson.Type{
													Name: "field2B",
													Fields: []jgoson.Field{
														{
															Name: "field2BA", Type: &jgoson.Type{Name: "string"},
														},
														{
															Name: "field2BB", Type: &jgoson.Type{
																Name: "field2BB",
																Fields: []jgoson.Field{
																	{
																		Name: "field2BBA", Type: &jgoson.Type{Name: "any"},
																	},
																	{
																		Name: "field2BBB", Type: &jgoson.Type{Name: "bool"},
																	},
																},
															},
														},
													},
													IsSlice: true,
												},
											},
											{
												Name: "field2C", Type: &jgoson.Type{Name: "bool"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Primitive",
			input: map[string]any{
				"name": "a",
				"fields": []any{
					map[string]any{
						"name": "Field",
						"type": "int",
					},
				},
			},
			expected: &jgoson.Type{
				Name: "Generated",
				Fields: []jgoson.Field{
					{
						Name: "fields",
						Type: &jgoson.Type{
							Name: "fields",
							Fields: []jgoson.Field{
								{
									Name: "name",
									Type: &jgoson.Type{Name: "string"},
								},
								{
									Name: "type",
									Type: &jgoson.Type{Name: "string"},
								},
							},
							IsSlice: true,
						},
					},
					{
						Name: "name",
						Type: &jgoson.Type{
							Name: "string",
						},
					},
				},
				IsSlice: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := jgoson.Parse(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseJSON(t *testing.T) {
	testCases := []struct {
		name            string
		jsonString      string
		expected        *jgoson.Type
		expectedErrType error
	}{
		{
			name: "Valid JSON",
			jsonString: `{
				"field1": "string",
				"field2": 10,	
				"field3": 3.14,
				"field4": true,
				"field5": null
			}`,
			expected: &jgoson.Type{
				Name: "Generated",
				Fields: []jgoson.Field{
					{
						Name: "field1",
						Type: &jgoson.Type{Name: "string"},
					},
					{
						Name: "field2",
						Type: &jgoson.Type{Name: "float64"},
					},
					{
						Name: "field3",
						Type: &jgoson.Type{Name: "float64"},
					},
					{
						Name: "field4",
						Type: &jgoson.Type{Name: "bool"},
					},
					{
						Name: "field5",
						Type: &jgoson.Type{Name: "any"},
					},
				},
				IsSlice: false,
			},
			expectedErrType: nil,
		},
		{
			name:            "Invalid JSON",
			jsonString:      "{invalid",
			expected:        nil,
			expectedErrType: &json.SyntaxError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.jsonString)
			result, err := jgoson.ParseJSON(reader)

			if tc.expectedErrType != nil {
				assert.ErrorAs(t, err, &tc.expectedErrType)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestType_ToGoInline(t *testing.T) {
	testCases := []struct {
		name            string
		typ             jgoson.Type
		expectedOutput  string
		expectedTagName string
		tagOpts         []string
	}{
		{
			name:            "nil fields",
			typ:             jgoson.Type{Name: "Generated"},
			expectedOutput:  "type Generated struct {\n}\n",
			expectedTagName: "tag",
		},
		{
			name: "dictionary",
			typ: jgoson.Type{
				Name: "Main",
				Fields: []jgoson.Field{
					{
						Name: "field1",
						Type: &jgoson.Type{Name: "int"},
					},
					{
						Name: "field2",
						Type: &jgoson.Type{Name: "float64"},
					},
					{
						Name: "field3",
						Type: &jgoson.Type{Name: "bool"},
					},
				},
			},
			expectedOutput:  "type Main struct {\n\tField1 int     `json:\"field1\"`\n\tField2 float64 `json:\"field2\"`\n\tField3 bool    `json:\"field3\"`\n}\n",
			expectedTagName: "json",
			tagOpts:         []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cfg := jgoson.Config{
				Tag:     tc.expectedTagName,
				TagOpts: tc.tagOpts,
			}
			tc.typ.ToGoInline(&buf, cfg)

			result, err := format.Source(buf.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, string(result))
		})
	}
}

func TestType_ToGo(t *testing.T) {
	testCases := []struct {
		name           string
		typ            jgoson.Type
		expectedOutput string
		expectedTag    string
	}{
		{
			name: "Empty Type",
			typ: jgoson.Type{
				Name:    "Generated",
				Fields:  nil,
				IsSlice: false,
			},
			expectedOutput: "type Generated struct {\n}\n",
			expectedTag:    "json",
		},
		{
			name: "Struct Type",
			typ: jgoson.Type{
				Name: "Generated",
				Fields: []jgoson.Field{
					{
						Name: "Field1",
						Type: &jgoson.Type{
							Name:    "Type1",
							Fields:  nil,
							IsSlice: false,
						},
					},
					{
						Name: "Field2",
						Type: &jgoson.Type{
							Name: "Type2",
							Fields: []jgoson.Field{
								{
									Name: "NestedField",
									Type: &jgoson.Type{
										Name:    "NestedType",
										Fields:  nil,
										IsSlice: false,
									},
								},
							},
							IsSlice: false,
						},
					},
				},
				IsSlice: false,
			},
			expectedOutput: "type Generated struct {\n\tField1 Type1 `json:\"field1\"`\n\tField2 Type2 `json:\"field2\"`\n}\n\n" +
				"type Type1 struct {\n}\n\n" +
				"type Type2 struct {\n\tNestedField NestedType `json:\"nested_field\"`\n}\n\n" +
				"type NestedType struct {\n}\n",
			expectedTag: "json",
		},
		{
			name: "Slice Type",
			typ: jgoson.Type{
				Name: "Generated",
				Fields: []jgoson.Field{
					{
						Name: "Field",
						Type: &jgoson.Type{
							Name:    "SliceType",
							Fields:  nil,
							IsSlice: false,
						},
					},
				},
				IsSlice: true,
			},
			expectedOutput: "type Generated struct {\n\tField []SliceType `json:\"field\"`\n}\n\n" +
				"type SliceType struct {\n}\n",
			expectedTag: "json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cfg := jgoson.Config{}
			tc.typ.ToGo(&buf, cfg)

			result := buf.String()
			assert.Equal(t, tc.expectedOutput, result)

			assert.Contains(t, result, fmt.Sprintf("`%s", tc.expectedTag))
		})
	}
}
