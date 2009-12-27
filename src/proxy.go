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
package proxy

import (
	"io";
	"os";
	"net";
)

const (
	defaultBufferSize = 4096;
)

func ListenOn(addr string) (result os.Error) {
	if listener, err := net.Listen("tcp", addr); err != nil {
		result = err
	} else {
		result = StartOn(listener)
	}
	return;
}

func StartOn(l net.Listener) os.Error {
	for {
		if conn, err := l.Accept(); err != nil {
			/* Not all errors from accept() mean we should kill
			the whole server, but for now let's go ahead and behave
			that way. TODO */
			return err
		} else {
			proxy := For(conn);
			go proxy.Start();
		}
	}
	return nil;
}

type Proxy struct {
	io.ReadWriteCloser;
}

func For(rwc io.ReadWriteCloser) *Proxy	{ return &Proxy{rwc} }

func (self *Proxy) Start() {
	self.Write([]byte{0x05, 0xFF});
	self.Close();
}
