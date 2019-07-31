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
