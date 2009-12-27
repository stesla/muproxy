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
	"os";
)

const (
	bufferSize	= 512;
	chanSize	= 8;
)

func newBuffer() []byte	{ return make([]byte, bufferSize)[0:0] }

func fillBuf(ch <-chan []byte, buf *bytes.Buffer) (isClosed bool) {
	bytes := <-ch;
	if bytes != nil {
		buf.Write(bytes)
	} else {
		isClosed = closed(ch)
	}
	return;
}

func mockConnection() (rwc *mockServer, client *mockClient) {
	in := make(chan []byte, chanSize);
	out := make(chan []byte, chanSize);
	rwc = &mockServer{in: in, out: out, buf: bytes.NewBuffer(newBuffer())};
	client = &mockClient{in: in, out: out, buf: bytes.NewBuffer(newBuffer())};
	return;
}

type mockServer struct {
	in	<-chan []byte;
	out	chan<- []byte;

	closed	bool;
	buf	*bytes.Buffer;
}

func (self *mockServer) Read(b []byte) (n int, err os.Error) {
	if self.buf.Len() >= len(b) {
		return self.buf.Read(b)
	}

	for !self.closed && self.buf.Len() < len(b) {
		self.closed = fillBuf(self.in, self.buf)
	}
	return self.buf.Read(b);
}

func (self *mockServer) Write(b []byte) (n int, err os.Error) {
	self.out <- b;
	return len(b), nil;
}

func (self *mockServer) Close() os.Error {
	close(self.out);
	return nil;
}

type mockClient struct {
	in	chan<- []byte;
	out	<-chan []byte;

	closed	bool;
	buf	*bytes.Buffer;
}

func (self *mockClient) Close()	{ close(self.in) }

func (self *mockClient) Closed() bool {
	if !self.closed {
		bytes, ok := <-self.out;
		if ok {
			self.buf.Write(bytes);
			self.closed = closed(self.out);
		}
	}
	return self.closed;
}

func (self *mockClient) Read(b []byte) (n int, err os.Error) {
	if self.buf.Len() >= len(b) {
		return self.buf.Read(b)
	}

	for !self.closed && self.buf.Len() < len(b) {
		self.closed = fillBuf(self.out, self.buf)
	}
	return self.buf.Read(b);
}

func (self *mockClient) Send(b []byte)	{ self.in <- b }
