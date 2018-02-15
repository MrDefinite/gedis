package resp

import (
	"io"
	"sync"
)

type bufferWriter struct {
	writer io.Writer
	buf    []byte
	mu     sync.Mutex
}

func createBufferWriter(writer io.Writer) *bufferWriter {
	bw := bufferWriter{}
	bw.reset(writer)

	return &bw
}

func (bw *bufferWriter) reset(writer io.Writer) {
	*bw = bufferWriter{
		writer: writer,
	}
	bw.buf = make([]byte, 0)
}

func (bw *bufferWriter) appendBytes(b []byte) {
	bw.mu.Lock()
	bw.buf = append(bw.buf, b...)
	bw.mu.Unlock()
}

func (bw *bufferWriter) appendString(s string) {
	bw.mu.Lock()
	bw.buf = append(bw.buf, s...)
	bw.mu.Unlock()
}

func (bw *bufferWriter) length() int {
	bw.mu.Lock()
	length := len(bw.buf)
	bw.mu.Unlock()
	return length
}

func (bw *bufferWriter) flush() error {
	bw.mu.Lock()
	defer bw.mu.Unlock()

	if len(bw.buf) == 0 {
		return nil
	}

	if _, err := bw.writer.Write(bw.buf); err != nil {
		return err
	}

	bw.buf = bw.buf[:0]
	return nil
}
