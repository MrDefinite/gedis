package resp

import (
	"errors"
	"bytes"
	"strconv"
)

type DataType int8

const (
	TypeNone         DataType = iota
	TypeSimpleString
	TypeBulkString
	TypeError
	TypeInteger
	TypeArray
)

type pausePos int8

const (
	noneParsing                pausePos = iota
	typeParsing
	simpleStringLengthChecking
	simpleStringParsing
	simpleStringEndChecking
	errorStringLengthChecking
	errorStringParsing
	errorStringEndChecking
	integerLengthChecking
	integerParsing
	integerEndChecking
)

const (
	simpleString = "+"
	gedisError   = "-"
	integers     = ":"
	bulkString   = "$"
	arrays       = "*"
	delimiter    = "\r\n"
)

const (
	// 512 MB
	MaxDataSize = 1024 * 1024 * 512
	// 16 MB
	MaxDataSizeReadPerTime = 1024 * 1024 * 16
)

var (
	malFormatErr        = errors.New("mal-format string received")
	nullHandlerErr      = errors.New("failed to init Parser: null handler")
	emptyBufferErr      = errors.New("failed to parse: empty buffer")
	unrecognizedTypeErr = errors.New("failed to parse: unrecognized byte")
	illegalTypeErr      = errors.New("failed to parse: illegal data type")
)

type Handler interface {
	handle()
}

/**
	For simple string, bulk string and error, data is string type
	For integer, data is int type
	For array, data is a list containing more nodes
 */
type Node struct {
	dataType DataType
	data     interface{}
}

//func isArray(in string) bool {
//	return len(in) >= 1 && in[0:1] == arrays
//}
//
//func isBulk(in string) bool {
//	return len(in) >= 1 && in[0:1] == bulkString
//}

//func isDelimiter(buffer *bytes.Buffer) bool {
//	return buffer.Len() >= 2 && bytes.Compare(delimiterBytes, buffer.Next(2)) == 0
//}

//func readStringToDelimiter(buffer *bytes.Buffer) (string, error) {
//	var data string
//	var c byte
//	var err error
//	for c, err = buffer.ReadByte(); c != '\r'; {
//		if err != nil {
//			return "", err
//		}
//		data += string(c)
//	}
//	c, err = buffer.ReadByte()
//	if err != nil {
//		return "", err
//	}
//	if c != '\n' {
//		return "", malFormatErr
//	}
//
//	return data, nil
//}
//
//func readString(buffer bytes.Buffer, length int) (string, error) {
//
//}

func readIntegerToDelimiter(buffer *bytes.Buffer) (int, error) {
	var data int
	var isFirstChar = true
	var isNeg = false
	var c byte
	var err error
	for c, err = buffer.ReadByte(); c != '\r'; {
		if err != nil {
			return 0, err
		}
		if isFirstChar {
			if c == '-' {
				isNeg = true
			}
			isFirstChar = false
			continue
		}
		data = data*10 + int(c-'0')
	}
	c, err = buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	if c != '\n' {
		return 0, malFormatErr
	}

	if isNeg {
		data = -data
	}

	return data, nil
}

//func parseLength(in *string) (int, error) {
//	// Parse the length
//	length := 0
//	for (*in)[0:1] != "\r" {
//		length = length*10 + int((*in)[0]-'0')
//		*in = (*in)[1:]
//	}
//
//	if len(*in) <= 1 {
//		return 0, malFormatErr
//	}
//	if !isDelimiter(*in) {
//		return 0, malFormatErr
//	}
//
//	*in = (*in)[2:]
//	return length, nil
//}

func (p *Parser) getLeftLength() int {
	totalLength := len(p.currentBuffer)
	return totalLength - p.offset - 1
}

func (p *Parser) parseType() error {
	// check the data length first
	left := p.getLeftLength()
	if left <= 0 {
		p.pauseTask(typeParsing)
		return nil
	}

	// There should always be left bytes when reading types
	next := p.currentBuffer[p.offset]
	p.offset++
	s := string(next)
	switch s {
	case arrays:
		return p.parseArray()
	case bulkString:
		return p.parseBulkString()
	case gedisError:
		return p.parseError()
	case integers:
		return p.parseInteger()
	case simpleString:
		return p.parseSimpleString()
	}

	return unrecognizedTypeErr
}

func (p *Parser) parseString(isError bool) error {
	// check length first
	left := p.getLeftLength()
	if left <= 0 {
		if isError {
			p.pauseTask(errorStringLengthChecking)
		} else {
			p.pauseTask(simpleStringLengthChecking)
		}
		return nil
	}

	var currentNode *Node
	currentPos := p.pos
	p.pos = noneParsing
	if currentPos == simpleStringParsing || currentPos == simpleStringEndChecking ||
		currentPos == errorStringParsing || currentPos == errorStringEndChecking {
		currentNode = p.currentNode
	} else {
		currentNode = &Node{
			data: "",
		}
		if isError {
			currentNode.dataType = TypeError
		} else {
			currentNode.dataType = TypeSimpleString
		}
		p.currentNode = currentNode
	}

	foundR := false
	for left = p.getLeftLength(); left > 0; {
		next := p.currentBuffer[p.offset]
		p.offset++

		s := string(next)
		switch s {
		case "\r":
			{
				if currentPos == simpleStringEndChecking ||
					currentPos == errorStringEndChecking {
					// previous pause stops at '\r', which is not a terminator
					currentNode.data = currentNode.data.(string) + s
				}

				// check if it is time to pause
				if p.getLeftLength() == 0 {
					// we don't know if it is end for now, pause it
					if isError {
						p.pauseTask(errorStringEndChecking)
					} else {
						p.pauseTask(simpleStringEndChecking)
					}
					return nil
				}

				// otherwise, continue the work
				foundR = true
				break
			}
		case "\n":
			{
				if currentPos == simpleStringEndChecking ||
					currentPos == errorStringEndChecking {
					// continue previous pause, and it is confirmed to be terminator
					return nil
				}

				if foundR {
					// found \r in previous iteration
					return nil
				}

				// no \r found, should be part of the data
				currentNode.data = currentNode.data.(string) + string(next)
				break
			}
		default:
			{
				currentNode.data = currentNode.data.(string) + string(next)
				break
			}
		}
	}

	return nil
}

func (p *Parser) parseError() error {
	return p.parseString(true)
}

func (p *Parser) parseSimpleString() error {
	return p.parseString(false)
}

func (p *Parser) parseInteger() error {
	// check length first
	left := p.getLeftLength()
	if left <= 0 {
		p.pauseTask(integerLengthChecking)
		return nil
	}

	var currentNode *Node
	currentPos := p.pos
	p.pos = noneParsing
	if currentPos == integerParsing {
		currentNode = p.currentNode
	} else {
		currentNode = &Node{
			dataType: TypeInteger,
			data:     0,
		}
		p.currentNode = currentNode
	}

	foundR := false
	for left = p.getLeftLength(); left > 0; {
		next := p.currentBuffer[p.offset]
		p.offset++

		s := string(next)
		switch s {
		case "\r":
			{
				if currentPos == integerEndChecking {
					// previous pause stops at '\r', which is not a terminator
					num, err := strconv.Atoi(s)
					if err != nil {
						return err
					}
					currentNode.data = currentNode.data.(int)*10 + num
				}

				// check if it is time to pause
				if p.getLeftLength() == 0 {
					// we don't know if it is end for now, pause it
					p.pauseTask(integerEndChecking)
					return nil
				}

				// otherwise, continue the work
				foundR = true
				break
			}
		case "\n":
			{
				if currentPos == integerEndChecking {
					// continue previous pause, and it is confirmed to be terminator
					return nil
				}

				if foundR {
					// found \r in previous iteration
					return nil
				}

				// no \r found, should be part of the data
				num, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				currentNode.data = currentNode.data.(int)*10 + num
				break
			}
		default:
			{
				num, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				currentNode.data = currentNode.data.(int)*10 + num
				break
			}
		}
	}

	return nil
}

func (p *Parser) parseBulkString() error {


	return nil
}

func (p *Parser) parseAndReturnError() error {
	return nil
}

func (p *Parser) parseArray() error {
	return nil
}

type Parser struct {
	isParsing     bool
	finishParsing bool

	currentBuffer []byte
	offset        int

	currentParsingType DataType
	currentNode        *Node
	pos                pausePos

	returnError      Handler
	returnFatalError Handler
	returnReply      Handler

	result *Node
}

func (p *Parser) Execute(data []byte) error {
	if p.isParsing == false {
		p.Reset()
		return p.executeNewTask(data)
	}

	if p.finishParsing {
		return malFormatErr
	}

	return p.continueTask(data)
}

func (p *Parser) executeNewTask(data []byte) error {
	if p.isParsing {
		return malFormatErr
	}

	p.isParsing = true
	p.currentBuffer = data
	return p.parseType()
}

func (p *Parser) continueTask(data []byte) error {
	if !p.isParsing || p.finishParsing {
		return malFormatErr
	}

	p.currentBuffer = data
	p.offset = 0

	// TODO
	return malFormatErr
}

func (p *Parser) pauseTask(pos pausePos) {
	p.pos = pos
}

func (p *Parser) Reset() {
	p.isParsing = false
	p.finishParsing = false
	p.currentBuffer = nil
	p.offset = 0
	p.currentParsingType = TypeNone
	p.currentNode = nil
}

func (p *Parser) Init(errorHandler Handler, fatalErrorHandler Handler, replyHandler Handler) error {
	if errorHandler == nil || replyHandler == nil {
		return nullHandlerErr
	}

	p.Reset()

	p.returnError = errorHandler
	p.returnReply = replyHandler
	if fatalErrorHandler == nil {
		p.returnFatalError = fatalErrorHandler
	} else {
		p.returnFatalError = errorHandler
	}

	return nil
}
