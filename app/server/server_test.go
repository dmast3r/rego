package main

import (
	"net"
	"testing"
	"time"
)

func TestHandleConnection(t *testing.T) {
	serverConn, clientConn := net.Pipe()

	go handleConnection(serverConn)

	sendPing := func() string {
		// Simulate client sending data
		cmd := "*1\r\n$4\r\nPING\r\n"
		_, err := clientConn.Write([]byte(cmd))
		if err != nil {
			t.Fatalf("Failed to write to client connection: %v", err)
		}

		// Give some time for the server to process the command
		time.Sleep(time.Second)

		// Read response from server
		buf := make([]byte, 1024)
		n, err := clientConn.Read(buf)
		if err != nil {
			t.Fatalf("Failed to read from client connection: %v", err)
		}

		return string(buf[:n])
	}

	// Send the first PING and check the response
	response := sendPing()
	if response != "+PONG\r\n" {
		t.Errorf("Expected '+PONG\\r\\n' for the first PING, got '%s'", response)
	}

	// Send the second PING and check the response
	response = sendPing()
	if response != "+PONG\r\n" {
		t.Errorf("Expected '+PONG\\r\\n' for the second PING, got '%s'", response)
	}

	clientConn.Close()
}
