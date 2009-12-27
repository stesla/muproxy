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
)

func readerFromBytes(b []int) *Reader {
	ch := make(chan int, len(b));
	for _,c := range b {
		ch <- c;
	}
	close(ch);
	return NewReader(ch);
}

func init() {
	Describe("chanio.Reader", func() {
		It("should return EAGAIN for a negative value", func(e Example) {
			r := readerFromBytes([]int{-1});
			n, err := r.Read(make([]byte, 1));
			e.Value(n).Should(Be(0));
			e.Value(err).Should(Be(os.EAGAIN));
		});

		It("should only return EAGAIN if nothing has been read", func(e Example) {
			r := readerFromBytes([]int{'q',-1});
			buf := make([]byte, 2);
			n, err := r.Read(buf);
			e.Value(n).Should(Be(1));
			e.Value(err).Should(Be(nil));
			e.Value(len(buf)).Should(Be(2));
			e.Value(string(buf[0:n])).Should(Be("q"));
		});

		It("should return EOF for a closed channel", func(e Example) {
			r := readerFromBytes([]int{});
			n, err := r.Read(make([]byte, 1));
			e.Value(n).Should(Be(0));
			e.Value(err).Should(Be(os.EOF));
		});

		It("should read bytes", func(e Example) {
			r := readerFromBytes([]int{'f','o','o'});
			buf := make([]byte, 3);
			n, err := r.Read(buf);
			e.Value(n).Should(Be(3));
			e.Value(err).Should(Be(nil));
			e.Value(string(buf)).Should(Be("foo"));
		});

		It("should read zeros", func(e Example) {
			r := readerFromBytes([]int{0, 0});
			buf := make([]byte, 2);
			n, err := r.Read(buf);
			e.Value(n).Should(Be(2));
			e.Value(err).Should(Be(nil));
			e.Value(string(buf)).Should(Be("\x00\x00"));
		})
	});

	Describe("chanio.Writer", func() {
		It("should always write a -1", func(e Example) {
			ch := make(chan int, 1);
			w := NewWriter(ch);
			n, _ := w.Write([]byte{});
			e.Value(n).Should(Be(0));
			e.Value(<-ch).Should(Be(-1));
		});

		It("should write bytes", func(e Example) {
			ch := make(chan int, 4);
			w := NewWriter(ch);
			nw, _ := w.Write([]byte{'b','a','r'});
			r := NewReader(ch);
			buf := make([]byte, 3);
			r.Read(buf);
			e.Value(nw).Should(Be(3));
			e.Value(string(buf)).Should(Be("bar"));
		});

		It("should close", func(e Example) {
			ch := make(chan int, 1);
			w := NewWriter(ch);
			r := NewReader(ch);
			w.Close();
			n, err := r.Read(make([]byte,1));
			e.Value(n).Should(Be(0));
			e.Value(err).Should(Be(os.EOF));
		});
	});
}
