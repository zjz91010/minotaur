# Minotaur

[![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/kercylan98/minotaur)
![](https://img.shields.io/badge/Email-kercylan@gmail.com-green.svg?style=flat)
![](https://komarev.com/ghpvc/?username=kercylan98)
<a target="_blank" href="https://goreportcard.com/report/github.com/kercylan98/minotaur"><img src="https://goreportcard.com/badge/github.com/kercylan98/minotaur?style=flat-square" /></a>

Minotaur 是一个基于Golang 1.20 编写的服务端开发支持库，其中采用了大量泛型设计，用于游戏服务器开发。

## 目录结构概况
```mermaid
mindmap
  root((Minotaur))
    /configuration 配置管理功能
    /game 游戏通用功能
      /builtin 游戏通用功能内置实现
    /notify 通知功能接口定义
    /planner 策划相关工具目录
      /pce 配置导表功能实现
    /server 网络服务器支持
      /cross 内置跨服功能实现
      /router 内置路由器功能实现
    /utils 工具结构函数目录
    /examples 示例代码目录
```

## Server 架构预览
![server-gdi.jpg](.github/images/server-gdi.jpg)

## 安装
注意：依赖于 **[Go](https://go.dev/) 1.20 +**

运行以下 Go 命令来安装软件包：`minotaur`
```sh
$ go get -u github.com/kercylan98/minotaur
```

## 用法
- 在`Minotaur`中大量使用了 **[泛型](https://go.dev/doc/tutorial/generics)** 、 **[观察者(事件)](https://www.runoob.com/design-pattern/observer-pattern.html)** 和 **[选项模式](https://juejin.cn/post/6844903729313873927)**，在使用前建议先进行相应了解；
- 项目文档可访问 **[pkg.go.dev](https://pkg.go.dev/github.com/kercylan98/minotaur)** 进行查阅；

### 本地文档
可使用 `godoc` 搭建本地文档服务器
#### 安装 godoc
```shell
git clone golang.org/x/tools
cd tools/cmd
go install ...
```
#### 使用 `godoc` 启动本地文档服务器
```shell
godoc -http=:9998 -play
```
#### Windows
```shell
.\local-doc.bat
```

#### Linux or MacOS
```shell
chmod 777 ./local-doc.sh
./local-doc.sh
```

#### 文档地址
- **[http://localhost:9998/pkg/github.com/kercylan98/minotaur/](http://localhost:9998/pkg/github.com/kercylan98/minotaur/)**
- **[https://pkg.go.dev/github.com/kercylan98/minotaur](https://pkg.go.dev/github.com/kercylan98/minotaur)**

### 简单回响服务器
创建一个基于`Websocket`创建的单线程回响服务器。
```go
package main

import (
	"github.com/kercylan98/minotaur/server"
)

func main() {
	srv := server.New(server.NetworkWebsocket)
	srv.RegConnectionReceivePacketEvent(func(srv *server.Server, conn *server.Conn, packet server.Packet) {
		conn.Write(packet)
	})
	if err := srv.Run(":9999"); err != nil {
		panic(err)
	}
}
```
访问 **[WebSocket 在线测试](http://www.websocket-test.com/)** 进行验证。
> Websocket地址: ws://127.0.0.1:9999

### 分流服务器
分流服务器可以将客户端分流到不同的分组上，每个分组中为串行处理，不同分组之间并行处理。
```go
package main

import "github.com/kercylan98/minotaur/server"

func main() {
	srv := server.New(server.NetworkWebsocket,
		server.WithShunt(func(guid int64) chan *server.Message {
			return make(chan *server.Message, 1024*100)
		}, func(conn *server.Conn) (guid int64, allowToCreate bool) {
			guid, allowToCreate = conn.GetData("roomId").(int64)
			return
		}),
	)
	srv.RegConnectionReceivePacketEvent(func(srv *server.Server, conn *server.Conn, packet []byte) {
		conn.Write(packet)
	})
	if err := srv.Run(":9999"); err != nil {
		panic(err)
	}
}
```
> 该示例中假设各房间互不干涉，故通过`server.WithShunt`将连接通过`roomId`进行分组，提高并发处理能力。

### 服务器死锁检测
`Minotaur`内置了服务器消息死锁检测功能，可通过`server.WithDeadlockDetect`进行开启。
```go
package main

import (
	"github.com/kercylan98/minotaur/server"
	"time"
)

func main() {
	srv := server.New(server.NetworkWebsocket,
		server.WithDeadlockDetect(time.Second*5),
	)
	srv.RegConnectionReceivePacketEvent(func(srv *server.Server, conn *server.Conn, packet []byte) {
		time.Sleep(10 * time.Second)
		conn.Write(packet)
	})
	if err := srv.Run(":9999"); err != nil {
		panic(err)
	}
}
```
> 在开启死锁检测的时候需要设置一个合理的死锁怀疑时间，该时间内消息没有处理完毕则会触发死锁检测，并打印`WARN`级别的日志输出。

### 计时器
在默认的`server.Server`不会包含计时器功能，可通过`server.WithTicker`进行开启，例如：
```go
package main

import "github.com/kercylan98/minotaur/server"

func main() {
	srv := server.New(server.NetworkWebsocket, server.WithTicker(50, false))
	if err := srv.Run(":9999"); err != nil {
		panic(err)
	}
}
```
也可以通过`timer.GetTicker`获取计时器进行使用，例如：
```go
package main

import (
	"fmt"
	"github.com/kercylan98/minotaur/utils/timer"
	"github.com/kercylan98/minotaur/utils/times"
	"sync"
)

func main() {
	var ticker = timer.GetTicker(10)
	var wait sync.WaitGroup
	wait.Add(3)
	ticker.Loop("LOOP", timer.Instantly, times.Second, timer.Forever, func() {
		fmt.Println("LOOP")
		wait.Done()
	})
	wait.Wait()
}
```

### 持续更新的示例项目
- **[Minotaur-Example](https://github.com/kercylan98/minotaur-example)**

### 参与贡献
请参考 **[CONTRIBUTING.md](CONTRIBUTING.md)** 贡献指南。

# JetBrains OS licenses

`Minotaur` had been being developed with `GoLand` IDE under the **free JetBrains Open Source license(s)** granted by JetBrains s.r.o., hence I would like to express my thanks here.

<a href="https://www.jetbrains.com/?from=minotaur" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png?_gl=1*1vt713y*_ga*MTEzMjEzODQxNC4xNjc5OTY3ODUw*_ga_9J976DJZ68*MTY4ODU0MDUyMy4yMC4xLjE2ODg1NDA5NDAuMjUuMC4w&_ga=2.261225293.1519421387.1688540524-1132138414.1679967850" width="250" align="middle"/></a>