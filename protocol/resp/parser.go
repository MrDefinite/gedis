package resp

import (
	"errors"
	"github.com/MrDefinite/gedis/basicdata"
	"io"
	"strconv"
)

const (
	// 512 MB
	MaxDataSize = 1024 * 1024 * 512
	// 16 MB
	MaxDataSizeReadPerTime = 1024 * 1024 * 16
)

var (
	ErrMalFormatInt      = errors.New("cannot parse integer from non gedis object encoded object")
	ErrMalFormatNodeType = errors.New("cannot parse with wrong node type")
)

/**
For simple string, bulk string and error, data is string type
For integer, data is int type
For array, data is a list containing more nodes
*/
type Node struct {
	dataType DataType
	data     interface{}
}

type ParsedResult struct {
	data *Node
}

type Parser struct {
	br *bufferReader
}

func CreateNewParser(reader io.Reader) *Parser {
	parser := Parser{}
	parser.br = createBufferReader(reader)

	return &parser
}

func (p *Parser) parse() (*ParsedResult, error) {
	// let's begin the parse work
	nodes, err := p.parseType()
	if err != nil {
		return nil, err
	}
	return &ParsedResult{data: nodes}, nil
}

func (p *Parser) parseType() (*Node, error) {
	if p.br.unread() < 1 {
		err := p.br.require(1)
		if err != nil {
			return nil, err
		}
	}

	b := p.br.buf[p.br.offset]
	p.br.offset += 1

	switch b {
	case simpleString:
		return p.parseSimpleString()
	case gedisError:
		return p.parseErrorString()
	case integers:
		return p.parseIntegerAsString()
	case bulkString:
		return p.parseBulkString()
	case arrays:
		return p.parseArray()
	}

	return nil, ErrMalFormat
}

func (p *Parser) parseString(isError bool) (*Node, error) {
	data, err := p.br.readLineBytes()
	if err != nil {
		return nil, err
	}

	strObj := basicdata.CreateStringObjectWithBytes(data)
	if isError {
		return &Node{dataType: TypeError, data: strObj}, nil
	}
	return &Node{dataType: TypeSimpleString, data: strObj}, nil
}

func (p *Parser) parseSimpleString() (*Node, error) {
	return p.parseString(false)
}

func (p *Parser) parseErrorString() (*Node, error) {
	return p.parseString(true)
}

func (p *Parser) parseIntegerAsString() (*Node, error) {
	data, err := p.br.readLineBytes()
	if err != nil {
		return nil, err
	}

	// treat it as string, someone can transform it later
	intObj := basicdata.CreateStringObjectWithBytes(data)
	return &Node{dataType: TypeInteger, data: intObj}, nil
}

func (p *Parser) parseBulkString() (*Node, error) {
	// Parse the length first
	lengthNode, err := p.parseIntegerAsString()
	if err != nil {
		return nil, err
	}

	bulkLength, err := p.parseStringObjToInt(lengthNode.data)
	if err != nil {
		return nil, err
	}

	// Ensure the data is there
	err = p.br.require(bulkLength)
	if err != nil {
		return nil, err
	}

	bulkBytes, err := p.br.readLineBytes()
	strObj := basicdata.CreateStringObjectWithBytes(bulkBytes)
	return &Node{dataType: TypeBulkString, data: strObj}, nil
}

func (p *Parser) parseArray() (*Node, error) {
	// First parse the length of the array
	lengthNode, err := p.parseIntegerAsString()
	if err != nil {
		return nil, err
	}

	arrayLength, err := p.parseStringObjToInt(lengthNode.data)
	if err != nil {
		return nil, err
	}

	node := Node{
		dataType: TypeArray,
	}
	node.data = make([]*Node, 0)
	for i := 0; i < arrayLength; i++ {
		newNode, err := p.parseType()
		if err != nil {
			return nil, err
		}

		switch nodePtr := node.data.(type) {
		case []*Node:
			node.data = append(nodePtr, newNode)
			break
		default:
			return nil, ErrMalFormatNodeType
		}
	}
	return &node, nil
}

func (p *Parser) parseStringObjToInt(obj interface{}) (int, error) {
	var length int
	switch objPtr := obj.(type) {
	case *basicdata.GedisObject:
		intStr, err := basicdata.GetStringValueFromObject(objPtr)
		if err != nil {
			return 0, err
		}
		length, err = strconv.Atoi(intStr)
		if err != nil {
			return 0, err
		}
	default:
		return 0, ErrMalFormatInt
	}

	return length, nil
}
