package susc

import (
	"fmt"
)

// StreamAccept asks SAM to accept a TCP-Like connection
func (c *Client) StreamAccept(id int32) (*Reply, error) {
	r, err := c.Command("STREAM ACCEPT ID=%d SILENT=false\n", id)
	if err != nil {
		return nil, err
	}

	// TODO: move check into Command()
	if r.Topic != "STREAM" || r.Type != "STATUS" {
		return nil, fmt.Errorf("Unknown Reply: %+v\n", r)
	}

	result := r.Pairs["RESULT"]
	if result != "OK" {
		return nil, fmt.Errorf("Reply Error")
	}

	return r, nil
}
