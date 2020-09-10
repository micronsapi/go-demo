package main

import (
	"go-demo/lib"
	"log"
)
/**
 * 微信登陆过程DEMO，只做演示使用
 *
 * @version  1.0
 * @date 2020-9-10
 * @author Ethan
 */
const (
	ADDR string ="8.129.188.188" //TCP地址
	PORT string ="8010" //TCP端口
	HOST string="http://8.129.188.188:8000" //HTTP请求地址
	TOKEN string ="" //购买的CODE
)


func main() {
	//建立TCP连接
	handShakeChan :=make(chan int)
	go lib.NewTcpClient(ADDR,PORT,TOKEN,handShakeChan)
<-handShakeChan
	//初始化HTTP请求
	httpMethod := lib.HttpInitialize(TOKEN,HOST)
	//查看是否在登陆状态
	ok,err :=httpMethod.GetLoginStatus()
	if err!=nil{
		log.Println("err:",err.Error())
	}
	if ok{
		//已经登陆
		log.Println("[INFO] 您已是登陆状态，无需再登陆")
	}else{
		//初始化实例
		log.Println("[DEBUG] 初始化运行实例")
		result,err:=httpMethod.Init()
		if err!=nil{
			log.Println("err:",err.Error())
		}
		log.Println("[DEBUG] ",result)

		//二维码登陆
		log.Println("[DEBUG] 初始化运行实例")
		qr,err :=httpMethod.GetQr()
		if err!=nil{
			log.Println("err:",err.Error())
		}
		log.Println("[DEBUG] 您的登陆二维码（请使用浏览器打开) : ",qr)
	}

	select {
		case <-lib.ExitChan:
			log.Println("系统退出")
			return
	}
}