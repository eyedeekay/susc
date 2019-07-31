package susc

import (
	"fmt"
)

func (c *Client) CreateStreamSession(id int32, dest, sigtype, options string) (string, error) {
	if dest == "" {
		dest = "TRANSIENT"
	}
	r, err := c.Command(
		"SESSION CREATE STYLE=STREAM ID=%d DESTINATION=%s %s %s\n",
		id,
		dest,
		sigtype,
		options,
	)
	if err != nil {
		return "", err
	}

	// TODO: move check into Command()
	if r.Topic != "SESSION" || r.Type != "STATUS" {
		return "", fmt.Errorf("Unknown Reply: %+v\n", r)
	}

	result := r.Pairs["RESULT"]
	if result != "OK" {
		return "", fmt.Errorf("Reply error")
	}
	return r.Pairs["DESTINATION"], nil
}
