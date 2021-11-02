package view

import (
	"encoding/binary"
)

/**
 *
 * @author hushengdong
 */
func ViewFile(fileType int, path string, offSet uint64) {

	switch fileType {
	case 1:
		viewConsumeQueueFile(path)
	case 2:
		viewCommitLogFile(path, offSet)
	case 3:
		viewIndexFile(path)
	default:
		panic("不支持的文件类型 可选范围 1:consumergroup文件 2:commitlog 3:index")
	}
}

func bytesToInt64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

func bytesToInt8(buf []byte) uint8 {
	return buf[0]
}

func bytesToInt32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

func bytesToInt16(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}
