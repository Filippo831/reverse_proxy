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

/*
this function sends a chunk once it reaches the defined size ("chunk_size" parameter)
or if a timeout ends ("chunk_timeout" parameter)
*/

func (chunkWriter *ChunkWriterStruct) Write(p []byte) (int, error) {
	bytesWritten := 0

	// while the payload to send has something in it run the function
	for len(p) > 0 {
		// get the length of the payload and if exceed the chunk_size write chunk_size size to bytes to write
		bytesToWrite := len(p)
		if bytesToWrite > chunkWriter.chunkSize {
			bytesToWrite = chunkWriter.chunkSize
		}

		// write only the amount of byte defined in bytesToWrite
		nBytes, err := chunkWriter.writer.Write(p[:bytesToWrite])

		// increment the bytes written counter
		bytesWritten += nBytes

		// delete the written bytes from the payload
		p = p[nBytes:]

		if err != nil {
			log.Print(err)
			return bytesWritten, err
		}

		// if the amount of byte written on the chunk writer exceed the chunkSize, flush it
		if chunkWriter.writer.Buffered() >= chunkWriter.chunkSize {
			chunkWriter.writer.Flush()
		}
	}
	return bytesWritten, nil

}

func (chunkWriter *ChunkWriterStruct) Flush() {
	chunkWriter.writer.Flush()
	if chunkWriter.flushTimer != nil {
		// every time a flush happen, reset the timer
		chunkWriter.flushTimer.Stop()
		chunkWriter.flushTimer = nil
	}
}

func (chunkWriter *ChunkWriterStruct) StartTimer() {
	// when the timer end counting, flush the data
	chunkWriter.flushTimer = time.AfterFunc(chunkWriter.timeout, func() {
		chunkWriter.Flush()
	})
}
