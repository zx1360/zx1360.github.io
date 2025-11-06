+++
date = '2025-11-05T16:48:41+08:00'
draft = false
title = '初试自建内网穿透'

categories=['编程', '网络']

tags=['go', '内网穿透', '云服务器']

+++

# 自建内网穿透

> 在校园网环境下向互联网暴露自己的http服务器.

## **背景**

### **需求**

希望我的手机能随时通过数据网络访问到我的pc.

我写的自用的安卓软件`TORRID`有一些功能在如果我的pc提供的http服务器能够在公网被随时访问那会方便很多:

- 打卡、随笔等页面产生的用户数据的备份、同步

- 浏览大体积媒体文件时, 通过http访问pc上的资源而非将漫画存于手机本地. (虽然正常来说都是这么做, 但之前由于没有办法随时随地连接到pc图方便看漫画就存在本地了.)

### **环境**

- 一台带有固定公网ip的低配低宽带的云服务器 (甚至还是共享型)
  - 由于云服务器实在是低配置而且宽带低的可怜, 所以它仅作为内网穿透的工具而使用, 仅看中了他的`固定的公网ip`这一优点.

- pc处于全锥型NAT校园网环境 (最宽松的NAT类型).
- 手机使用对称型NAT数据网络 (最严格限制的类型).

校园网环境下的设备只有内网ip, 只有主动向外部网络发起请求的时候才会获得临时的`公网ip+随机端口`

### 实现效果

任何连入互联网的设备都能访问到我位于校园网内的pc提供的http服务.

## 实现

### 实现原理

云服务器运行程序:

- 监听8000端口, 接收pc发来的经过NAT后的ip和端口号.
- 监听8080端口, 响应手机的http请求告知其记录的ip:port信息.

pc端运行程序:

- 在某一端口运行http服务器, 运行数据备份同步, 响应漫画数据等业务数据逻辑. (此时服务器仅运行于pc本地及路由器网络下如寝室网内)
- 通过设置`端口复用`, 同样使用该端口向云服务器的8000端口发送`长TCP`连接并以小于1/2超时时长的间隔持续不断的发送空数据的`心跳包`维持这个连接使从校园网获取到的临时port能够持续维系下去.
- 手机端的应用则可以通过访问云服务器的8080端口`获知pc的网络地址`, 并实现向pc的单向网络通信.

### 云服务器端代码

```go
package main

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// 存储PC的公网NAT地址（IP:端口）
var pcAddr string

func main() {
	// 1. 启动TCP服务，接收PC的长连接并记录其NAT地址
	go func() {
		listener, err := net.Listen("tcp", ":8000")
		if err != nil {
			log.Fatalf("TCP监听失败: %v", err)
		}
		defer listener.Close()

		for {
			// 只处理第一个PC连接（自用场景）
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("接收连接失败: %v", err)
				continue
			}
			// 记录PC的公网NAT地址（conn.RemoteAddr()返回的是NAT映射后的地址）
			pcAddr = conn.RemoteAddr().String()
			log.Printf("PC已连接，NAT地址: %s", pcAddr)

			// 保持连接（读取数据防止连接被关闭，PC会发心跳）
			go func(c net.Conn) {
				defer c.Close()
				c.SetReadDeadline(time.Now().Add(25 * time.Second))
				buf := make([]byte, 1024)
				for {
					_, err := c.Read(buf)
					if err != nil {
						if !strings.Contains(err.Error(), "closed") {
							log.Printf("PC连接断开: %v", err)
						}
						pcAddr = "" // 清空地址
						return
					}
					// 读超时时间设置.
					c.SetReadDeadline(time.Now().Add(45 * time.Second))
				}
			}(conn)
		}
	}()

	// 2. 启动HTTP服务，供手机查询PC的NAT地址
	http.HandleFunc("/get-pc-addr", func(w http.ResponseWriter, r *http.Request) {
		if pcAddr == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("PC未连接"))
			return
		}
		w.Write([]byte(pcAddr)) // 返回格式: "公网IP:端口"
	})

	log.Println("云服务器启动，TCP端口:8000，HTTP查询端口:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

- 云服务器在这里的作用仅相当于`所有人都知道位置`的`留言板`. pc端在此处留下自己的网络地址, 手机端根据留言板上的留言向pc端发送网络请求,
- 我去感觉我这个比喻打得`相当`恰当啊.

### pc端代码

```go
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	serverIP   = "替换为实际云服务器IP" // 替换为实际云服务器IP
	serverPort = 8000              // 云服务器TCP端口（与服务端对应）
	localPort  = 7274              // 本地端口（同时用于HTTP服务和长连接）
	interval   = 20
)

func main() {
	// 1. 启动本地HTTP服务（端口，供手机访问）
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("手机访问PC成功！"))
		})

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
		if err != nil {
			log.Fatalf("本地HTTP服务启动失败: %v", err)
		}
		log.Printf("本地HTTP服务已启动，端口: %d", localPort)
		log.Fatal(http.Serve(listener, mux))
	}()

	// 2. 与云服务器建立长连接并维持NAT映射
	for {
		dialer := net.Dialer{
			LocalAddr: &net.TCPAddr{Port: localPort},
			Timeout:   10 * time.Second,
		}
		conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
		if err != nil {
			log.Printf("连接云服务器失败，5秒后重试: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("与云服务器建立长连接成功，开始发送心跳")

		heartbeatTicker := time.NewTicker(interval * time.Second)
		// 使用标志位控制内层循环，避免直接return导致流程不清晰
		running := true

		// 单独的资源释放函数，确保在各种退出路径下都能正确释放资源
		cleanup := func() {
			heartbeatTicker.Stop()
			conn.Close()
			log.Println("长连接断开，准备重连")
		}

		// 内层循环处理心跳逻辑，使用标志位控制退出
		for running {
			select {
			case <-heartbeatTicker.C:
				fmt.Println("续命成功.")
				_, err := conn.Write([]byte("heartbeat"))
				if err != nil {
					log.Printf("心跳发送失败: %v", err)
					running = false // 设置退出标志
				}
			}
		}

		// 执行资源清理
		cleanup()
		// 等待一小段时间再重连，避免频繁重试
		time.Sleep(2 * time.Second)
	}
}
```

- "启动本地http服务"处做修改承担起实际业务, 响应漫画资源和数据备份同步等.

### 移动端代码

-> 后续有时间了改改我的安卓应用`TORRID`, 利用好这个工具.

## 完毕

**摸索的过程**

很早就有这个想法了, 看到29元/年的云服务器这还说啥了, 我下单就是了.

前后大概一周的时间:

- 稍微了解了`Go`这门语言, 作为服务端语言, 可编译为单一的二进制文件很适合部署到云服务器上.
- 了解了网络相关的入门知识, 了解了`内网穿透`, `p2p`的实现原理.
- 最后接近一整天的时间, 实际上手实现这一过程.

**助手**

豆包用起来的话, 感觉他啥都知道, 但有时候就是一句话都不多讲, 问什么讲什么, 提出什么需求就只解决哪个, 

`比如如下这一过程`:

1. 他告诉我为了使经NAT后的临时端口保持不被销毁, 需要pc端持续不断地向云服务器发送网络请求. 
2. 让他生成响应代码, 他返回了但端口还是会不断变化, 问来问去才知道网络请求是会随机采取本地端口向外连接, 由于NAT映射的四元组规则, 这会直接导致新分配另外一个临时端口. 
3. 让他每次都使用同一本地端口请求, 它返回的代码运行一会儿就崩了, 说是端口占用, 因为他只实现了我提出的**每次使用同一本地端口**的需求, 而甚至不多思考一步**关闭原先端口**
4. ......
5. 后来了解到了`TCP长连接`, `心跳包`等概念, 终于让它生成可用的代码了.

- 如上只是pc端代码的`向云服务器留言`这一功能的摸索过程......

总而言之, 豆包好用, 这一实现借助它才得以完成, 但是摸索的这一过程常常陷入低效的碰壁, 哪天考虑考虑别的响应更快, 思考的更远的AI呢.

**可恶**

![Snipaste_2025-11-06_09-09-34](assets/Snipaste_2025-11-06_09-09-34.png)

令人`佩服`的效率和耐心...

骗你的, 不止六七小时.
