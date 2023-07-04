package main

import "testing"

func TestRecursion(t *testing.T) {
	t.Run("one field", func(t *testing.T) {
		recursion("", nil, 0)
	})
}
