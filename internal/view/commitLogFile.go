package view

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

/**
 * 解析 commitLog文件
 * @author hushengdong
 */

type commitLog struct {
	totalSize   uint32 // 代表这个消息的大小
	magicCode   uint32 // MAGICCODE = daa320a7
	bodyCrc     uint32 // 消息体BODY CRC 当broker重启recover时会校验
	queueId     uint32 //
	flag        uint32 //
	queueOffSet uint64 // 这个值是个自增值不是真正的consume queue的偏移量，
	// 可以代表这个consumeQueue队列或者tranStateTable队列中消息的个数，若是非事务消息或者commit事务消息，
	// 可以通过这个值查找到consumeQueue中数据，QUEUEOFFSET * 20才是偏移地址；若是PREPARED或者Rollback事务，
	// 则可以通过该值从tranStateTable中查找数据
	physicalOffSet uint64 // 代表消息在commitLog中的物理起始地址偏移量
	sysFlag        uint32 // 指明消息是事物事物状态等消息特征，二进制为四个字节从右往左数：当4个字节均为0（值为0）时表示非事务消息；
	// 当第1个字节为1（值为1）时表示表示消息是压缩的（Compressed）；
	// 当第2个字节为1（值为2）表示多消息（MultiTags）；
	// 当第3个字节为1（值为4）时表示prepared消息；
	// 当第4个字节为1（值为8）时表示commit消息；
	// 当第3/4个字节均为1时（值为12）时表示rollback消息；
	// 当第3/4个字节均为0时表示非事务消息；
	bomTimestamp   uint64 // 消息产生端(producer)的时间戳
	bomHost        uint64 // 消息产生端(producer)地址(address:port)
	storeTimestamp uint64 // 消息在broker存储时间

	storeHostAddress uint64 // 消息存储到broker的地址(address:port)

	reconsumeTimes            uint64 // 消息被某个订阅组重新消费了几次（订阅组之间独立计数）,因为重试消息发送到了topic名字为%retry%groupName的队列queueId=0的队列中去了，成功消费一次记录为0；
	preparedTransactionOffset uint64 // 表示是prepared状态的事物消息
	bodyLength                uint32 // 消息体大小值
	body                      string // 消息体内容
	topicLength               uint8  // topic名称内容大小
	topic                     string // topic的内容值
	propertiesLength          uint16 // 属性值大小
	properties                string // 属性值
}

var _commitLogIndex uint64 = 0

func viewCommitLogFile(path string, absPos uint64) {

	//messageMagicCode := 0xdaa320a7
	//blankmagiccode := 0xcbd43194
	//
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic("读取index文件失败: " + path)
	}
	totalSize := readInt32(&buf, absPos)
	magicCode := readInt32(&buf, absPos)
	bodyCrc := readInt32(&buf, absPos)
	queueId := readInt32(&buf, absPos)
	flag := readInt32(&buf, absPos)
	queueOffSet := readInt64(&buf, absPos)
	physicalOffSet := readInt64(&buf, absPos)
	sysFlag := readInt32(&buf, absPos)
	bornTimeStamp := readInt64(&buf, absPos)

	bornIpByte := readByteByLen(&buf, absPos, uint64(4))
	bornIp := convertByte2Ip(bornIpByte)
	bornPort := readInt32(&buf, absPos)
	bornHost := fmt.Sprintf("%s:%d", bornIp, bornPort)

	storeTimestamp := readInt64(&buf, absPos)

	storeIpByte := readByteByLen(&buf, absPos, uint64(4))
	storeIp := convertByte2Ip(storeIpByte)
	storePort := readInt32(&buf, absPos)
	storeHost := fmt.Sprintf("%s:%d", storeIp, storePort)

	reconsumeTimes := readInt32(&buf, absPos)
	preparedTransactionOffset := readInt64(&buf, absPos)
	bodyLength := readInt32(&buf, absPos)
	body := readStringByLength(&buf, absPos, uint64(bodyLength))
	topicLength := readInt8(&buf, absPos)
	topic := readStringByLength(&buf, absPos, uint64(topicLength))
	propertiesLength := readInt16(&buf, absPos)
	properties := readStringByLength(&buf, absPos, uint64(propertiesLength))
	fmt.Println("----------commitlog---------文件数据 absPos:", absPos)
	fmt.Println("totalSize:", totalSize)
	fmt.Println("magiccode:", strconv.FormatInt(int64(magicCode), 16))
	fmt.Println("bodyCrc:", bodyCrc)
	fmt.Println("queueId:", queueId)
	fmt.Println("flag:", flag)
	fmt.Println("queueOffSet(consumeQueue队列或者tranStateTable队列中消息的个数):", queueOffSet)
	fmt.Println("physicalOffSet:", physicalOffSet)
	fmt.Println("sysFlag:", sysFlag)
	fmt.Println("bornTimeStamp:", bornTimeStamp)
	fmt.Println("bornHost:", bornHost)
	fmt.Println("storeTimestamp:", storeTimestamp)
	fmt.Println("storeHost:", storeHost)
	fmt.Println("reconsumeTimes:", reconsumeTimes)
	fmt.Println("preparedTransactionOffset:", preparedTransactionOffset)
	fmt.Println("bodyLength:", bodyLength)
	fmt.Println("body:", body)
	fmt.Println("topicLength:", topicLength)
	fmt.Println("topic:", topic)
	fmt.Println("propertiesLength:", propertiesLength)
	fmt.Println("properties:", properties)
}

func convertByte2Ip(ipByte []byte) string {

	length := len(ipByte)
	if length == 4 {
		// ipv4
		return convertByte2IpV4(ipByte)
	} else if length == 16 {
		// ipv6
		return convertByte2IpV6(ipByte)
	}
	return ""
}

func convertByte2IpV6(ipByte []byte) string {

	return "暂时不支持ipv6"
}

func convertByte2IpV4(addr []byte) string {

	template := "%d.%d.%d.%d"
	return fmt.Sprintf(template, addr[0], addr[1], addr[2], addr[3])
}

func readStringByLength(buf *[]byte, absPos uint64, len uint64) string {

	start := absPos + _commitLogIndex
	end := start + len
	res := string((*buf)[start:end])
	_commitLogIndex = _commitLogIndex + len
	return res
}

func readByteByLen(buf *[]byte, absPos uint64, len uint64) []byte {

	start := absPos + _commitLogIndex
	end := start + len
	_commitLogIndex = _commitLogIndex + len
	return (*buf)[start:end]
}

func readInt32(buf *[]byte, absPos uint64) uint32 {

	start := absPos + _commitLogIndex
	end := start + 4
	res := bytesToInt32((*buf)[start:end])
	_commitLogIndex = _commitLogIndex + 4
	return res
}

func readInt8(buf *[]byte, absPos uint64) uint8 {

	start := absPos + _commitLogIndex
	end := start + 1
	res := bytesToInt8((*buf)[start:end])
	_commitLogIndex = _commitLogIndex + 1
	return res
}

func readInt64(buf *[]byte, absPos uint64) uint64 {

	start := absPos + _commitLogIndex
	end := start + 8
	res := bytesToInt64((*buf)[start:end])
	_commitLogIndex = _commitLogIndex + 8
	return res
}

func readInt16(buf *[]byte, absPos uint64) uint16 {

	start := absPos + _commitLogIndex
	end := start + 2
	res := bytesToInt16((*buf)[start:end])
	_commitLogIndex = _commitLogIndex + 2
	return res
}
