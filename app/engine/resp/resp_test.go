package resp

import (
	"bytes"
	"net"
	"reflect"
	"testing"
	"time"
)

// mockConn is a mock net.Conn implementation to simulate reading from a connection.
type mockConn struct {
	*bytes.Buffer
}

func (m mockConn) Close() error                       { return nil }
func (m mockConn) LocalAddr() net.Addr                { return nil }
func (m mockConn) RemoteAddr() net.Addr               { return nil }
func (m mockConn) SetDeadline(t time.Time) error      { return nil }
func (m mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m mockConn) SetWriteDeadline(t time.Time) error { return nil }

// TestParseSingle tests the DecodeRESP method for various RESP types.
func TestParseSingle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		err      error
	}{
		{"array of bulk string", "*1\r\n$4\r\nping\r\n", []interface{}{"ping"}, nil},
		{"SimpleString", "+OK\r\n", "OK", nil},
		{"Error", "-Error message\r\n", "Error message", nil},
		{"Integer", ":12345\r\n", int64(12345), nil},
		{"BulkString", "$4\r\nping\r\n", "ping", nil},
		{"Array", "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", []interface{}{"foo", "bar"}, nil},
		// Add more test cases as necessary
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			conn := mockConn{Buffer: bytes.NewBufferString(tc.input)}
			parser := NewRespParser(conn)

			result, err := parser.DecodeRESP()
			if !reflect.DeepEqual(result, tc.expected) || (err != nil && tc.err == nil) || (err == nil && tc.err != nil) {
				t.Errorf("Failed %s: expected %v, got %v, expected error %v, got error %v", tc.name, tc.expected, result, tc.err, err)
			}
		})
	}
}
