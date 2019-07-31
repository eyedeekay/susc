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
