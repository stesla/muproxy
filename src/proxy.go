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
	MethodNoPassword	= 0x00;
	MethodMUProxy		= 0x80;
	MethodNone		= 0xFF;
	Version			= 0x05;

	AddrIP4 = 0x01;

	ReqConnect = 0x01;
	ReqBind = 0x02;
	ReqAssociate = 0x03;

	CommandNotSupported = 0x07;
)

const (
	defaultBufferSize	= 4096;
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
	var err os.Error;
	err = self.negotiateMethod();
	if err != nil {
		self.Close();
		return;
	}
	self.Write([]byte{Version, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00});
	self.Close();
}

func (self *Proxy) negotiateMethod() os.Error {
	if method, err := self.selectMethod(); err != nil {
		return err
	} else {
		if _, err := self.Write([]byte{Version, method}); err != nil {
			return err
		}
		if method == MethodNone {
			self.Close()
		}
	}
	return nil;
}

func (self *Proxy) selectMethod() (method byte, error os.Error) {
	buf := make([]byte, defaultBufferSize)[0:2];
	if _, error = io.ReadFull(self, buf); error != nil {
		method = MethodNone
	} else {
		if buf[0] != Version {
			method, error = MethodNone, os.NewError("Invalid Version in request")
		} else {
			buf = buf[0:buf[1]];
			if _, error = io.ReadFull(self, buf); error != nil {
				method = MethodNone
			} else {
				method = findAcceptableMethod(buf)
			}
		}
	}
	return;
}

func findAcceptableMethod(methods []byte) (method byte) {
	for _, method = range methods {
		if method == MethodNoPassword || method == MethodMUProxy {
			return
		}
	}
	return MethodNone;
}
