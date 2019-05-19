package go_fastdfs

import (
	"bufio"
	"io"
)

type CallBacker interface {
	send(connect *Connection)
}
type UploadStream struct {
	reader   *bufio.Reader
	fileSize int64
}

func (up *UploadStream) send(connect *Connection) {
	bytes := make([]byte, up.fileSize)
	_, err := up.reader.Read(bytes)
	connect.conn.Write(bytes)
	if err != nil && err == io.EOF {
		return
	}
}

type UploadBuf struct {
	fileBuf []byte
	offset  int64
	length  int64
}

func (up *UploadBuf) send(connect *Connection) {
	connect.conn.Write(up.fileBuf[up.offset : up.offset+up.length])
}
