package resp

import (
	"errors"
	"fmt"
	"github.com/dmast3r/rego/app/engine/data_structures"
	"strconv"
	"strings"
)

// RedisCmd /**
// Represent a Redis command with its string and arguments
type RedisCmd struct {
	Cmd  string
	Args []string
}

func (cmd *RedisCmd) GetRepr() string {
	return fmt.Sprintf("command: %s\nargs: %v", cmd.Cmd, cmd.Args)
}

func GetRedisCmd(cmdString []interface{}) (RedisCmd, error) {
	tokens := toArrayString(cmdString)
	return RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func (cmd *RedisCmd) Execute() (string, error) {
	switch cmd.Cmd {
	case "PING":
		return "+PONG\r\n", nil
	case "ECHO":
		return fmt.Sprintf("$%d\r\n%s\r\n", len(cmd.Args[0]), cmd.Args[0]), nil
	case "SET":
		var expiry int64 = -1
		var expiryParseError error

		if len(cmd.Args) == 4 && strings.ToUpper(cmd.Args[2]) == "PX" {
			expiry, expiryParseError = strconv.ParseInt(cmd.Args[3], 10, 64)
			if expiryParseError != nil {
				return "", expiryParseError
			}
		}

		err := data_structures.GetHashMap().Set(cmd.Args[0], cmd.Args[1], expiry)
		if err != nil {
			return "", err
		}

		return "+OK\r\n", nil
	case "GET":
		val, exists, err := data_structures.GetHashMap().Get(cmd.Args[0])
		if err != nil {
			return "", err
		}

		if !exists {
			return "$-1\r\n", nil
		}

		strVal := val.(string)
		return fmt.Sprintf("$%d\r\n%s\r\n", len(strVal), strVal), nil
	}
	return "", errors.New("invalid command")
}

func toArrayString(cmdString []interface{}) []string {
	result := make([]string, len(cmdString))
	for i := range result {
		result[i] = cmdString[i].(string)
	}
	return result
}
