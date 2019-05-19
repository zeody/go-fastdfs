package go_fastdfs

import (
	"bytes"
	"encoding/binary"
	"log"
)

// StorageServer show a simple storageServer info
type StorageServer struct {
	addr      string
	port      uint64
	storePath byte
}

func innerParseAddr(buf []byte, storePath byte) StorageServer {
	storageServertemp := StorageServer{}
	storageServertemp.addr = string(bytes.Trim(buf[:FDFS_IPADDR_SIZE-1], "\x00"))
	storageServertemp.port = binary.BigEndian.Uint64(buf[FDFS_IPADDR_SIZE-1 : FDFS_IPADDR_SIZE+7])
	storageServertemp.storePath = storePath
	return storageServertemp
}

// ParseStorageServer response body like below
// byte 0-15 -->groupName len is 16
// byte 16-30 -->server len is 15
// byte 31-38 -->port len is 8
// byte 39 --> storePath  len is 1
func ParseStorageServer(buf []byte) StorageServer {
	storageServer := StorageServer{}
	storageServer.addr = string(bytes.Trim(buf[FDFS_GROUP_NAME_MAX_LEN:FDFS_GROUP_NAME_MAX_LEN+FDFS_IPADDR_SIZE-1], "\x00"))
	storageServer.port = binary.BigEndian.Uint64(buf[FDFS_GROUP_NAME_MAX_LEN+FDFS_IPADDR_SIZE-1 : FDFS_GROUP_NAME_MAX_LEN+FDFS_IPADDR_SIZE+7])
	storageServer.storePath = buf[FDFS_GROUP_NAME_MAX_LEN+FDFS_IPADDR_SIZE+7]
	return storageServer
}

// ParseStoragesServer response body like below
// byte 0-15---->groupName len is 16
// byte 16-30----> server1 len is 15
// byte 31-38----> port len is 8
// byte 39 ---->storePath len is 1
// byte 40-54----> server2 len is 15 remember to trim result.
// byte 55-62----> port len is 8
// .....
func ParseStoragesServer(buf []byte) []StorageServer {
	bufLen := int64(len(buf))
	ipAddrsLen := bufLen - FDFS_GROUP_NAME_MAX_LEN - 1
	recordLen := FDFS_IPADDR_SIZE - 1 + FDFS_PROTO_PKG_LEN_SIZE
	if bufLen < TRACKER_QUERY_STORAGE_FETCH_BODY_LEN {
		log.Panic("get body buf error")
	}
	if ipAddrsLen%(recordLen) != 0 {
		log.Panic("ip addrLen buf error")
	}
	num := ipAddrsLen / (recordLen)
	if num > 16 {
		log.Panic("server count large than 16")
	}
	storePath := buf[FDFS_GROUP_NAME_MAX_LEN+FDFS_IPADDR_SIZE+7]
	storageServers := make([]StorageServer, 0, num)
	offset := FDFS_GROUP_NAME_MAX_LEN
	for i := 0; int64(i) < num; i++ {
		storageServers = append(storageServers, innerParseAddr(buf[offset:offset+FDFS_IPADDR_SIZE+FDFS_PROTO_PKG_LEN_SIZE], storePath))
		offset = offset + FDFS_IPADDR_SIZE + FDFS_PROTO_PKG_LEN_SIZE
	}
	return storageServers
}

func ParseUploadResult(buf []byte) []string {
	groupName := string(buf[:FDFS_GROUP_NAME_MAX_LEN-1])
	remoteFileName := string(buf[FDFS_GROUP_NAME_MAX_LEN-1:])
	return []string{groupName, remoteFileName}
}
