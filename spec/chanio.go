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
package main //chanio

import (
//	"fmt";
	"os";
)

type Reader struct {
	ch <-chan int;
}

func NewReader(ch <-chan int) *Reader { return &Reader{ch} }

func (self *Reader) Read(b []byte) (n int, error os.Error) {
	for len(b) > 0 {
		switch c := <-self.ch; {
		case c < 0:
			if n > 0 {
				b = b[0:0];
				return
			} else {
				return n, os.EAGAIN
			}
		case c == 0 && closed(self.ch):
			return n, os.EOF
		default:
			b[0] = byte(c);
			b = b[1:];
		}
		n++;
	}
	return
}

type Writer struct {
	ch chan<- int;
}

func NewWriter(ch chan<- int) *Writer { return &Writer{ch} }

func (self *Writer) Write(b []byte) (n int, error os.Error) {
	for _, c := range b {
		self.ch <- int(c);
	}
	self.ch <- -1;
	return len(b), nil
}

func (self *Writer) Close() os.Error {
	close(self.ch);
	return nil;
}
