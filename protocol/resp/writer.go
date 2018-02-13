package resp

import (
	"io"
	"strconv"
)

type Writer struct {
	bw *bufferWriter
}

func CreateNewWriter(writer io.Writer) *Writer {
	w := Writer{}
	w.bw = createBufferWriter(writer)

	return &w
}

func (w *Writer) Write() error {
	err := w.bw.flush()
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) AppendSimpleString(s string) {
	w.bw.appendString(string(simpleString))
	w.bw.appendString(s)
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendSimpleStringBytes(b []byte) {
	w.bw.appendString(string(simpleString))
	w.bw.appendBytes(b)
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendInteger(i int) {
	w.bw.appendString(string(integers))
	w.bw.appendString(strconv.Itoa(i))
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendIntegerString(s string) {
	w.bw.appendString(string(integers))
	w.bw.appendString(s)
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendIntegerBytes(b []byte) {
	w.bw.appendString(string(integers))
	w.bw.appendBytes(b)
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendError(s string) {
	w.bw.appendString(string(gedisError))
	w.bw.appendString(s)
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendErrorBytes(b []byte) {
	w.bw.appendString(string(gedisError))
	w.bw.appendBytes(b)
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendBulkString(s string) {
	w.bw.appendString(string(bulkString))
	w.AppendInteger(len(s))
	w.bw.appendString(s)
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendBulkStringBytes(b []byte) {
	w.bw.appendString(string(bulkString))
	w.AppendInteger(len(b))
	w.bw.appendBytes(b)
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendArrayLength(i int) {
	w.bw.appendString(string(arrays))
	w.bw.appendString(strconv.Itoa(i))
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendArrayLengthString(s string) {
	w.bw.appendString(string(arrays))
	w.bw.appendString(s)
	w.bw.appendString(Crlf)
}

func (w *Writer) AppendArrayLengthBytes(b []byte) {
	w.bw.appendString(string(arrays))
	w.bw.appendBytes(b)
	w.bw.appendString(Crlf)
}
