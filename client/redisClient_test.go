package main

import (
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	tt := []struct {
		test    string
		payload []byte
		want    []byte
	}{
		{
			"Test SET",
			[]byte("SET A B\n"),
			[]byte("OK!"),
		},
		{
			"Test GET",
			[]byte("GET A\n"),
			[]byte("B"),
		},
		{
			"Test GET",
			[]byte("GET FAKE_KEY\n"),
			[]byte("No such key!"),
		},
		{
			"Test LPUSH",
			[]byte("LPUSH mlist a\n"),
			[]byte("1"),
		},
		{
			"Test LPUSH",
			[]byte("LPUSH mlist b\n"),
			[]byte("2"),
		},
		{
			"Test HSET",
			[]byte("HSET mmp b a\n"),
			[]byte("OK!"),
		},
		{
			"TesHGET",
			[]byte("HGET mmp b a\n"),
			[]byte("a"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.test, func(t *testing.T) {
			conn, err := net.Dial("tcp", ":9000")
			if err != nil {
				t.Error("could not connect to TCP server: ", err)
			}
			defer conn.Close()
			if _, err := conn.Write(tc.payload); err != nil {
				t.Error("could not write payload to TCP server:", err)
			}
			out := make([]byte, 1024)
			if n, err := conn.Read(out); err == nil {
				if string(out[0:n]) != string(tc.want) {
					t.Errorf("Got %s Expected %s\n", string(out[0:n]), string(tc.want))
				}
			} else {
				t.Error("could not read from connection")
			}
		})
	}
}
