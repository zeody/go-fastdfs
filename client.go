package go_fastdfs

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

const DEFAULT_STORAGE_LEN = 10

var (
	trackerAdds []string
	storageAdds []string
)

// FdfsClient is a interface define
type FdfsClient interface {
	getConnect(addr string, connType int) *Connection
}
type Client struct {
	connPool *ConnectionPool
}

func Default() *Client {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//解析文档
	conf := ParseYaml()
	trackerServers := strings.Split(conf.TrackerServers, ",")
	//创建连接池
	connPool, err := NewConnectionPool(trackerServers, len(trackerServers)*20, len(trackerServers)*40)
	if err != nil {
		log.Panic("create connPool error")
	}

	//初始化
	client := Client{connPool}
	infos := client.getStorageInfos()
	for index := range infos {
		str := fmt.Sprintf("%s:%d", infos[index].addr, infos[index].port)
		client.connPool.NewStorageConnection(str, StorageConnectFlag)
	}
	return &client
}

//获取连接
func (client *Client) getConnect(addr string, connType int) *Connection {
	for {
		select {
		case con := <-client.connPool.tConns:
			//check timeout
			if !isTimeout(con.createTime, time.Now(), con.timeout) {
				//simple check
				return con
			} else if HeartBeat(con) {
				//real check
				return con
			} else {
				con.conn.Close()
				//创建新的连接
				client.connPool.NewTrackerConnection(addr, con.connType)
				break
			}
		case con := <-client.connPool.sConns:
			//check timeout
			if !isTimeout(con.createTime, time.Now(), con.timeout) {
				//simple check
				return con
			} else if HeartBeat(con) {
				//real check
				return con
			} else {
				con.conn.Close()
				//创建新的连接
				client.connPool.NewStorageConnection(addr, con.connType)
				break
			}
			//获取连接1秒超时
		case <-time.After(time.Second):
			return nil
		}

	}
}

//get storage info   ----------------------- start
func (client *Client) getStorageInfo() StorageServer {
	header := Header{0, TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE, 0}
	resCheck := ResCheck{TRACKER_PROTO_CMD_RESP, TRACKER_QUERY_STORAGE_STORE_BODY_LEN}
	connect := client.getConnect(trackerAdds[rand.Intn(len(trackerAdds))], TrackerConnectFlag)
	storage := innerGetResByHeader(connect, header, resCheck)
	//return connection
	client.connPool.tConns <- connect
	return ParseStorageServer(storage)
}

//get storage info with groupName
func (client *Client) getStorageInfoWithGroupName(groupName string) StorageServer {
	header := Header{FDFS_GROUP_NAME_MAX_LEN, TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE, 0}
	if groupName == "" {
		log.Panic("groupName is nil")
	}
	connect := client.getConnect(trackerAdds[rand.Intn(len(trackerAdds))], TrackerConnectFlag)
	resCheck := ResCheck{TRACKER_PROTO_CMD_RESP, TRACKER_QUERY_STORAGE_STORE_BODY_LEN}
	storage := innerGetResByHeader(connect, header, resCheck)
	//return connection
	client.connPool.tConns <- connect
	return ParseStorageServer(storage)
}

//get storages info
func (client *Client) getStorageInfos() []StorageServer {
	header := Header{0, TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ALL, 0}
	resCheck := ResCheck{TRACKER_PROTO_CMD_RESP, -1}
	connect := client.getConnect(trackerAdds[rand.Intn(len(trackerAdds))], TrackerConnectFlag)
	storage := innerGetResByHeader(connect, header, resCheck)
	//return connection
	client.connPool.tConns <- connect
	return ParseStoragesServer(storage)
}

//get storages info with groupName
func (client *Client) getStorageInfosWithGroupName(groupName string) []StorageServer {
	header := Header{FDFS_GROUP_NAME_MAX_LEN, TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ALL, 0}
	resCheck := ResCheck{TRACKER_PROTO_CMD_RESP, -1}
	if groupName == "" {
		log.Panic("groupName is nil")
	}
	connect := client.getConnect(trackerAdds[rand.Intn(len(trackerAdds))], TrackerConnectFlag)
	storage := innerGetResByHeader(connect, header, resCheck)
	//return connection
	client.connPool.tConns <- connect
	return ParseStoragesServer(storage)
}

//get storage info --------------------- end

//get storage to use ------------------- start
func (client *Client) getFetchStorage(groupName string, fileName string) StorageServer {
	header := Header{FDFS_GROUP_NAME_MAX_LEN + int64(len(fileName)), TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE, 0}
	body := joinBytes([]byte(groupName), []byte(fileName))
	resCheck := ResCheck{TRACKER_PROTO_CMD_RESP, -1}
	connect := client.getConnect(trackerAdds[rand.Intn(len(trackerAdds))], TrackerConnectFlag)
	storage := innerGetResByBody(connect, header, body, resCheck)
	//return connection
	client.connPool.tConns <- connect
	return ParseStorageServer(storage)
}

func (client *Client) getUpdateStorage(groupName string, fileName string) StorageServer {
	header := Header{FDFS_GROUP_NAME_MAX_LEN + int64(len(fileName)), TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE, 0}
	body := joinBytes([]byte(groupName), []byte(fileName))
	resCheck := ResCheck{TRACKER_PROTO_CMD_RESP, -1}
	connect := client.getConnect(trackerAdds[rand.Intn(len(trackerAdds))], TrackerConnectFlag)
	storage := innerGetResByBody(connect, header, body, resCheck)
	//return connection
	client.connPool.tConns <- connect
	return ParseStorageServer(storage)
}

func (client *Client) getFetchStorages(groupName string, fileName string) []StorageServer {
	header := Header{FDFS_GROUP_NAME_MAX_LEN + int64(len(fileName)), TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE, 0}
	body := joinBytes([]byte(groupName), []byte(fileName))
	resCheck := ResCheck{TRACKER_PROTO_CMD_RESP, -1}
	connect := client.getConnect(trackerAdds[rand.Intn(len(trackerAdds))], TrackerConnectFlag)
	storage := innerGetResByBody(connect, header, body, resCheck)
	//return connection
	client.connPool.tConns <- connect
	return ParseStoragesServer(storage)
}
func (client *Client) getFetchStorageByFileId(fileId string) StorageServer {
	result := splitFileId(fileId)
	groupName := result[0]
	fileName := result[1]
	return client.getFetchStorage(groupName, fileName)
}
func (client *Client) getFetchStoragesByFileId(fileId string) []StorageServer {
	result := splitFileId(fileId)
	groupName := result[0]
	fileName := result[1]
	return client.getFetchStorages(groupName, fileName)
}

//get storage to use ----------------- end
func (client *Client) uploadFile(localFile string, fileExtName string, nameValue map[string]string) []string {
	groupName := ""
	cmd := STORAGE_PROTO_CMD_UPLOAD_FILE
	file, err := os.Open(localFile)

	if err != nil {
		log.Panic("open file error")
	}
	defer file.Close()
	info, e := os.Stat(localFile)
	if e != nil {
		log.Panic("cannot get fileInfo")
	}
	reader := bufio.NewReader(file)
	stream := &UploadStream{reader, info.Size()}
	if fileExtName == "" {
		index := strings.Index(localFile, ".")
		if len(localFile)-index <= int(FDFS_FILE_EXT_NAME_MAX_LEN)+1 {
			s := localFile[index+1:]
			fileExtName = s
		}
	}
	return client.innerUploadFile(cmd, groupName, "", "", fileExtName, info.Size(), stream, nameValue)
}

func (client *Client) uploadFileByte(fileBuff []byte, fileExtName string, nameValue map[string]string) []string {
	groupName := ""
	cmd := STORAGE_PROTO_CMD_UPLOAD_FILE
	loadBuf := &UploadBuf{fileBuff, 0, int64(len(fileBuff))}
	return client.innerUploadFile(cmd, groupName, "", "", fileExtName, int64(len(fileBuff)), loadBuf, nameValue)
}

//upload file implement
func (client *Client) innerUploadFile(cmd byte, groupName string, masterFile string, prefixName string, fileExtName string, fileSize int64, callBack CallBacker, nameValue map[string]string) []string {
	// check param, rules are groupName,masterFile,prefixName,fileExtName --> write or update
	var (
		uploadSlave bool
		storage     StorageServer
		bodyLen     int64
		offset      int
		body        []byte //body [8]byte-->len(masterFile),[8]byte-->fileSize
		//temp buf to save [8]byte
		buf8  = make([]byte, 8)
		buf16 = make([]byte, 16)
		buf6  = make([]byte, 6)
	)
	uploadSlave = groupName != "" && masterFile != "" && prefixName != "" && fileExtName != ""
	if uploadSlave {
		storage = client.getUpdateStorage(groupName, masterFile)
		bodyLen = 2*FDFS_PROTO_PKG_LEN_SIZE + int64(FDFS_FILE_PREFIX_MAX_LEN) + int64(FDFS_FILE_EXT_NAME_MAX_LEN) + int64(len(masterFile)) + fileSize
		//int64 to byte
		binary.BigEndian.PutUint64(buf8, uint64(len(masterFile)))
		offset = len(buf8)
		body = make([]byte, bodyLen)
		copy(body, buf8)
		binary.BigEndian.PutUint64(buf8, uint64(fileSize))
		copy(body[offset:], buf8)
		offset = 16 //2*FDFS_PROTO_PKG_LEN_SIZE
		copy(buf16, []byte(prefixName))
		copy(body[offset:], buf16)
		offset = 32
		copy(buf6, []byte(fileExtName))
		copy(body[offset:], buf6)
		offset = 38
		copy(buf8, []byte(masterFile))
		copy(body[offset:], buf8)
	} else {
		storage = client.getStorageInfo()
		bodyLen = 1 + FDFS_PROTO_PKG_LEN_SIZE + int64(FDFS_FILE_EXT_NAME_MAX_LEN) + fileSize
		offset = 1
		body = make([]byte, bodyLen)
		body[0] = storage.storePath
		binary.BigEndian.PutUint64(buf8, uint64(fileSize))
		copy(body[offset:], buf8)
		offset = 9 //1 + FDFS_PROTO_PKG_LEN_SIZE
		copy(buf6, []byte(fileExtName))
		copy(body[offset:], buf6)
	}
	header := Header{bodyLen, cmd, 0}
	resCheck := ResCheck{STORAGE_PROTO_CMD_RESP, -1}
	connect := client.getConnect(fmt.Sprintf("%s:%d", storage.addr, storage.port), StorageConnectFlag)
	resultBody := innerGetResByBodyCallback(connect, header, body, callBack, resCheck)
	uploadResult := ParseUploadResult(resultBody)
	//return connection
	client.connPool.sConns <- connect
	return uploadResult
}

//getStorage use inner only send header
func innerGetResByHeader(connect *Connection, header Header, resCheck ResCheck) []byte {
	headerBytes := packageHeader(header)
	connect.send(headerBytes)
	receiveHeader, err := connect.receiveHeader(resCheck)
	if err != nil {
		log.Panicf("%v", err)
	}
	bytes, _ := connect.receiveBody(receiveHeader.pkLen)
	return bytes
}

func innerGetResByBody(connect *Connection, header Header, body []byte, resCheck ResCheck) []byte {
	bodyBytes := packageBody(header, body)
	connect.send(bodyBytes)

	receiveHeader, err := connect.receiveHeader(resCheck)
	if err != nil {
		log.Panicf("%v", err)
	}
	bytes, _ := connect.receiveBody(receiveHeader.pkLen)
	return bytes
}

func innerGetResByBodyCallback(connect *Connection, header Header, body []byte, backer CallBacker, resCheck ResCheck) []byte {
	bodyBytes := packageBody(header, body)
	connect.send(bodyBytes)
	//send file stream
	backer.send(connect)
	receiveHeader, err := connect.receiveHeader(resCheck)
	if err != nil {
		log.Panicf("%v", err)
	}
	bytes, _ := connect.receiveBody(receiveHeader.pkLen)
	return bytes
}
