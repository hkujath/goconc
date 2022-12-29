package main

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {

	tests := []struct {
		name  string
		input []string
		want  []cmdArgs
	}{
		{
			name:  "Two commands",
			input: []string{"abc", "def", "::", "abcd", "defg"},
			want: []cmdArgs{
				{name: "abc", args: []string{"def"}},
				{name: "abcd", args: []string{"defg"}},
			}},
		{
			name:  "One command",
			input: []string{"abc", "def", "hij"},
			want: []cmdArgs{
				{name: "abc", args: []string{"def", "hij"}},
			}},
		{
			name:  "No commands",
			input: []string{},
			want:  nil},
		{
			name:  "Two commands with short declaration",
			input: []string{"abc", "def", ":", "defg"},
			want: []cmdArgs{
				{name: "abc", args: []string{"def"}},
				{name: "abc", args: []string{"defg"}},
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseArgs(tt.input)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArgs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
