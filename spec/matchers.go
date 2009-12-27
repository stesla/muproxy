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
	"io";
	"os";
)

func mockConnTest(val interface{}, test func(*mockClient) os.Error) (err os.Error) {
	if conn, ok := val.(*mockClient); !ok {
		err = os.NewError("Not a mockClient")
	} else {
		err = test(conn)
	}
	return;
}

func Receive(b []byte) receiveMatcher	{ return receiveMatcher(b) }

type receiveMatcher []byte

func (b receiveMatcher) Should(val interface{}) (err os.Error) {
	return mockConnTest(val, func(conn *mockClient) os.Error {
		expected := ([]byte)(b);
		actual := make([]byte, len(expected));
		io.ReadFull(conn, actual);
		if !bytes.Equal(expected, actual) {
			return os.NewError(fmt.Sprintf("expected `%v` to be `%v`", actual, expected))
		}
		return nil;
	})
}
func (receiveMatcher) ShouldNot(val interface{}) os.Error {
	return os.NewError("matcher not implemented")
}

const BeClosed closedMatcher = true

type closedMatcher bool

func (closedMatcher) Should(val interface{}) os.Error {
	return mockConnTest(val, func(conn *mockClient) os.Error {
		if !conn.Closed() {
			return os.NewError("expected connection to be closed")
		}
		return nil;
	})
}
func (closedMatcher) ShouldNot(val interface{}) os.Error {
	return mockConnTest(val, func(conn *mockClient) os.Error {
		if conn.Closed() {
			return os.NewError("expected connection not to be closed")
		}
		return nil;
	})
}
