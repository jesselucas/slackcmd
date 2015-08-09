/*
	Boolean only flags for Slack Slash Commands
*/
package slack

import (
	"fmt"
)

type FlagSet struct {
	Flags []Flag
	Usage string // help message
}

func (fs *FlagSet) addFlag(f Flag) {
	fs.Flags = append(fs.Flags, f)
}

func (fs *FlagSet) String() string {

	s := fmt.Sprintf(
		"`%v` \n",
		fs.Usage,
	)

	s += "```Flags: \n"
	for _, flag := range fs.Flags {
		s += fmt.Sprint(flag)
	}
	s += "```"

	return s
}

type Flag struct {
	Name      string // full name
	ShortName string // single letter name
	Usage     string // help message
	Callback  func()
}

func (f Flag) String() string {
	return fmt.Sprintf(
		"• --%v (-%v): %v \n",
		f.Name,
		f.ShortName,
		f.Usage,
	)
}

func SetFlag(fs *FlagSet, name string, shortname string, usage string, callback func()) {
	f := Flag{name, shortname, usage, callback}

	fs.addFlag(f)

}

func ParseFlags(fs *FlagSet, flags []string) (h bool, s string) {
	for _, flag := range flags {
		// Test for help command
		if flag == "help" || flag == "h" {
			return true, fmt.Sprint(fs)
		}

		// check each flag passed with all registerd
		for _, fsFlag := range fs.Flags {
			if flag == fsFlag.Name || flag == fsFlag.ShortName {
				fsFlag.Callback()
			}
		}
	}

	// Return string without flags
	return false, ""
}
