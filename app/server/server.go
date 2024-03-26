package main

import (
	"fmt"
	"github.com/dmast3r/rego/app/engine/resp"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	/**
	* Create a listener on port 6379.
	 */
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	/*
	* Make a buffered channel of size 1. This channel will listen to termination events.
	 */
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	/**
	* Continuously listen to new connections on this port. If a new connection has been established, then spin a new
	* goroutine to handle it.
	 */
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				continue
			}
			go handleConnection(conn)
		}
	}()

	<-sigChan
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("Error closing the connection!")
		}
	}()

	lastInteractionTime := time.Now()
	for time.Since(lastInteractionTime) <= 5*time.Minute {
		respDecodedValue, err := resp.NewRespParser(conn).DecodeRESP()

		if err != nil {
			if err != io.EOF {
				fmt.Println("Error parsing the result input", err)
			}
			continue
		}

		lastInteractionTime = time.Now()

		cmdString, ok := respDecodedValue.([]interface{})
		if !ok {
			fmt.Println("expected the command to be an array but wasn't")
		}

		cmd, err := resp.GetRedisCmd(cmdString)
		if err != nil {
			fmt.Println("Failed to obtain a valid command from the command string", err)
		}

		result, err := cmd.Execute()
		if err != nil {
			fmt.Printf("Failed to execute command %s with error: %v", cmd.GetRepr(), err)
		}

		_, err = conn.Write([]byte(result))
		if err != nil {
			fmt.Println("Couldn't write to the connection object. Error occurred", err)
		}
	}
}
