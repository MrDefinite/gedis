package resp

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
)

const (
	defaultBufferCap = 1024
)

var (
	ErrDataTooLong = errors.New("Cannot find the trailing /n for data too long")
)

/**
-----------------------------------------------------
....offset......size.......slice len........slice cap
      |          |             |
      v          v             v
    current    size of      length of
    read       total        the slice
    offset     buffer       could be
                            longer than
                            buffer size
-----------------------------------------------------
*/
type bufferReader struct {
	reader io.Reader
	buf    []byte

	offset int
	size   int
}

func createBufferReader(reader io.Reader) *bufferReader {
	br := bufferReader{}
	br.reset(reader)

	return &br
}

func (br *bufferReader) unread() int {
	return br.size - br.offset
}

func (br *bufferReader) unused() int {
	return cap(br.buf) - br.size
}

func (br *bufferReader) require(required int) error {
	if required-br.unused() <= 0 {
		return nil
	}

	// We need to read more now!
	br.compact()
	if br.unused() < required {
		newSlice := make([]byte, br.size+required)
		copy(newSlice, br.buf)
		br.buf = newSlice
	}

	n, err := io.ReadAtLeast(br.reader, br.buf[br.size:], required)
	if err != nil {
		return err
	}
	br.size += n

	return nil
}

// Move the unread bytes to the beginning place of buf
func (br *bufferReader) compact() {
	if br.offset == 0 {
		// Already there
		return
	}

	// Move byte from offset:size to 0:size-offset
	length := br.unread()
	for i := 0; i < length; i++ {
		br.buf[i] = br.buf[br.offset+i]
	}
	br.offset = 0
	br.size = length
}

func (br *bufferReader) fill() error {
	br.compact()
	if br.unused() <= 0 {
		return nil
	}

	n, err := br.reader.Read(br.buf[br.size:])
	if err != nil {
		return err
	}
	br.size += n
	return nil
}

func (br *bufferReader) reset(reader io.Reader) {
	*br = bufferReader{
		reader: reader,
		offset: 0,
		size:   0,
	}
	br.buf = make([]byte, defaultBufferCap)
}

func (br *bufferReader) skip(skipSize int) {
	if br.unread() >= skipSize {
		br.offset += skipSize
	}
}

func (br *bufferReader) readByte() (byte, error) {
	if br.unread() < 1 {
		err := br.require(1)
		if err != nil {
			return 0, err
		}
	}

	b := br.buf[br.offset]
	br.offset += 1
	return b, nil
}

func (br *bufferReader) readNextBytes(size int) ([]byte, error) {
	err := br.require(size)
	if err != nil {
		return nil, err
	}

	data := br.buf[br.offset : br.offset+size]
	br.offset += size
	return data, nil
}

// Not include the trailing '/r/n'
func (br *bufferReader) readLineBytes() ([]byte, error) {
	index := bytes.IndexByte(br.buf, '\n')
	if index <= 0 {
		if err := br.fill(); err != nil {
			return nil, err
		}
	}

	index = bytes.IndexByte(br.buf, '\n')
	if index <= 0 {
		return nil, ErrDataTooLong
	}

	// Check if prev byte is '\r'
	pre := index - 1
	if br.buf[pre] != '\r' {
		return nil, ErrMalFormat
	}

	data := br.buf[br.offset:pre]
	br.offset = index + 1
	return data, nil
}

func (br *bufferReader) readLine() (string, error) {
	data, err := br.readLineBytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}
