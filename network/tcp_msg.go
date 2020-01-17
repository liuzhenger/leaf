package network

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

/*
 MsgParser 数据协议解析器
 按照指定的数据包规则,解析二进制数据
 数据格式：
 	data 业务数据
 	len 业务数据字节长度（本身是1或2或4个字节存储）
*/

// --------------
// | len | data |
// --------------
type MsgParser struct {
	lenMsgLen    int    // data数据的字节个数用几个字节存储
	minMsgLen    uint32 // data数据最少字节个数
	maxMsgLen    uint32 // data数据最长字节个数
	littleEndian bool   // 字节序列存储方式，大小端
}

func NewMsgParser() *MsgParser {
	p := new(MsgParser)
	p.lenMsgLen = 2
	p.minMsgLen = 1
	p.maxMsgLen = 4096
	p.littleEndian = false

	return p
}

// It's dangerous to call the method on reading or writing
func (p *MsgParser) SetMsgLen(lenMsgLen int, minMsgLen uint32, maxMsgLen uint32) {
	if lenMsgLen == 1 || lenMsgLen == 2 || lenMsgLen == 4 {
		p.lenMsgLen = lenMsgLen
	}
	if minMsgLen != 0 {
		p.minMsgLen = minMsgLen
	}
	if maxMsgLen != 0 {
		p.maxMsgLen = maxMsgLen
	}

	var max uint32
	switch p.lenMsgLen {
	case 1:
		max = math.MaxUint8
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if p.minMsgLen > max {
		p.minMsgLen = max
	}
	if p.maxMsgLen > max {
		p.maxMsgLen = max
	}
}

// It's dangerous to call the method on reading or writing
func (p *MsgParser) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// goroutine safe
func (p *MsgParser) Read(conn *TCPConn) ([]byte, error) {
	var b [4]byte
	bufMsgLen := b[:p.lenMsgLen]

	// read len
	if _, err := io.ReadFull(conn, bufMsgLen); err != nil {
		return nil, err
	}

	// parse len
	var msgLen uint32
	switch p.lenMsgLen {
	case 1:
		msgLen = uint32(bufMsgLen[0])
	case 2:
		if p.littleEndian {
			msgLen = uint32(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = uint32(binary.BigEndian.Uint16(bufMsgLen))
		}
	case 4:
		if p.littleEndian {
			msgLen = binary.LittleEndian.Uint32(bufMsgLen)
		} else {
			msgLen = binary.BigEndian.Uint32(bufMsgLen)
		}
	}

	// check len
	// 客户端发送的数据（data部分）长度要大于等于配置的最小长度（minMsgLen），小于等于配置的最大长度（maxMsgLen）
	if msgLen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil, errors.New("message too short")
	}

	// data
	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return nil, err
	}

	return msgData, nil
}

// goroutine safe
func (p *MsgParser) Write(conn *TCPConn, args ...[]byte) error {
	// get len
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}

	// check len
	if msgLen > p.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return errors.New("message too short")
	}

	msg := make([]byte, uint32(p.lenMsgLen)+msgLen)

	// write len
	switch p.lenMsgLen {
	case 1:
		msg[0] = byte(msgLen)
	case 2:
		if p.littleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case 4:
		if p.littleEndian {
			binary.LittleEndian.PutUint32(msg, msgLen)
		} else {
			binary.BigEndian.PutUint32(msg, msgLen)
		}
	}

	// write data
	l := p.lenMsgLen
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}

	conn.Write(msg)

	return nil
}
