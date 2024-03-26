# rego

`rego` is a lightweight, Redis-like data store implementation written in Go. It aims to provide a simple and intuitive API for interacting with key-value pairs, supporting basic commands like SET, GET, and handling key expiries.

## Features

- SET/GET commands: Store and retrieve data associated with a key.
- Expiry: Set expiration on keys.
- PING command: Check the connection to the server.
- ECHO command: Echo back the given string.