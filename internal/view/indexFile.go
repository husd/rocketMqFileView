package view

import (
	"fmt"
	"io/ioutil"
)

/**
 * 读取 rocketmq的 index 文件
 * @author hushengdong
 */
func viewIndexFile(path string) {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic("读取index文件失败: " + path)
	}
	// header占40个字节
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
	slotCount := 0 // 这个是存在索引的slot的数量
	max := hashSlotNum
	slotNum := 0

	slotListNodeCount := 0
	for slotNum < max {
		// 程序里是根据 hashcode(key) % 5000000 算出来的 这里我们挨个查看
		absSlotPos := calcSlotPosition(slotNum)
		slotValue := bytesToInt32(buf[absSlotPos : absSlotPos+hash_slot_size])
		if slotValue > uint32(0) {
			// 读取指定的文件的位置 index 固定占 20个字节
			for nextIndexToRead := slotValue; ; {
				absIndexPos := calcIndexPosition(nextIndexToRead)
				hashCode := bytesToInt32(buf[absIndexPos : absIndexPos+4])
				commitLogOffSet := bytesToInt64(buf[absIndexPos+4 : absIndexPos+4+8])
				timeDiff := bytesToInt32(buf[absIndexPos+4+8 : absIndexPos+4+8+4]) // timeDiff偏移量是秒 这个细节要注意一下
				prevIndexRead := bytesToInt32(buf[absIndexPos+4+8+4 : absIndexPos+4+8+4+4])

				//absTimestamp := firstMsgTime + uint64(timeDiff) * uint64(1000) // 索引文件里是 相对于第一个消息的时间的偏移量 单位是秒，所以要乘1000 还原为毫秒
				fmt.Println("---index文件明细数据: slotNum:", slotNum, " hashcode:", hashCode, " commitLogOffSet:", commitLogOffSet, " timestamp:", timeDiff, " prevIndexRead:", prevIndexRead)
				slotListNodeCount++
				nextIndexToRead = prevIndexRead
				if prevIndexRead <= 0 {
					break
				}
			}
			slotCount++
		}
		slotNum++
	}
	fmt.Println("-------end slotCount:", slotCount, " 链表节点数量:", slotListNodeCount)
}

// 第n个slot在indexFile中的起始位置是这样:40+(n-1)*4
func calcSlotPosition(slotNum int) int {

	return index_header_size + slotNum*hash_slot_size
}

// 第n个index在indexFile中的起始位置是这样: 40+5000000*4+(n-1)*20
func calcIndexPosition(nextIndexToRead uint32) int {

	return index_header_size + hashSlotNum*hash_slot_size + int(nextIndexToRead)*index_size
}

const index_header_size int = 40
const hashSlotNum int = 5000000
const hash_slot_size int = 4
const index_size int = 20
