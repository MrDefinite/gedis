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

func parseLength(in *string) (int, error) {
	length := 0
	for (*in)[0:1] != "\r" {
		length = length * 10 + int((*in)[0] - '0')
		*in = (*in)[1:]
	}

	if len(*in) <= 1 {
		return 0, errors.New("mal-format string received")
	}
	if !isDelimiter(*in) {
		return 0, errors.New("mal-format string received")
	}

	*in = (*in)[2:]
	return length, nil
}

func isArray(in string) bool {
	return len(in) >= 1 && in[0:1] == ARRAYS
}

func isBulk(in string) bool {
	return len(in) >= 1 && in[0:1] == BULK_STRINGS
}

func isDelimiter(in string) bool {
	return len(in) >= 2 && in[0:2] == DELIMITER
}

func ParseRequest(input string) (*RequestCmd, error) {
	requestCmd := RequestCmd{
		Name: "",
		Params: make([]string, 0),
	}

	if !isArray(input) {
		return nil, errors.New("mal-format string received")
	}
	input = input[1:]

	length, err := parseLength(&input)
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
		if !isBulk(input) {
			return nil, errors.New("mal-format string received")
		}
		input = input[1:]

		bulkLength, err := parseLength(&input)
		if err != nil {
			return nil, err
		}

		// 2. Get bulk string
		if len(input) < bulkLength {
			return nil, errors.New("mal-format string received")
		}
		bulk := input[0:bulkLength]
		input = input[bulkLength:]
		if !isDelimiter(input) {
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

