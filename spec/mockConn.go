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

type MockConn struct {
	closed bool;
	input, output *bytes.Buffer;
}

func newMockConn() (result *MockConn) {
	result = &MockConn{};
	result.input = bytes.NewBufferString("");
	result.output = bytes.NewBufferString("");
	return
}

func (self *MockConn) Close() os.Error {
	self.closed = true;
	return nil;
}

func (self *MockConn) Read(bytes []byte) (int, os.Error) {
	return self.input.Read(bytes);
}

func (self *MockConn) Write(bytes []byte) (int, os.Error) {
	return self.output.Write(bytes);
}

func (self *MockConn) Send(s string) {
	self.input = bytes.NewBufferString(s);
}

func (self *MockConn) ExtractBytes() (result []byte) {
	result = self.output.Bytes();
	self.output.Reset();
	return
}

func (self *MockConn) Closed() bool {
	return self.closed;
}
