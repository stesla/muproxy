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
	"../src/proxy";
)

func init() {
	Describe("method negotiation", func() {
		Before(beforeProxySpec);

		/* Refuse all methods except:
		- 0x00 passwordless
		- 0x80 muproxy auth + world + character
		*/
		for c := 0x01; c < 0xff; c++ {
			if c == 0x80 {
				continue
			}
			method := byte(c);
			It(fmt.Sprintf("should refuse method 0x%X", c), func(e Example) {
				conn := getClient(e);
				conn.Write(methodRequest(method));
				e.Value(conn).Should(Receive(methodResponse(proxy.MethodNone)));
				e.Value(conn).Should(BeEndOfFile);
			});
		}

		It("should accept passwordless", func(e Example) {
			conn := getClient(e);
			conn.Write(methodRequest(proxy.MethodNoPassword));
			e.Value(conn).Should(Receive(methodResponse(proxy.MethodNoPassword)));
		});

		It("should deny bind requests", func(e Example) {
			conn := getClient(e);
			negotiateMethod(e, conn);
			conn.Write(socksMsg(proxy.ReqBind, proxy.AddrIP4, "10.1.2.3", 4000));
			e.Value(conn).Should(Receive(socksMsg(proxy.CommandNotSupported, proxy.AddrIP4, "0.0.0.0", 0)));
		});

		It("should deny associate requests", func(e Example) {
			conn := getClient(e);
			negotiateMethod(e, conn);
			conn.Write(socksMsg(proxy.ReqAssociate, proxy.AddrIP4, "10.1.2.3", 4000));
			e.Value(conn).Should(Receive(socksMsg(proxy.CommandNotSupported, proxy.AddrIP4, "0.0.0.0", 0)));
		});

		It("should accept connect requests", nil);

		It("should accept muproxy auth + world + character", func(e Example) {
			conn := getClient(e);
			conn.Write(methodRequest(0x80));
			e.Value(conn).Should(Receive(methodResponse(proxy.MethodMUProxy)));
		});

		After(afterProxySpec);
	})
}
