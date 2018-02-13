package resp

import "github.com/pkg/errors"

type DataType int8

const (
	defaultBufferCap = 1024
)

const (
	TypeSimpleString = iota
	TypeBulkString
	TypeError
	TypeInteger
	TypeArray
)

const (
	simpleString = '+'
	gedisError   = '-'
	integers     = ':'
	bulkString   = '$'
	arrays       = '*'
)

const (
	Crlf = "/r/n"
)

var (
	ErrBufferOverflow = errors.New("Parser internal error: read buffer overflow")
	ErrMalFormat      = errors.New("Protocol error: mal-format string received")
	//emptyBufferErr      = errors.New("failed to parse: empty buffer")
	//unrecognizedTypeErr = errors.New("failed to parse: unrecognized byte")
	//illegalTypeErr      = errors.New("failed to parse: illegal data type")
)
