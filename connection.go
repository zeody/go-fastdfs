package go_fastdfs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

const (
	TrackerConnectFlag = 1
	StorageConnectFlag = 2
)

// ConnectionPool,
type ConnectionPool struct {
	minConns int
	maxConns int
	// use a chan to imitate pool
	tConns chan *Connection
	sConns chan *Connection
}

// Connection
type Connection struct {
	host       string
	port       string
	connType   int
	conn       net.Conn
	timeout    time.Duration
	createTime time.Time
}
type Header struct {
	pkLen  int64
	cmd    byte
	status byte
}

// ResCheck is to valid when receive package
type ResCheck struct {
	exceptCmd     byte
	exceptBodyLen int64
}

//create a ConnectionPool
func NewConnectionPool(addr []string, minPool int, maxPool int) (*ConnectionPool, error) {
	cp := &ConnectionPool{
		minConns: minPool,
		maxConns: maxPool,
		tConns:   make(chan *Connection, maxPool),
		sConns:   make(chan *Connection, maxPool),
	}
	//初始化最小connectionPool
	for index := range addr {
		for i := 0; i < minPool; i++ {
			cp.NewTrackerConnection(addr[index], TrackerConnectFlag)
		}
	}
	return cp, nil
}

//create connection
func (cp *ConnectionPool) NewTrackerConnection(addr string, connType int) {
	con := innerNewConnection(addr, connType)
	trackerAdds = append(trackerAdds, addr)
	cp.tConns <- con
}
func (cp *ConnectionPool) NewStorageConnection(addr string, connType int) {
	con := innerNewConnection(addr, connType)
	storageAdds = append(storageAdds, addr)
	cp.sConns <- con
}
func innerNewConnection(addr string, connType int) (res *Connection) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*30)
	if err != nil {
		//create error
		log.Printf("%+v", err)
		return
	}
	split := strings.Split(addr, ":")
	con := &Connection{
		host:       split[0],
		port:       split[1],
		connType:   connType,
		conn:       conn,
		timeout:    time.Minute,
		createTime: time.Now(),
	}
	return con
}

//send content
func (con *Connection) send(bytes []byte) {
	con.conn.Write(bytes)
}

//receive a header
func (con *Connection) receiveHeader(resCheck ResCheck) (Header, error) {
	buf := make([]byte, FDFS_PROTO_HEADER_LEN)
	var header Header
	_, err := io.ReadFull(con.conn, buf)
	if err != nil {
		log.Printf("%v", err)
		return header, err
	}
	if len(buf) != 10 {
		return header, errors.New("header receive data not equal 10")
	}
	if buf[PROTO_HEADER_CMD_INDEX] != resCheck.exceptCmd {
		return header, errors.New("header recevie cmd not equal exceptCmd")
	}
	if buf[PROTO_HEADER_STATUS_INDEX] != 0 {
		return Header{0, 0, buf[PROTO_HEADER_STATUS_INDEX]}, nil
	}
	header = Header{}
	buffer := bytes.NewBuffer(buf[:FDFS_PROTO_PKG_LEN_SIZE])
	binary.Read(buffer, binary.BigEndian, &header.pkLen)
	if header.pkLen < 0 {
		return header, errors.New("receive body len < 0")
	}
	if resCheck.exceptBodyLen >= 0 && header.pkLen != resCheck.exceptBodyLen {
		return header, errors.New("receive pklen not equal exceptLen")
	}
	cmd, _ := buffer.ReadByte()
	status, _ := buffer.ReadByte()
	header.status = status
	header.cmd = cmd
	return header, nil
}

//receive a body
func (con *Connection) receiveBody(pkLen int64) ([]byte, int64) {
	result := make([]byte, 0, pkLen)
	buf := make([]byte, 1024)
	var total int64
	for {
		n, err := con.conn.Read(buf)
		total += int64(n)
		result = append(result, buf[:n]...)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if total == pkLen {
			break
		}
	}
	return result[:total], total
}

// heart beat
func HeartBeat(con *Connection) bool {
	header := Header{0, FDFS_PROTO_CMD_ACTIVE_TEST, 0}
	byte := packageHeader(header)
	resCheck := ResCheck{TRACKER_PROTO_CMD_RESP, 0}
	con.send(byte)
	receiveHeader, err := con.receiveHeader(resCheck)
	if err != nil {
		log.Print(err)
		return false
	}
	return receiveHeader.status == 0
}
