package main

import (
	"fmt"
)

type Beats1 struct {
}

func (b Beats1) Request(sc *SlashCommand) (*CommandPayload, error) {
	// create payload
	cp := &CommandPayload{
		Channel:       fmt.Sprintf("@%v", sc.UserName),
		Username:      "Beats1",
		Emoji:         ":metal:",
		SlashResponse: false,
		SendPayload:   true,
	}

	cp.Text = "Beats1"

	// get the latest aacs
	// http://itsliveradiobackup.apple.com/streams/hub02/session02/64k/prog.m3u8

	// sort on aac and check the latest id tag data

	return cp, nil
}
