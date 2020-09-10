//@Title:  HTTP请求
//@author: Ethan
//@date: 2020/9/10
//@Description:
//
package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

const (
	POST ="POST"
	GET ="GET"
)


type HttpClient struct{
	Token string
	Host string
}

type ResponseData struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
	Msg string `json:"msg"`
}

var ins *HttpClient
var once sync.Once

//初始化http请求
func HttpInitialize(token ,host string) *HttpClient {
	once.Do(func(){
			 ins = &HttpClient{
				 Token: token,
				 Host:  host,
			 }
		 })
	 return ins
 }


//发送http请求
func (req *HttpClient)httpRequest(method, url string, buffer []byte) ([]byte, error) {
	//提交请求
	reqest, err := http.NewRequest(method, req.Host+url, bytes.NewReader(buffer))
	if err != nil {
		return nil, err
	}
	reqest.Close = true
	////增加header选项
	reqest.Header.Add("Token", req.Token)
	reqest.Header.Add("Content-Type", "application/json")

	//处理返回结果
	response, err := http.DefaultClient.Do(reqest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Http Request Failed!,status=%v",response.StatusCode))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//获取微信登陆状态
//[GET] /api/v1/login/islogin
func (req *HttpClient) GetLoginStatus()(bool,error) {
	body,err :=req.httpRequest(GET,"/api/v1/login/islogin",nil)
	if err!=nil{
		return false,err
	}
	res,ok := req.responseHandle(body)
	if ok{
		return res.(bool),nil
	}
	return false,nil
}

//初始化实例
//[GET] /api/v1/login/init
func(req *HttpClient) Init() (string,error) {
	body,err :=req.httpRequest(GET,"/api/v1/login/init",nil)
	if err!=nil{
		return "",err
	}
	res,_ := req.responseHandle(body)
	return res.(string),nil
}

//获取登陆二维码
func (req *HttpClient)GetQr() (string,error){
	body,err :=req.httpRequest(GET,"/api/v1/login/loginqrcode",nil)
	if err!=nil{
		return "",err
	}
	res,_ := req.responseHandle(body)
	return res.(string),nil
}


//解析响应内容
func(req *HttpClient)responseHandle(body []byte) (interface{},bool){
	res := &ResponseData{}
	json.Unmarshal(body,&res)
	if res.Code == 0{
		return res.Data,true
	}else{
		return res.Msg,false
	}
}