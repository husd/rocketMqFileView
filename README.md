# rocketMqFileView
解析rocketMQ的存储文件，例如commitlog index consumerGroup等文件

## 编译运行

编译： go build -o mqFileView husd.com/view/main
添加执行权限： chmod +x mqFileView
执行： ./mqFileView -h 查看帮助信息:

```shell
Usage of ./mqFileView:
  -f string
    	文件地址 (default "/Users/hushengdong/store/commitlog/00000000000000000000")
  -offSet uint
    	只对-t=2 有效，表示读取commitlog的什么位置的数据，这个commitLogOffSet不能随意写，必须是有效的offSet
  -t int
    	文件的类型 1:consumergroup文件 2:commitlog 3:index (default 2)
```

## 解析 commitlog

```shell
    # ./mqFileView -f /Users/hushengdong/store/commitlog/00000000000000000000 -offSet=0 -t=2
    ------------- start -----------------path: /Users/hushengdong/store/commitlog/00000000000000000000  fileType: 2
----------commitlog---------文件数据 absPos: 492
消息的大小: 164
magiccode: daa320a7
bodyCrc: 1032136437
queueId: 3
flag: 0
queueOffSet(consumeQueue队列或者tranStateTable队列中消息的个数): 0
physicalOffSet: 492
sysFlag: 0
bornTimeStamp: 1635821124970
bornHost: 10.2.144.15:57443
storeTimestamp: 1635821124972
storeHost: 10.2.144.15:10911
reconsumeTimes: 0
preparedTransactionOffset: 0
bodyLength: 16
body: Hello RocketMQ 3
topicLength: 15
topic: hushengdong_uat
propertiesLength: 42
properties: KEYS1234563WAITtrueTAGS3abcdeewefwef

// -f 表示读取的commitlog的文件，-offSet表示从哪个位置开始读，一般你可以不填，先写0，然后看看totalSize是多少，下次就可以使用
// -offSet=totalSize 就可以查看下一条消息了，这个程序仅仅是研究学习使用的，需要你对offSet的作用有一个直观的认识
```

## 读取consumequeue

```shell
./mqFileView -f /Users/hushengdong/store/consumequeue/hushengdong_uat/0/00000000000000000000 -t=1
./mqFileView -f /Users/hushengdong/store/consumequeue/hushengdong_uat/1/00000000000000000000 -t=1
./mqFileView -f /Users/hushengdong/store/consumequeue/hushengdong_uat/2/00000000000000000000 -t=1
./mqFileView -f /Users/hushengdong/store/consumequeue/hushengdong_uat/3/00000000000000000000 -t=1
./mqFileView -f /Users/hushengdong/store/consumequeue/hushengdong_uat/4/00000000000000000000 -t=1

// 输出结果:

------------- start -----------------path: /Users/hushengdong/store/consumequeue/hushengdong_uat/1/00000000000000000000  fileType: 1
commitLogOffset:164 size:164 messageTagHashCode:18446744072091510419
commitLogOffset:820 size:164 messageTagHashCode:1503314071
commitLogOffset:1476 size:164 messageTagHashCode:329702043
commitLogOffset:2141 size:167 messageTagHashCode:18446744072502135556
commitLogOffset:2809 size:167 messageTagHashCode:1913939208
commitLogOffset:3477 size:167 messageTagHashCode:18446744072583382945
commitLogOffset:4145 size:167 messageTagHashCode:1995186597
commitLogOffset:4813 size:167 messageTagHashCode:821574569
commitLogOffset:5481 size:167 messageTagHashCode:2076433986
commitLogOffset:6149 size:167 messageTagHashCode:902821958
commitLogOffset:6817 size:167 messageTagHashCode:18446744071572265695
commitLogOffset:7485 size:167 messageTagHashCode:984069347
commitLogOffset:8153 size:167 messageTagHashCode:18446744073520008935
commitLogOffset:8821 size:167 messageTagHashCode:1065316736
commitLogOffset:9489 size:167 messageTagHashCode:18446744073601256324
commitLogOffset:10157 size:167 messageTagHashCode:1146564125
commitLogOffset:10825 size:167 messageTagHashCode:18446744073682503713
commitLogOffset:11493 size:167 messageTagHashCode:18446744072508891685
commitLogOffset:12161 size:167 messageTagHashCode:54199486
commitLogOffset:12829 size:167 messageTagHashCode:18446744072590139074
commitLogOffset:13497 size:167 messageTagHashCode:135446875
commitLogOffset:14165 size:167 messageTagHashCode:18446744072671386463
commitLogOffset:14833 size:167 messageTagHashCode:2083190115
commitLogOffset:15501 size:167 messageTagHashCode:18446744072752633852
commitLogOffset:16169 size:167 messageTagHashCode:18446744071579021824
commitLogOffset:16840 size:170 messageTagHashCode:740543188
commitLogOffset:17520 size:170 messageTagHashCode:18446744073276482776
commitLogOffset:18200 size:170 messageTagHashCode:18446744072102870748
commitLogOffset:18880 size:170 messageTagHashCode:18446744073357730165
commitLogOffset:19560 size:170 messageTagHashCode:18446744072184118137
commitLogOffset:20240 size:170 messageTagHashCode:18446744073438977554
commitLogOffset:20920 size:170 messageTagHashCode:18446744072265365526
commitLogOffset:21600 size:170 messageTagHashCode:1677169178
commitLogOffset:22280 size:170 messageTagHashCode:18446744072346612915
```
如果输出结果为空，那么可能是consumequeue中没有数据

## 读取indexFile

索引文件目前解析的结果，没有太大的参考含义，看代码可能会更好点，对go不熟悉就算了。
着重理解下表头的一些统计信息,索引呢文件明细的timestamp=0,是因为存储的是相对于索引文件的第一个消息的偏移量

```shell
./mqFileView -f /Users/hushengdong/store/index/20211102104524952 -t=3

# 输出结果
------------- start -----------------path: /Users/hushengdong/store/index/20211102104524952  fileType: 3
--------------IndexHeader-----------------
8位 该索引文件的第一个消息(Message)的存储时间(落盘时间): 1635821124941
8位 该索引文件的最后一个消息(Message)的存储时间(落盘时间): 1635821125612
8位 该索引文件第一个消息(Message)的在CommitLog(消息存储文件)的物理位置偏移量(可以通过该物理偏移直接获取到该消息): 0
8位 该索引文件最后一个消息(Message)的在CommitLog(消息存储文件)的物理位置偏移量: 22280
4位 该索引文件目前的hash slot的个数: 134
4位 索引文件目前的索引个数: 135
---index文件明细数据: slotNum: 271004  hashcode: 1955271004  commitLogOffSet: 1640  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271005  hashcode: 1955271005  commitLogOffSet: 1807  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271006  hashcode: 1955271006  commitLogOffSet: 1974  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271007  hashcode: 1955271007  commitLogOffSet: 2141  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271008  hashcode: 1955271008  commitLogOffSet: 2308  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271009  hashcode: 1955271009  commitLogOffSet: 2475  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271010  hashcode: 1955271010  commitLogOffSet: 2642  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271011  hashcode: 1955271011  commitLogOffSet: 2809  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271012  hashcode: 1955271012  commitLogOffSet: 2976  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271013  hashcode: 1955271013  commitLogOffSet: 3143  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271035  hashcode: 1955271035  commitLogOffSet: 3310  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271036  hashcode: 1955271036  commitLogOffSet: 3477  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271037  hashcode: 1955271037  commitLogOffSet: 3644  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271038  hashcode: 1955271038  commitLogOffSet: 3811  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271039  hashcode: 1955271039  commitLogOffSet: 3978  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271040  hashcode: 1955271040  commitLogOffSet: 4145  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271041  hashcode: 1955271041  commitLogOffSet: 4312  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271042  hashcode: 1955271042  commitLogOffSet: 4479  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271043  hashcode: 1955271043  commitLogOffSet: 4646  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271044  hashcode: 1955271044  commitLogOffSet: 4813  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271066  hashcode: 1955271066  commitLogOffSet: 4980  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271067  hashcode: 1955271067  commitLogOffSet: 5147  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271068  hashcode: 1955271068  commitLogOffSet: 5314  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271069  hashcode: 1955271069  commitLogOffSet: 5481  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271070  hashcode: 1955271070  commitLogOffSet: 5648  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271071  hashcode: 1955271071  commitLogOffSet: 5815  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271072  hashcode: 1955271072  commitLogOffSet: 5982  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271073  hashcode: 1955271073  commitLogOffSet: 6149  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271074  hashcode: 1955271074  commitLogOffSet: 6316  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271075  hashcode: 1955271075  commitLogOffSet: 6483  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271097  hashcode: 1955271097  commitLogOffSet: 6650  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271098  hashcode: 1955271098  commitLogOffSet: 6817  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271099  hashcode: 1955271099  commitLogOffSet: 6984  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271100  hashcode: 1955271100  commitLogOffSet: 7151  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271101  hashcode: 1955271101  commitLogOffSet: 7318  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271102  hashcode: 1955271102  commitLogOffSet: 7485  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271103  hashcode: 1955271103  commitLogOffSet: 7652  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271104  hashcode: 1955271104  commitLogOffSet: 7819  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271105  hashcode: 1955271105  commitLogOffSet: 7986  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271106  hashcode: 1955271106  commitLogOffSet: 8153  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271128  hashcode: 1955271128  commitLogOffSet: 8320  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271129  hashcode: 1955271129  commitLogOffSet: 8487  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271130  hashcode: 1955271130  commitLogOffSet: 8654  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271131  hashcode: 1955271131  commitLogOffSet: 8821  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271132  hashcode: 1955271132  commitLogOffSet: 8988  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271133  hashcode: 1955271133  commitLogOffSet: 9155  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271134  hashcode: 1955271134  commitLogOffSet: 9322  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271135  hashcode: 1955271135  commitLogOffSet: 9489  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271136  hashcode: 1955271136  commitLogOffSet: 9656  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271137  hashcode: 1955271137  commitLogOffSet: 9823  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271159  hashcode: 1955271159  commitLogOffSet: 9990  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271160  hashcode: 1955271160  commitLogOffSet: 10157  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271161  hashcode: 1955271161  commitLogOffSet: 10324  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271162  hashcode: 1955271162  commitLogOffSet: 10491  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271163  hashcode: 1955271163  commitLogOffSet: 10658  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271164  hashcode: 1955271164  commitLogOffSet: 10825  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271165  hashcode: 1955271165  commitLogOffSet: 10992  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271166  hashcode: 1955271166  commitLogOffSet: 11159  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271167  hashcode: 1955271167  commitLogOffSet: 11326  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271168  hashcode: 1955271168  commitLogOffSet: 11493  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271190  hashcode: 1955271190  commitLogOffSet: 11660  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271191  hashcode: 1955271191  commitLogOffSet: 11827  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271192  hashcode: 1955271192  commitLogOffSet: 11994  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271193  hashcode: 1955271193  commitLogOffSet: 12161  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271194  hashcode: 1955271194  commitLogOffSet: 12328  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271195  hashcode: 1955271195  commitLogOffSet: 12495  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271196  hashcode: 1955271196  commitLogOffSet: 12662  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271197  hashcode: 1955271197  commitLogOffSet: 12829  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271198  hashcode: 1955271198  commitLogOffSet: 12996  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271199  hashcode: 1955271199  commitLogOffSet: 13163  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271221  hashcode: 1955271221  commitLogOffSet: 13330  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271222  hashcode: 1955271222  commitLogOffSet: 13497  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271223  hashcode: 1955271223  commitLogOffSet: 13664  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271224  hashcode: 1955271224  commitLogOffSet: 13831  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271225  hashcode: 1955271225  commitLogOffSet: 13998  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271226  hashcode: 1955271226  commitLogOffSet: 14165  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271227  hashcode: 1955271227  commitLogOffSet: 14332  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271228  hashcode: 1955271228  commitLogOffSet: 14499  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271229  hashcode: 1955271229  commitLogOffSet: 14666  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271230  hashcode: 1955271230  commitLogOffSet: 14833  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271252  hashcode: 1955271252  commitLogOffSet: 15000  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271253  hashcode: 1955271253  commitLogOffSet: 15167  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271254  hashcode: 1955271254  commitLogOffSet: 15334  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271255  hashcode: 1955271255  commitLogOffSet: 15501  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271256  hashcode: 1955271256  commitLogOffSet: 15668  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271257  hashcode: 1955271257  commitLogOffSet: 15835  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271258  hashcode: 1955271258  commitLogOffSet: 16002  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271259  hashcode: 1955271259  commitLogOffSet: 16169  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271260  hashcode: 1955271260  commitLogOffSet: 16336  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 271261  hashcode: 1955271261  commitLogOffSet: 16503  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859028  hashcode: 483859028  commitLogOffSet: 16670  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859029  hashcode: 483859029  commitLogOffSet: 16840  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859030  hashcode: 483859030  commitLogOffSet: 17010  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859031  hashcode: 483859031  commitLogOffSet: 17180  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859032  hashcode: 483859032  commitLogOffSet: 17350  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859033  hashcode: 483859033  commitLogOffSet: 17520  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859034  hashcode: 483859034  commitLogOffSet: 17690  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859035  hashcode: 483859035  commitLogOffSet: 17860  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859036  hashcode: 483859036  commitLogOffSet: 18030  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859037  hashcode: 483859037  commitLogOffSet: 18200  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859059  hashcode: 483859059  commitLogOffSet: 18370  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859060  hashcode: 483859060  commitLogOffSet: 18540  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859061  hashcode: 483859061  commitLogOffSet: 18710  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859062  hashcode: 483859062  commitLogOffSet: 18880  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859063  hashcode: 483859063  commitLogOffSet: 19050  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859064  hashcode: 483859064  commitLogOffSet: 19220  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859065  hashcode: 483859065  commitLogOffSet: 19390  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859066  hashcode: 483859066  commitLogOffSet: 19560  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859067  hashcode: 483859067  commitLogOffSet: 19730  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859068  hashcode: 483859068  commitLogOffSet: 19900  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859090  hashcode: 483859090  commitLogOffSet: 20070  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859091  hashcode: 483859091  commitLogOffSet: 20240  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859092  hashcode: 483859092  commitLogOffSet: 20410  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859093  hashcode: 483859093  commitLogOffSet: 20580  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859094  hashcode: 483859094  commitLogOffSet: 20750  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859095  hashcode: 483859095  commitLogOffSet: 20920  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859096  hashcode: 483859096  commitLogOffSet: 21090  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859097  hashcode: 483859097  commitLogOffSet: 21260  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859098  hashcode: 483859098  commitLogOffSet: 21430  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859099  hashcode: 483859099  commitLogOffSet: 21600  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859121  hashcode: 483859121  commitLogOffSet: 21770  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859122  hashcode: 483859122  commitLogOffSet: 21940  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859123  hashcode: 483859123  commitLogOffSet: 22110  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 3859124  hashcode: 483859124  commitLogOffSet: 22280  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 4663396  hashcode: 629663396  commitLogOffSet: 1476  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 4663397  hashcode: 629663397  commitLogOffSet: 1312  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 4663398  hashcode: 629663398  commitLogOffSet: 1148  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 4663399  hashcode: 629663399  commitLogOffSet: 984  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 4663400  hashcode: 629663400  commitLogOffSet: 820  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 4663401  hashcode: 629663401  commitLogOffSet: 656  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 4663402  hashcode: 629663402  commitLogOffSet: 492  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 4663403  hashcode: 629663403  commitLogOffSet: 328  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 4663404  hashcode: 629663404  commitLogOffSet: 164  timestamp: 0  prevIndexRead: 0
---index文件明细数据: slotNum: 4663405  hashcode: 629663405  commitLogOffSet: 0  timestamp: 0  prevIndexRead: 0
-------end slotCount: 134  链表节点数量: 134
```