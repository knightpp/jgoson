package jgoson

import (
	"reflect"
	"sort"
)

func parseRecursive(value any, parentName string, depth int) *Type {
	if m, ok := value.(map[string]any); ok {
		t := &Type{
			Name:   parentName,
			Fields: nil,
		}

		for _, pair := range extractAndSortKV(m) {
			t.Fields = append(t.Fields, Field{
				Name: pair.K,
				Type: parseRecursive(pair.V, pair.K, depth+1),
			})
		}
		return t
	} else if s, ok := value.([]any); ok {
		if len(s) == 0 {
			return &Type{
				Name:    "any",
				IsSlice: true,
			}
		}

		newT := parseRecursive(s[0], parentName, depth+1)
		newT.IsSlice = true
		return newT
	} else {
		typ := reflect.TypeOf(value)
		if typ == nil {
			return &Type{
				Name: "any",
			}
		}
		return &Type{
			Name: typ.String(),
		}
	}
}

type kvPair struct {
	K string
	V any
}

func extractAndSortKV(m map[string]any) []kvPair {
	pairs := make([]kvPair, len(m))
	i := 0
	for k, v := range m {
		pairs[i] = kvPair{
			K: k,
			V: v,
		}
		i++
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].K < pairs[j].K
	})

	return pairs
}
