package view

import (
	"encoding/binary"
)

/**
 *
 * @author hushengdong
 */
func ViewFile(fileType int, path string, count int) {

	switch fileType {
	case 1:
		viewConsumerGroupFile(path, count)
	case 2:
		viewCommitLogFile(path, count)
	case 3:
		viewIndexFile(path, count)
	default:
		panic("不支持的文件类型 可选范围 1:consumergroup文件 2:commitlog 3:index")
	}
}

func bytesToInt64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

func bytesToInt32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}
