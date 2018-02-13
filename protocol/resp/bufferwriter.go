package resp

import "io"

type bufferWriter struct {
	writer io.Writer
	buf    [1024]byte
}
