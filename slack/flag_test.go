package slack

import (
	"fmt"
	"reflect"
	"testing"
)

func TestVar(t *testing.T) {
	// Create FlagSet to store flags
	fs := &FlagSet{}

	tests := []struct {
		name      string
		shortname string
		usage     string
		command   func()
	}{
		{"help", "h", "List all flags", func() { fmt.Println("help func") }},
		{"channel", "c", "Post command result to current channel", func() { fmt.Println("channel func") }},
	}

	for i, test := range tests {
		// Set the flag. Must do this before testing against actual
		SetFlag(fs, test.name, test.shortname, test.usage, test.command)

		actual := fs.Flags[i]

		if actual.Name != test.name {
			t.Errorf("Test %d errored. Name should be %v but is %v '", i, fs.Flags[i].Name, test.name)
		}
		if actual.ShortName != test.shortname {
			t.Errorf("Test %d errored. ShortName should be %v but is %v '", i, fs.Flags[i].ShortName, test.shortname)
		}
		if actual.Usage != test.usage {
			t.Errorf("Test %d errored. Usage should be %v but is %v '", i, fs.Flags[i].Usage, test.usage)
		}
	}
}

func TestParseFlags(t *testing.T) {
	fs := &FlagSet{}
	SetFlag(fs, "channel", "c", "channel flag usage", func() {})

	tests := []struct {
		commands       string
		actualCommands string
	}{
		{"golang links -channel", "golang links"},
		{"golang links -private", "golang links"},
		{"wiki -help", fmt.Sprint(fs)},
	}

	for _, test := range tests {
		// Separate flags
		commands, flags := SeparateFlags(test.commands)

		// test ParseFlags
		h, s := ParseFlags(fs, flags)

		// help message was returned
		if h {
			if reflect.DeepEqual(s, test.actualCommands) == false {
				t.Errorf("Test errored. Usage should be %v but is %v '", test.actualCommands, s)
			}
		} else {
			if reflect.DeepEqual(commands, test.actualCommands) == false {
				t.Errorf("Test errored. Usage should be %v but is %v '", test.actualCommands, commands)
			}
		}
	}

}
