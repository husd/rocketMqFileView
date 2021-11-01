package view

import (
	"fmt"
	"io/ioutil"
)

/**
 *
 * @author hushengdong
 */
func viewIndexFile(path string, count int) {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic("读取index文件失败: " + path)
	}
	// header占40个字节

	//8位 该索引文件的第一个消息(Message)的存储时间(落盘时间)
	//8位 该索引文件的最后一个消息(Message)的存储时间(落盘时间)
	//8位 该索引文件第一个消息(Message)的在CommitLog(消息存储文件)的物理位置偏移量(可以通过该物理偏移直接获取到该消息)
	//8位 该索引文件最后一个消息(Message)的在CommitLog(消息存储文件)的物理位置偏移量
	//4位 该索引文件目前的hash slot的个数
	//4位 索引文件目前的索引个数
	firstMsgTime := bytesToInt64(buf[0:8])
	lastMsgTime := bytesToInt64(buf[8:16])
	firstMsgOffset := bytesToInt64(buf[16:24])
	lastMsgOffset := bytesToInt64(buf[24:32])
	hashSlotCount := bytesToInt32(buf[32:36])
	indexCount := bytesToInt32(buf[36:40])
	// header占40个字节
	fmt.Println("--------------IndexHeader-----------------")
	fmt.Println("8位 该索引文件的第一个消息(Message)的存储时间(落盘时间):", firstMsgTime)
	fmt.Println("8位 该索引文件的最后一个消息(Message)的存储时间(落盘时间):", lastMsgTime)
	fmt.Println("8位 该索引文件第一个消息(Message)的在CommitLog(消息存储文件)的物理位置偏移量(可以通过该物理偏移直接获取到该消息):", firstMsgOffset)
	fmt.Println("8位 该索引文件最后一个消息(Message)的在CommitLog(消息存储文件)的物理位置偏移量:", lastMsgOffset)
	fmt.Println("4位 该索引文件目前的hash slot的个数:", hashSlotCount)
	fmt.Println("4位 索引文件目前的索引个数:", indexCount)

	// slot array ，每个4字节，有50W个，占 4*50W 个字节 这里是根据 key 的hashcode算出来的 ，这里我们遍历下，找出最多10个key的索引
	total := 0
	max := 5000000
	slotNum := 1
	for slotNum < max && total < count {
		// 程序里是根据 hashcode(key) % 5000000 算出来的 这里我们挨个查看
		slotPos := calcSlotPosition(slotNum)
		indexNum := bytesToInt32(buf[slotPos : slotPos+4])
		if indexNum > uint32(0) {
			indexPos := calcIndexPosition(indexNum)
			//读取指定的文件的位置 index 固定占 20个字节
			tempBuf := buf[indexPos : indexPos+20]
			hashCode := bytesToInt32(tempBuf[0:4])
			commitLogOffSet := bytesToInt64(tempBuf[4:12])
			timestamp := bytesToInt32(tempBuf[12:16])
			nextIndexOffSet := bytesToInt32(tempBuf[16:20])
			fmt.Println("---index slot table pos:", slotNum, " hashcode:", hashCode, " commitLogOffSet:", commitLogOffSet, " timestamp:", timestamp, " nextIndexOffSet:", nextIndexOffSet)
			total++
		}
		slotNum++
	}
	//fmt.Println("-------end current:",slotNum)
}

//key-->计算hash值-->hash值对500万取余算出对应的slot序号-->
//根据40+(n-1)*4(公式1)算出该slot在文件中的位置-->读取slot值，也就是index序号-->
// 根据40+500000*4+(s-1)*20(公式2)算出该index在文件中的位置-->读取该index-->
// 将key的hash值以及传入的时间范围与index的keyHash值以及timeDiff值进行比对。

// 第n个slot在indexFile中的起始位置是这样:40+(n-1)*4
func calcSlotPosition(n int) int {

	if n <= 0 {
		panic("无效的slot数量")
	}
	return 40 + (n-1)*4
}

// 第n个index在indexFile中的起始位置是这样: 40+5000000*4+(n-1)*20
func calcIndexPosition(n uint32) int {

	if n <= 0 {
		panic("无效的slot数量")
	}
	return 40 + 5000000*4 + (int(n)-1)*20
}
