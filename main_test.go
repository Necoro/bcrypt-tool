package main

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alecthomas/kong"
)

func TestHelpStrings(t *testing.T) {
	argParser, err := kong.New(&CLI, kongParserOptions()...)
	assert.NoError(t, err, "building CLI interface")
	for _, cmd := range argParser.Model.Children {
		t.Run(cmd.Name, func(t *testing.T) {
			assert.NotZero(t, cmd.Help, "help text should not be empty")
		})
	}
}
