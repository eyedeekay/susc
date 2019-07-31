package susc

import (
	"fmt"
)

// StreamTCPConnect asks SAM for a TCP-Like connection to dest, has to be called on a new Client
func (c *Client) StreamTCPConnect(id int32, dest string) error {
	r, err := c.Command("STREAM CONNECT ID=%d DESTINATION=%s\n", id, dest)
	if err != nil {
		return err
	}

	// TODO: move check into Command()
	if r.Topic != "STREAM" || r.Type != "STATUS" {
		return fmt.Errorf("Unknown Reply: %+v\n", r)
	}

	result := r.Pairs["RESULT"]
	if result != "OK" {
		return fmt.Errorf("Reply Error")
	}

	return nil
}
