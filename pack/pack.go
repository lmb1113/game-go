package pack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

type Package struct {
	Version        [2]byte // 协议版本
	Length         int16   // 数据部分长度
	MsgType        uint8   // 数据类型
	Timestamp      int64   // 时间戳
	HostnameLength int16   // 主机名长度
	Hostname       []byte  // 主机名
	Msg            []byte  // 数据部分长度
}

func (p *Package) Pack(writer io.Writer) error {
	var err error
	err = binary.Write(writer, binary.BigEndian, &p.Version)
	err = binary.Write(writer, binary.BigEndian, &p.Length)
	err = binary.Write(writer, binary.BigEndian, &p.MsgType)
	err = binary.Write(writer, binary.BigEndian, &p.Timestamp)
	err = binary.Write(writer, binary.BigEndian, &p.HostnameLength)
	err = binary.Write(writer, binary.BigEndian, &p.Hostname)
	err = binary.Write(writer, binary.BigEndian, &p.Msg)
	return err
}

func (p *Package) Unpack(reader io.Reader) error {
	var err error
	err = binary.Read(reader, binary.BigEndian, &p.Version)
	err = binary.Read(reader, binary.BigEndian, &p.Length)
	err = binary.Read(reader, binary.BigEndian, &p.MsgType)
	err = binary.Read(reader, binary.BigEndian, &p.Timestamp)
	err = binary.Read(reader, binary.BigEndian, &p.HostnameLength)
	p.Hostname = make([]byte, p.HostnameLength)
	err = binary.Read(reader, binary.BigEndian, &p.Hostname)
	p.Msg = make([]byte, p.Length-8-1-p.HostnameLength-2)
	err = binary.Read(reader, binary.BigEndian, &p.Msg)
	return err
}
func (p *Package) String() string {
	return fmt.Sprintf(" version:%s length:%d msg_type%d timestamp:%d hostname:%s msg:%s",
		p.Version,
		p.Length,
		p.MsgType,
		p.Timestamp,
		p.Hostname,
		p.Msg,
	)
}

func Send(conn net.Conn, msgType uint8, name string, data []byte) {
	packData := &Package{
		Version:        [2]byte{'V', '1'},
		MsgType:        msgType,
		Timestamp:      time.Now().Unix(),
		HostnameLength: int16(len(name)),
		Hostname:       []byte(name),
		Msg:            data,
	}
	packData.Length = 8 + 2 + 1 + packData.HostnameLength + int16(len(packData.Msg))
	buf := new(bytes.Buffer)
	err := packData.Pack(buf)
	if err != nil {
		return
	}
	conn.Write(buf.Bytes())
}
