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
	"os";
	. "specify";
	"../src/proxy";
)

type join struct {
	ch <-chan bool;
}

func getConnection(c Context) *MockConn {
	conn, ok := c.GetField("connection").(*MockConn);
	if !ok {
		c.Error(os.NewError("connection not a *MockConn"))
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
	c := newMockConn();
	e.SetField("connection", c);
	p := proxy.For(c);
	e.SetField("proxy", p);
	ch := make(chan bool, 1);
	e.SetField("join", &join{ch});
	go func() {
		p.Start();
		ch <- true;
	}();
}

func afterProxySpec(c Context) {
	conn := getConnection(c);
	ch := getJoin(c);
	conn.Close();
	if val, ok := <-ch; !(val && ok) {
		c.Error(os.NewError("Proxy did not exit on close"))
	}
}