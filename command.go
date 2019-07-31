package susc

import (
    "bufio"
	"fmt"
)

// Command is a helper to send one command and return the reply as a string
func (c *Client) Command(str string, args ...interface{}) (*Reply, error) {
	if _, err := fmt.Fprintf(c.TCPConn, str, args...); err != nil {
		return nil, err
	}
    reader := bufio.NewReader(c.TCPConn)
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	return ParseReply(string(line))
}
