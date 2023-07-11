package jgoson_test

import (
	"fmt"
	"testing"

	"github.com/knightpp/jgoson"
)

func TestConfig_FillDefaults(t *testing.T) {
	t.Run("Default Tag", func(t *testing.T) {
		config := &jgoson.Config{}
		config.FillDefaults()

		expectedTag := "json"
		if config.Tag != expectedTag {
			t.Errorf("Expected Tag to be %s, but got %s", expectedTag, config.Tag)
		}
	})

	t.Run("Default TagNameFn", func(t *testing.T) {
		config := &jgoson.Config{}
		config.FillDefaults()

		expectedFn := jgoson.LowerCamelCaseToSnakeCase
		if !isFuncEqual(expectedFn, config.TagNameFn) {
			t.Errorf("Expected TagNameFn to be %p, but got %p", expectedFn, config.TagNameFn)
		}
	})

	t.Run("Default TagOpts", func(t *testing.T) {
		config := &jgoson.Config{}
		config.FillDefaults()

		expectedOpts := []string{"omitempty"}
		if len(config.TagOpts) != len(expectedOpts) {
			t.Errorf("Expected TagOpts length to be %d, but got %d", len(expectedOpts), len(config.TagOpts))
		}

		for i, opt := range expectedOpts {
			if config.TagOpts[i] != opt {
				t.Errorf("Expected TagOpts[%d] to be %s, but got %s", i, opt, config.TagOpts[i])
			}
		}
	})

	t.Run("Default StructNameFn", func(t *testing.T) {
		config := &jgoson.Config{}
		config.FillDefaults()

		expectedFn := jgoson.SnakeCaseToUpperCamelCase
		if !isFuncEqual(config.StructNameFn, expectedFn) {
			t.Errorf("Expected StructNameFn to be %p, but got %p", expectedFn, config.StructNameFn)
		}
	})

	t.Run("Default StructFieldFn", func(t *testing.T) {
		config := &jgoson.Config{}
		config.FillDefaults()

		expectedFn := jgoson.SnakeCaseToUpperCamelCase
		if !isFuncEqual(config.StructFieldFn, expectedFn) {
			t.Errorf("Expected StructFieldFn to be %p, but got %p", expectedFn, config.StructFieldFn)
		}
	})
}

func TestLowerCamelCaseToSnakeCase(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty String", "", ""},
		{"No Uppercase", "test", "test"},
		{"Single Uppercase", "testCase", "test_case"},
		{"Multiple Uppercase", "helloWorldGo", "hello_world_go"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := jgoson.LowerCamelCaseToSnakeCase(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, result)
			}
		})
	}
}

func TestSnakeCaseToUpperCamelCase(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty String", "", ""},
		{"No Underscore", "test", "Test"},
		{"Single Underscore", "test_case", "TestCase"},
		{"Multiple Underscores", "__hello__world___", "HelloWorld"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := jgoson.SnakeCaseToUpperCamelCase(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, result)
			}
		})
	}
}

func isFuncEqual(a, b func(string) string) bool {
	return fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b)
}
