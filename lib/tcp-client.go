//@Title:  TCP客户端
//@author: Ethan
//@date: 2020/9/10
//@Description:
//
package lib

import (
	"fmt"
	"github.com/wondayer/tcpx"
	"log"
	"net"
	"time"
)

const (
	TCP_MSGID_HANDSHAKE_REQ int32 = 1 //与服务器握手
	TCP_MSGID_HANDSHAKE_RESP int32 = 101 //握手数据响应
	TCP_MSGID_SYS_RESP int32 = 102 //系统推送响应
	TCP_MSGID_PUSH_RESP int32 = 103 //微信业务数据响应
	TCP_MSGID_HEARTBEAT int32 = 1392 //心跳包 10s
)


var ExitChan =make(chan bool) //异常退出channel

//创建一个tcp客户端
//token 购买的CODE
//addr tcp地址
//port tcp端口
func NewTcpClient(addr ,port ,token string,handShakeChan chan<- int){
	log.Println("[INFO] 创建连接TCP")

	//建立新的tcp连接
	Conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		log.Println("[ERROR] TCP连接异常：err:",err.Error())
		return
	}
	//发送握手消息
	sendbuf, _ := tcpx.PackWithMarshaller(tcpx.Message{
		MessageID: TCP_MSGID_HANDSHAKE_REQ,
		Body:      []byte(token),
	})
	_, _ = Conn.Write(sendbuf)
	//开始监听
	for{
		//接收服务器推送消息
		messageID, body, err := tcpx.UnPackFromReader(Conn)
		if err != nil {
			tcpExceptionHandle(err)
			return
		}
		switch messageID {
			//握手成功
			case TCP_MSGID_HANDSHAKE_RESP:
				hanshakeHandle(string(body),Conn)
				handShakeChan<-1
			//系统推送
			case TCP_MSGID_SYS_RESP:
				systemMsgHandle(string(body),Conn)
			//微信业务
			case TCP_MSGID_PUSH_RESP:
				wechatMsgHandle(string(body),Conn)
		}
	}
}

//握手成功处理
func hanshakeHandle(body string,conn net.Conn){
	log.Println("[INFO] tcp握手成功")
	//创建心跳包
	go (func() {
		for  {
			_, _ = conn.Write(tcpx.PackHeartbeat())
			time.Sleep( 10*time.Second )
		}
	})()
}

//系统推送
func systemMsgHandle(body string,conn net.Conn){
	log.Println("[DEBUG] TCP_MSGID_SYS_RESP :",body)
}

//微信推送消息
func wechatMsgHandle(body string,conn net.Conn){
	log.Println("[DEBUG] TCP_MSGID_PUSH_RESP :",body)
}

//tcp异常处理
func tcpExceptionHandle(err error){
	log.Println("[ERROR] tcp连接异常，err：",err.Error())
	ExitChan<-true
}