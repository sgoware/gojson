# gojson

[![build](https://img.shields.io/badge/build-1.01-brightgreen)](https://github.com/StellarisW/StellarisW)[![go-version](https://img.shields.io/badge/go-%3E%3D1.8-30dff3?logo=go)](https://github.com/StellarisW/StellarisW)[![OpenTracing Badge](https://img.shields.io/badge/OpenTracing-enabled-blue.svg)](http://opentracing.io)[![Go Report Card](https://goreportcard.com/badge/github.com/emirpasic/gods)](https://goreportcard.com/report/github.com/emirpasic/gods)[![PyPI](https://img.shields.io/badge/License-BSD_2--Clause-green.svg)](https://github.com/emirpasic/gods/blob/master/LICENSE)

> 一个强大的json框架

# 💡  简介

gojson是一个支持数据多种方式读取,智能解析,操作便捷的一个json框架

# 🚀 功能

- `json`序列与反序列化
- `json`指定字段查找
- 支持`json`数据操作 (插入,删除,排序)
- 支持发起请求并从`response body`获取json数据

# 🌟 亮点

- 性能出色
    - `json`对象相对于其他项目携带的数据(字段)更少,加快底层数据传递的速度
    - 使用优化后的递归函数,`struct`映射更快

- 功能强大
    - 支持将任意形式`(嵌套,指针,切片,数组,map,空接口)`的结构体`unmarshal`成`json`格式
    - 支持将其他格式`(yaml,toml,xml...)`的数据转化成`json`格式
    - 支持从文件读取数据
    
- 操作便捷

    

# ⚙ 代码结构

<details>
<summary>展开查看</summary>
<pre><code>
    ├── internal          		   (内部工具包)
    	├── conv                   (数据转换)
    		├── byte.go
    		├── map.go
    		├── string.go
    	├── mutex                  (读写锁)
    		├── mutes.go
    	├── regex                  (正则匹配)
    		├── regex.go 
    	├── types                  (类型操作)
    		├── interface.go 
    ├── const.go                   (常量定义)
    ├── err.go                     (错误定义)
    ├── gojson.go                  (用户可操作函数)
    ├── load.go                    (数据加载相关的函数)
    ├── option.go                  (选项相关的函数)             
</pre></code>
</details>

# 📌 TODO

- [ ] json的序列化
    - [x] string,[]byte的序列化
            - [x] json格式
            - [ ] 其他类型的格式(toml,yaml,xml,ini)
            - [x] 结构体序列化
            - [ ] 切片,数据序列化
            - [ ] map序列化
- [ ] json的反序列化(将json数据unmarshal到结构体)
- [ ] json的数据操作
    - [ ] 查找
    - [ ] 插入
    - [ ] 删除
    - [ ] 排序

- [ ] 发起http请求并读取

# 🛠 环境要求

- golang 版本 >= 1.18
- 

# 🎬 开始



# 📊 性能测试



# 📔 参考文献

[CSDN Golang自定义结构体转map](https://blog.csdn.net/pyf09/article/details/111027686?ops_request_misc=%257B%2522request%255Fid%2522%253A%2522165856096916782395381810%2522%252C%2522scm%2522%253A%252220140713.130102334.pc%255Fall.%2522%257D&request_id=165856096916782395381810&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~all~first_rank_ecpm_v1~pc_rank_34-7-111027686-null-null.142^v33^pc_rank_34,185^v2^control&utm_term=go%20%E7%BB%93%E6%9E%84%E4%BD%93%E8%BD%AC%E6%8D%A2%E6%88%90map%5Bstring%5Dinterface%7B%7D&spm=1018.2226.3001.4187)

[GitHub structs](https://github.com/fatih/structs/)

# 🎈 结语



# 🔑 JetBrains 开源证书支持

`gojson` 项目一直以来都是在 JetBrains 公司旗下的 GoLand 集成开发环境中进行开发，基于 **free JetBrains Open Source license(s)** 正版免费授权，在此表达我的谢意。

<a href="https://www.jetbrains.com/?from=gnet" target="_blank"><img src="https://raw.githubusercontent.com/panjf2000/illustrations/master/jetbrains/jetbrains-variant-4.png" width="250" align="middle"/></a>