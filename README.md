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

## 实现
首先需要完成database的一些存储相关设计并实现以配置文件redis.toml启动项目
然后实现dict、list、sds、set这四个对象的部分操作并完成测试
然后可以转入连接阶段，目前可以考虑与redis-cli对接，也可以直接用网页方式，也可以两种一起（待定）
然后需要实现aof持久化

完成上面部分可以考虑进行sortedset的开发
然后就是一些容器的打包和对应环境的配置