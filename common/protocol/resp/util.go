package resp

import (
	"strings"
	"github.com/sirupsen/logrus"
	"strconv"
	"github.com/pkg/errors"
)

var log = logrus.New()

const (
	SIMPLE_STRING = "+"
	ERRORS        = "-"
	INTEGERS      = ":"
	BULK_STRINGS  = "$"
	ARRAYS        = "*"
	DELIMITER     = "\r\n"
)

type RequestCmd struct {
	cmdName *string
	params  []*string
}

type ResponseCmd struct {
}

func ParseResponse(input *string) string {
	//tokens := strings.Split(*input, " ")
	//length := len(tokens)

	return ""
}

func ParseRequest(input *string) (*RequestCmd, error) {
	*input = strings.TrimSpace(*input)

	requestCmd := RequestCmd{}

	lengthStr := ""
	// Get cmd length
	for *input != "" {
		if string([]rune(*input)[0]) == ARRAYS {
			*input = (*input)[1:]
			continue
		}
		if (*input)[0:2] == DELIMITER {
			*input = (*input)[2:]
			break
		}
		lengthStr += string([]rune(*input)[0])
		*input = (*input)[1:]
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, err
	}
	log.Infof("Length of the cmd and params is: %d", length)

	// Begin to parse
	for i := 0; i < length; i++ {
		// Read next bulk string
		// 1. Get string length
		mark := string([]rune(*input)[0])
		if mark != BULK_STRINGS {
			return nil, errors.Errorf("mal-format string '%s' received", *input)
		}
		*input = (*input)[1:]

		bulkLengthStr := ""
		for *input != "" {
			if (*input)[0:2] == DELIMITER {
				*input = (*input)[2:]
				break
			}
			bulkLengthStr += string([]rune(*input)[0])
			*input = (*input)[1:]
		}

		// 2. Get bulk string
		bulkLength, err := strconv.Atoi(bulkLengthStr)
		if err != nil {
			return nil, err
		}
		bulk := (*input)[0:bulkLength]
		*input = (*input)[bulkLength:]
		if (*input)[0:2] != DELIMITER {
			return nil, errors.Errorf("mal-format string '%s' received", *input)
		}
		*input = (*input)[2:]

		if i == 0 {
			requestCmd.cmdName = &bulk
		} else {
			requestCmd.params = append(requestCmd.params, &bulk)
		}
	}

	return &requestCmd, nil
}

func EncodeCmd(cmd *string) string {
	*cmd = strings.TrimSpace(*cmd)
	tokens := strings.Split(*cmd, " ")
	cmdLength := len(tokens)

	encodedCmd := ARRAYS + string(cmdLength) + DELIMITER
	for i := 1; i < cmdLength; i++ {
		token := tokens[i]
		if token == "" {
			continue
		}
		encodedCmd += BULK_STRINGS + string(len(token)) + DELIMITER
		encodedCmd += token + DELIMITER
	}

	return encodedCmd
}
