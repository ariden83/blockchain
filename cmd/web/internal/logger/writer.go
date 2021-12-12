package logger

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"fmt"
	"io"
	"net"
)

// See http://docs.graylog.org/en/2.4/pages/gelf.html.
const (
	// ChunkSize sets to maximal chunk size in bytes.
	chunkSize = 8192
	// ChunkDataSize is ChunkSize minus header's size.
	chunkDataSize = chunkSize - 12
	// MaxChunkCount maximal chunk per message count.
	maxChunkCount = 128
)

var (
	// ChunkedMagicBytes chunked message magic bytes.
	// See http://docs.graylog.org/en/2.4/pages/gelf.html.
	chunkedMagicBytes = []byte{0x1e, 0x0f}
)

// GELFWriter is a writter following the GELF specs.
// It will write in udp using a gzip compression.
type GELFWriter struct {
	conn net.Conn
}

// NewGELFWriter will create a new GELFWriter.
// We expect a valid udp address.
func NewGELFWriter(addr string) (*GELFWriter, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}

	return &GELFWriter{conn: conn}, nil
}

// Sync implements WriteSyncer.
func (w *GELFWriter) Sync() error {
	return nil
}

// Write implements io.Writer.
func (w *GELFWriter) Write(buf []byte) (int, error) {
	b, err := w.compress(buf)
	if err != nil {
		return 0, err
	}

	if count := w.chunkCount(b); count > 1 {
		return w.writeChunked(count, b)
	}

	n, err := w.conn.Write(b)
	if err != nil {
		return n, err
	}

	if n != len(b) {
		return n, fmt.Errorf("writed %d bytes but should %d bytes", n, len(b))
	}

	return n, nil
}

func (w *GELFWriter) compress(b []byte) ([]byte, error) {
	var buffer bytes.Buffer

	writer := gzip.NewWriter(&buffer)
	if _, err := writer.Write(b); err != nil {
		return make([]byte, 0), err
	}
	if err := writer.Close(); err != nil {
		return make([]byte, 0), err
	}

	out := buffer.Bytes()
	return out, nil
}

// chunkCount calculate the number of GELF chunks.
func (w *GELFWriter) chunkCount(b []byte) int {
	lenB := len(b)
	if lenB <= chunkSize {
		return 1
	}

	return len(b)/chunkDataSize + 1
}

// writeChunked send message by chunks.
func (w *GELFWriter) writeChunked(count int, b []byte) (n int, err error) {
	if count > maxChunkCount {
		return 0, fmt.Errorf("need %d chunks but shold be later or equal to %d", count, maxChunkCount)
	}

	// Generate random messageID
	messageID := make([]byte, 8)
	if n, err = io.ReadFull(rand.Reader, messageID); err != nil || n != 8 {
		return 0, fmt.Errorf("rand.Reader: %d/%s", n, err)
	}

	bytesLeft := len(b)
	buffer := bytes.NewBuffer(make([]byte, 0, chunkSize))

	for i := 0; i < count; i++ {
		off := i * chunkDataSize
		chunkLen := chunkDataSize
		if chunkLen > bytesLeft {
			chunkLen = bytesLeft
		}

		buffer.Reset()
		buffer.Write(chunkedMagicBytes)
		buffer.Write(messageID)
		buffer.WriteByte(uint8(i))
		buffer.WriteByte(uint8(count))
		buffer.Write(b[off : off+chunkLen])

		if n, err = w.conn.Write(buffer.Bytes()); err != nil {
			return len(b) - bytesLeft + n, err
		}

		if n != len(buffer.Bytes()) {
			n = len(b) - bytesLeft + n
			return n, fmt.Errorf("writed %d bytes but should %d bytes", n, len(b))
		}

		bytesLeft -= chunkLen
	}

	if bytesLeft != 0 {
		return len(b) - bytesLeft, fmt.Errorf("error: %d bytes left after sending", bytesLeft)
	}

	return len(b), nil
}
