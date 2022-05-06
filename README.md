# qymkv

## 说明
本项目有部分源码来自于godis，感谢godis作者的开源贡献
https://github.com/HDT3213/godis

## 学习与借鉴
这是一些我们项目中借鉴的资料：
黄健宏《redis设计与实现》
https://www.cnblogs.com/Finley/category/1598973.html
https://github.com/HDT3213/godis
https://space.bilibili.com/1324259795/channel/seriesdetail?sid=642777
https://hardcore.feishu.cn/mindnotes/bmncn1pO2ZhEyFkBgbQ2ttXncsc

## 开发了解
首先需要了解redis的各种命令的功能，因为redis的命令比较成熟，所以我们的功能基本对标他即可。
然后就是需要熟悉本项目各个源码包，比较重要的就是dict/concurrent.go。
然后阅读下面的实现，了解各种开发需求。
在开发的过程中，可以边了解上面所发的redis资料，了解什么是内存kv，为什么要这样设计？
然后如果在了解过程中有很想实现的功能而我没有提到的，可以在群里进行讨论，大家看看是否有实现的可能。

## 实现
首先需要完成database的一些存储相关设计并实现以配置文件redis.toml启动项目
然后实现dict、list、sds、set这四个对象的部分操作并完成测试（具体如何操作可以查看database包和datastruct包）

目前不打算对接redis-cli，因为如果使用redis-cli的话，一切格式都得按其规定来传输，这样就被限制了

因此这些命令需要对外提供多个接口来操作不同的对象：
这些接口会放在operation.go中,如：
```go
type StringOperation interface {
	Set(key, val string) *Reply
}
```
这时就需要定义一个对外的struct（假设是 StringOp）来实现这个interface，如果想要这种请求： set x=1，
那自然就需要调用 StringOp.Set(key,val)来进行set操作(具体看network package里的内容)
其他对象也类似如此

完成上面部分可以考虑进行sortedset的开发以及各种其他功能的扩展，比如key的超时时间之类的...