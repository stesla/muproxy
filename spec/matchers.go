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
	"fmt";
	"os";
	"strings";
)

func Receive(s string) *receiveMatcher { return (*receiveMatcher)(&s); }
type receiveMatcher string;
func (s receiveMatcher) Bytes() []byte { return strings.Bytes((string)(s)) }
func (s receiveMatcher) Should(val interface{}) (err os.Error) {
	if conn, ok := val.(*MockConn); !ok {
		err = os.NewError("Not a MockConn")
	} else {
		expected := s.Bytes();
		actual := conn.ExtractBytes();
		if !bytes.Equal(expected, actual) {
			err = os.NewError(fmt.Sprintf("expected `%v` to be `%v`", actual, expected))
		}
	}
	return
}
func (s receiveMatcher) ShouldNot(val interface{}) os.Error { return os.NewError("matcher not implemented") }

const BeClosed closedMatcher = true;
type closedMatcher bool;
func (closedMatcher) Should(val interface{}) (err os.Error) {
	if conn, ok := val.(*MockConn); !ok {
		err = os.NewError("Not a MockConn")
	} else {
		if !conn.Closed() {
			err = os.NewError("expected connection to be closed");
		}
	}
	return
}
func (closedMatcher) ShouldNot(val interface{}) os.Error { return os.NewError("matcher not implemented") }
