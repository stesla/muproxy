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
	"fmt";
	. "specify";
)

func init() {
	Describe("method negotiation", func() {
		Before(beforeProxySpec);

		/* Refuse all methods except:
		- 0x00 passwordless
		- 0x80 muproxy auth + world + character
		*/
		for c := 0x01; c < 0xFF; c++ {
			if c == 0x80 {
				continue
			}
			method := byte(c);
			It(fmt.Sprintf("should refuse method 0x%X", c), func(e Example) {
				conn := getClient(e);
				conn.Send(methodRequest(method));
				e.Value(conn).Should(Receive(methodResponse(0xFF)));
				e.Value(conn).Should(BeClosed);
			});
		}

		It("should accept passwordless", func(e Example) {
			conn := getClient(e);
			conn.Send(methodRequest(0x00));
			e.Value(conn).Should(Receive(methodResponse(0x00)));
			e.Value(conn).ShouldNot(BeClosed);
		});

		It("should accept muproxy auth + world + character", func(e Example) {
			conn := getClient(e);
			conn.Send(methodRequest(0x80));
			e.Value(conn).Should(Receive(methodResponse(0x80)));
			e.Value(conn).ShouldNot(BeClosed);
		});

		After(afterProxySpec);
	})
}
