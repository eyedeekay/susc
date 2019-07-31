package susc

// Write implements the TCPConn Write method.
func (c *Client) Write(b []byte) (int, error) {
    return c.TCPConn.Write(b)
}
