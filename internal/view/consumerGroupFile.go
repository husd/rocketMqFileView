package view

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

/**
 *
 * @author hushengdong
 */
func viewConsumerGroupFile(path string, count int) {

	fd, err := os.Open(path)
	defer fd.Close()
	if err != nil {
		panic("读取consumer group 文件失败: " + path)
	}
	r := bufio.NewReader(fd)
	// 1次读20个字节
	buf := make([]byte, 20)
	template := "commitLogOffset:%d size:%d messageTagHashCode:%d "
	c := 0
	for {
		n, err := io.ReadFull(r, buf)
		if err == io.EOF {
			break
		}
		if n != 20 {
			panic("读取错误")
		}
		commitLogOffset := bytesToInt64(buf[0:8])
		size := bytesToInt32(buf[8:12])
		messageTagHashCode := bytesToInt64(buf[12:20])
		fmt.Println(fmt.Sprintf(template, commitLogOffset, size, messageTagHashCode))
		if c >= count {
			break
		}
		c++
	}
}
