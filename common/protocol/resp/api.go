//package resp
//
//import (
//	"strings"
//	"strconv"
//	"errors"
//)
//
//type RequestCmd struct {
//	Name   string
//	Params []string
//}
//
//type ResponseCmd struct {
//}
//
//func ParseResponse(res []byte) (*Node, bool, error) {
//
//
//
//	return nil, nil
//}
//
//func ParseRequest(input string) (*RequestCmd, error) {
//	requestCmd := RequestCmd{
//		Name:   "",
//		Params: make([]string, 0),
//	}
//
//	if !isArray(input) {
//		return nil, errors.New("mal-format string received")
//	}
//	input = input[1:]
//
//	length, err := parseLength(&input)
//	if err != nil {
//		return nil, err
//	}
//
//	// Begin to parse
//	for i := 0; i < length; i++ {
//		if len(input) < 1 {
//			return nil, errors.New("mal-format string received")
//		}
//
//		// Read next bulk string
//		// 1. Get string length
//		if !isBulk(input) {
//			return nil, errors.New("mal-format string received")
//		}
//		input = input[1:]
//
//		bulkLength, err := parseLength(&input)
//		if err != nil {
//			return nil, err
//		}
//
//		// 2. Get bulk string
//		if len(input) < bulkLength {
//			return nil, errors.New("mal-format string received")
//		}
//		bulk := input[0:bulkLength]
//		input = input[bulkLength:]
//		if !isDelimiter(input) {
//			return nil, errors.New("mal-format string received")
//		}
//		input = input[2:]
//
//		if i == 0 {
//			requestCmd.Name = bulk
//		} else {
//			requestCmd.Params = append(requestCmd.Params, bulk)
//		}
//	}
//
//	return &requestCmd, nil
//}
//
//func EncodeCmd(cmd string) string {
//	cmd = strings.TrimSpace(cmd)
//	tokens := strings.Split(cmd, " ")
//	cmdLength := len(tokens)
//
//	encodedCmd := arrays + strconv.Itoa(cmdLength) + delimiter
//	for i := 0; i < cmdLength; i++ {
//		token := tokens[i]
//		if token == "" {
//			continue
//		}
//		encodedCmd += bulkString + strconv.Itoa(len(token)) + delimiter
//		encodedCmd += token + delimiter
//	}
//
//	return encodedCmd
//}
//
//func EncodeResponse(response string) string {
//	return response
//}
