package susc

import (
	"fmt"
)

func (c *Client) Hello() error {
	r, err := c.Command("HELLO VERSION MIN=3.0 MAX=3.2\n")
	if err != nil {
		return err
	}

	if r.Topic != "HELLO" {
		return fmt.Errorf("Unknown Reply: %+v\n", r)
	}

	if r.Pairs["RESULT"] != "OK" {
		return fmt.Errorf("Handshake did not succeed\nReply:%+v\n", r)
	}

	return nil
}
