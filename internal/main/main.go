package main

import (
	"flag"
	"fmt"
	"husd.com/view/view"
)

/**
 *
 * @author hushengdong
 */
func main() {

	p := flag.String("f", "/Users/hushengdong/store/commitlog/00000000000000000000", "文件地址")
	fileType := flag.Int("t", 2, "文件的类型 1:consumequeue文件 2:commitlog 3:index")
	offSet := flag.Uint64("offSet", uint64(0), "只对-t=2 有效，表示读取commitlog的什么位置的数据，这个commitLogOffSet不能随意写，必须是有效的offSet")
	flag.Parse()
	fmt.Println("------------- start -----------------path:", *p, " fileType:", *fileType)
	view.ViewFile(*fileType, *p, *offSet)
}
