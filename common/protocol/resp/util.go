package resp

import (
	"strings"
	"strconv"
	"errors"
)


const (
	SIMPLE_STRING = "+"
	ERRORS        = "-"
	INTEGERS      = ":"
	BULK_STRINGS  = "$"
	ARRAYS        = "*"
	DELIMITER     = "\r\n"
)

type RequestCmd struct {
	Name   string
	Params []string
}

type ResponseCmd struct {
}

func ParseResponse(input string) string {
	//tokens := strings.Split(*input, " ")
	//length := len(tokens)

	return ""
}

func parseLength(input string) (int, string, error) {
	length := 0
	for input[0] != '\r' {
		length = length * 10 + int(input[0] - '0')
		input = input[1:]
	}

	if len(input) <= 1 {
		return 0, "", errors.New("mal-format string received")
	}
	if input[0:2] != DELIMITER {
		return 0, "", errors.New("mal-format string received")
	}

	input = input[2:]
	return length, input, nil
}

func ParseRequest(input string) (*RequestCmd, error) {
	input = strings.TrimSpace(input)

	requestCmd := RequestCmd{}

	if string([]rune(input)[0]) != ARRAYS {
		return nil, errors.New("mal-format string received")
	}
	input = input[1:]

	length, input, err := parseLength(input)
	if err != nil {
		return nil, err
	}

	// Begin to parse
	for i := 0; i < length; i++ {
		if len(input) < 1 {
			return nil, errors.New("mal-format string received")
		}

		// Read next bulk string
		// 1. Get string length
		mark := string([]rune(input)[0])
		if mark != BULK_STRINGS {
			return nil, errors.New("mal-format string received")
		}
		input = input[1:]

		bulkLength, input, err := parseLength(input)
		if err != nil {
			return nil, err
		}

		// 2. Get bulk string
		if len(input) < bulkLength {
			return nil, errors.New("mal-format string received")
		}
		bulk := input[0:bulkLength]
		input = input[bulkLength:]
		if len(input) < 2 || input[0:2] != DELIMITER {
			return nil, errors.New("mal-format string received")
		}
		input = input[2:]

		if i == 0 {
			requestCmd.Name = bulk
		} else {
			requestCmd.Params = append(requestCmd.Params, bulk)
		}
	}

	return &requestCmd, nil
}

func EncodeCmd(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	tokens := strings.Split(cmd, " ")
	cmdLength := len(tokens)

	encodedCmd := ARRAYS + strconv.Itoa(cmdLength) + DELIMITER
	for i := 0; i < cmdLength; i++ {
		token := tokens[i]
		if token == "" {
			continue
		}
		encodedCmd += BULK_STRINGS + strconv.Itoa(len(token)) + DELIMITER
		encodedCmd += token + DELIMITER
	}

	return encodedCmd
}

func EncodeResponse(response string) string {
	return response
}

