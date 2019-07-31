# susc
Simplest Useful SAM Streaming Client.

This isn't intended to be very "good" right now, but rather to illustrate the
simplest ways the concepts of SAM map onto it's clearnet equivalents. It's been
created as a set of examples for Def Con 27. When it's done being a basic
example it might become a socket library, but probably not. sam3 is better.
accept.go
============

```
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
```
client.go
============

```
package susc

import (
	"net"
    	"encoding/binary"
    	"encoding/base32"
	"encoding/base64"
)

type Client struct {
	*net.TCPConn
}

var (
	i2pB64enc *base64.Encoding = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-~")
	i2pB32enc *base32.Encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")
)

func NewClient() (*Client, error) {
	//var err error
    var c Client
    samaddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:7656")
    if err != nil {
		return nil, err
	}
	c.TCPConn, err = net.DialTCP("tcp", nil, samaddr)
	if err != nil {
		return nil, err
	}
	return &c, nil
}


// Base64 returns the base64 of the local tunnel
func Base64(destination string) string {
    if destination != "" {
		s, _ := i2pB64enc.DecodeString(destination)
		alen := binary.BigEndian.Uint16(s[385:387])
		return i2pB64enc.EncodeToString(s[:387+alen])
	}
	return ""
}
```
command.go
============

```
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
```
conn_test.go
============

```
package susc

import (
	"testing"
	//"time"
)

func TestMe(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		panic(err)
	}
	err = c.Hello()
	if err != nil {
		panic(err)
	}
	dest, err := c.CreateStreamSession(int32(1337), "", "", "")
	if err != nil {
		panic(err)
	}
	t.Log(dest)
	c.StreamAccept(int32(1337))

	c2, err := NewClient()
	if err != nil {
		panic(err)
	}
	err = c2.Hello()
	if err != nil {
		panic(err)
	}
	dest2, err := c2.CreateStreamSession(int32(1776), "", "", "")
	if err != nil {
		panic(err)
	}
	t.Log("Connecting", dest2, "to", Base64(dest))
	c2.StreamTCPConnect(int32(1776), Base64(dest))

	bytes, err := c2.Write([]byte("HELLO SAM"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Wrote", bytes, "bytes")
	reply, err := c.ReadLine()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Read", reply)

}
```
connect.go
============

```
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
```
hello.go
============

```
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
```
read.go
============

```
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
```
reply.go
============

```
package susc

import (
	"fmt"
	"strings"
)

type Reply struct {
	Topic string
	Type  string

	Pairs map[string]string
}

func ParseReply(line string) (*Reply, error) {
	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return nil, fmt.Errorf("Malformed Reply.\n%s\n", line)
	}

	r := &Reply{
		Topic: parts[0],
		Type:  parts[1],
		Pairs: make(map[string]string, len(parts)-2),
	}

	for _, v := range parts[2:] {
		kvPair := strings.SplitN(v, "=", 2)
		if kvPair != nil {
			if len(kvPair) != 2 {
				return nil, fmt.Errorf("Malformed key-value-pair.\n%s\n", kvPair)
			}
		}

		r.Pairs[kvPair[0]] = kvPair[len(kvPair)-1]
	}

	return r, nil
}
```
session.go
============

```
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
```
write.go
============

```
package susc

// Write implements the TCPConn Write method.
func (c *Client) Write(b []byte) (int, error) {
    return c.TCPConn.Write(b)
}
```
