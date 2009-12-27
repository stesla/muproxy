/*
Copyright (c) 2009 Samuel Tesla <samuel.tesla@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"bytes";
	"encoding/binary";
	"io";
	"os";
	"strconv";
	"strings";
	. "specify";
	"../src/proxy";
)

type join struct {
	ch <-chan bool;
}

func getClient(c Context) *mockConn {
	conn, ok := c.GetField("client").(*mockConn);
	if !ok {
		c.Error(os.NewError("client is not a *mockConn"))
	}
	return conn;
}

func getJoin(c Context) <-chan bool {
	join, ok := c.GetField("join").(*join);
	if !ok {
		c.Error(os.NewError("join not a *join"))
	}
	return join.ch;
}

func beforeProxySpec(e Example) {
	client, server := mockConnection();
	e.SetField("client", client);
	p := proxy.For(server);
	e.SetField("proxy", p);
	ch := make(chan bool, 1);
	e.SetField("join", &join{ch});
	go func() {
		p.Start();
		ch <- true;
	}();
}

func afterProxySpec(c Context) {
	conn := getClient(c);
	ch := getJoin(c);
	conn.Close();
	if val, ok := <-ch; !(val && ok) {
		c.Error(os.NewError("Proxy did not exit on close"))
	}
}

func methodRequest(m byte) []byte {
	return []byte{proxy.Version, 1, m};
}

func methodResponse(m byte) []byte {
	return []byte{proxy.Version, m};
}

func negotiateMethod(c Context, rw io.ReadWriter) {
	rw.Write(methodRequest(proxy.MethodNoPassword));
	if _, err := rw.Read(make([]byte,3)); err != nil {
		c.Error(err)
	}
	return
}

func socksMsg(mtype, atype byte, addr string, port uint16) []byte {
	buf := bytes.NewBuffer([]byte{proxy.Version, mtype, 0x00, atype});
	switch atype {
	case proxy.AddrIP4:
		buf.Write(ip4Bytes(addr));
	default:
		panic("not implemented");
	}
	if err := binary.Write(buf, binary.BigEndian, port); err != nil {
		panic("error encoding port");
	}
	return buf.Bytes();
}

func ip4Bytes(addr string) (buf []byte) {
	buf = make([]byte, 4);
	octets := strings.Split(addr, ".", 0);
	if len(octets) != 4 {
		panic("expected dotted quad");
	}
	for i, s := range octets {
		if n, err := strconv.Atoi(s); err != nil {
			panic(err.String());
		} else {
			buf[i] = byte(n)
		}
	}
	return
}
