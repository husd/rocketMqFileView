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

	p := flag.String("f", "/Users/hushengdong/store/index/20211101165452503", "文件地址")
	fileType := flag.Int("t", 3, "文件的类型 1:consumergroup文件 2:commitlog 3:index")
	count := flag.Int("n", 10, "文件较大 -n指定最多打印多少条记录 ")
	flag.Parse()
	fmt.Println("------------- start -----------------path:", *p, " fileType:", *fileType)
	view.ViewFile(*fileType, *p, *count)
}
