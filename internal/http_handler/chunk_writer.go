package http_handler

import (
	"bufio"
	"io"
	"log"
	"time"

	readconfiguration "github.com/Filippo831/reverse_proxy/internal/read_configuration"
)

type ChunkWriterStruct struct {
	writer     *bufio.Writer
	flushTimer *time.Timer
	chunkSize  int
	timeout    time.Duration
}

func ChunkedWriter(w io.Writer, conf readconfiguration.Server) *ChunkWriterStruct {
	chunkWriter := &ChunkWriterStruct{
		writer:    bufio.NewWriterSize(w, conf.ChunkSize*1024),
		chunkSize: conf.ChunkSize * 1024,
		timeout:   time.Duration(conf.ChunkTimeout * int(time.Millisecond)),
	}
	return chunkWriter
}

func (chunkWriter *ChunkWriterStruct) Write(p []byte) (int, error) {
	bytesWritten := 0

	for len(p) > 0 {
		bytesToWrite := len(p)
		if bytesToWrite > chunkWriter.chunkSize {
			bytesToWrite = chunkWriter.chunkSize
		}

		nBytes, err := chunkWriter.writer.Write(p[:bytesToWrite])

		bytesWritten += nBytes

		p = p[nBytes:]

		if err != nil {
			log.Print(err)
			return bytesWritten, err
		}

		if chunkWriter.writer.Buffered() >= chunkWriter.chunkSize {
			chunkWriter.writer.Flush()
		}
	}
	return bytesWritten, nil

}

func (chunkWriter *ChunkWriterStruct) Flush() {
	chunkWriter.writer.Flush()
	if chunkWriter.flushTimer != nil {
		chunkWriter.flushTimer.Stop()
		chunkWriter.flushTimer = nil
	}
}

func (chunkWriter *ChunkWriterStruct) StartTimer() {
	chunkWriter.flushTimer = time.AfterFunc(chunkWriter.timeout, func() {
		chunkWriter.Flush()
	})
}
