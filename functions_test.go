package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	type test struct {
		name  string
		input string
		want  []string
	}
	tests := []test{
		{name: "basic test", input: "just some words", want: []string{"just", "some", "words"}},
		{name: "capital letters", input: "Hello There", want: []string{"hello", "there"}},
		{name: "trailing punctuation marks", input: "Whats YOUR NAME?!", want: []string{"whats", "your", "name?!"}},
		{name: "trailing spaces", input: "this string has spaces  ", want: []string{"this", "string", "has", "spaces"}},
		{name: "included spaces", input: "larger  spaces", want: []string{"larger", "spaces"}},
	}
	for _, tc := range tests {
		got := cleanInput(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("test failed: %v, expected: %v, got: %v", tc.name, tc.want, got)
		}
	}
}
