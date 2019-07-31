package susc

import (
    "bufio"
)

func (c *Client) ReadLine() (string, error) {
    reader := bufio.NewReader(c.TCPConn)
	bytes, _, err := reader.ReadLine()
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}
