package proto

// 自定义数据包协议

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

// Encode msg
func Encode(msg string) ([]byte, error) {
	// 读取数据的长度, 转换成int32类型(占用4个字节)
	var length = int32(len(msg))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}

	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(msg))
	if err != nil {
		return nil, err
	}

	return pkg.Bytes(), nil
}

// Decode msg
func Decode(reader *bufio.Reader) (string, error) {
	// 读取消息头 获取body length
	lengthByte, _ := reader.Peek(4) // 读取前4个字节的数据 encode 的时候 头部是int32 是四个字节的
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	// buffer 返回缓冲中现有的可读数据的字节数
	if int32(reader.Buffered()) < length + 4 {
		return "", err
	}

	// 读取真正的消息数据
	pack := make([]byte, int(length + 4))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}