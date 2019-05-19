package go_fastdfs

import (
	"bytes"
	"encoding/binary"
	"log"
	"strings"
	"time"
)

const (
	SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR = "/"
)

//check timeout
func isTimeout(start time.Time, end time.Time, timeout time.Duration) bool {
	return start.Add(timeout).After(end)
}

//split fileId to []string
func splitFileId(fileId string) []string {
	index := strings.Index(fileId, SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR)
	if index <= 0 || index == len(fileId)-1 {
		log.Panic("fileId is not valid")
	}
	result := strings.Split(fileId, SPLIT_GROUP_NAME_AND_FILENAME_SEPERATOR)
	return result
}

//convert byte to int64
func byteToInt64(bytes []byte) int64 {
	return int64(binary.BigEndian.Uint64(bytes))
}

//convert int64 to byte
func int64ToByte(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func packageHeader(header Header) []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, header.pkLen)
	buffer.WriteByte(header.cmd)
	buffer.WriteByte(header.status)
	return buffer.Bytes()
}
func packageBody(header Header, body []byte) []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, header.pkLen)
	buffer.WriteByte(header.cmd)
	buffer.WriteByte(header.status)
	buffer.Write(body)
	return buffer.Bytes()
}

func joinBytes(pByte ...[]byte) []byte {
	buffer := new(bytes.Buffer)
	for i := range pByte {
		buffer.Write(pByte[i])
	}
	return buffer.Bytes()
}
