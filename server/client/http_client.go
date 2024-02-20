// Package client http 客户端
package client

import (
	"crypto/tls"
	"fmt"
	 "log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/link1st/go-stress-testing/helper"
	"github.com/link1st/go-stress-testing/model"
	httplongclinet "github.com/link1st/go-stress-testing/server/client/http_longclinet"
	"golang.org/x/net/http2"
)

// logErr err
var logErr = log.New(os.Stderr, "", 0)

// HTTPRequest HTTP 请求
// method 方法 GET POST
// url 请求的url
// body 请求的body
// headers 请求头信息
// timeout 请求超时时间
func HTTPRequest(chanID uint64, request *model.Request) (resp *http.Response, requestTime uint64, err error) {
	method := request.Method
	url := request.URL

	urlParam := "?appId=%s&userId=%d&roomId=%d&liveMode=2&doremeVersion=3.3.6&sslMode=1&brand=TCL&transactionId=169691953524700&deviceId=N_6b317d8a97308460ss&osName=android&deviceType=T602DL&osVersion=2.3.4"
	arr := [3]string{"Lizhi_Heiye_20210727", "Lizhi_PP_20191010", "Lizhi_Xiaoximi_20220407"}
	randomIndex := rand.Intn(len(arr))
	appId := arr[randomIndex]
	userId := rand.New(rand.NewSource(time.Now().UnixNano())).Uint32()
	roomId := rand.New(rand.NewSource(time.Now().UnixNano())).Uint32()
	urlParam = fmt.Sprintf(urlParam, appId, userId, roomId)
	url = url + urlParam
	
	//fmt.Println("url:", url)

	body := request.GetBody()
	timeout := request.Timeout
	headers := request.Headers

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	// 在req中设置Host，解决在header中设置Host不生效问题
	if _, ok := headers["Host"]; ok {
		req.Host = headers["Host"]
	}
	// 设置默认为utf-8编码
	if _, ok := headers["Content-Type"]; !ok {
		if headers == nil {
			headers = make(map[string]string)
		}
		headers["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	var client *http.Client
	if request.Keepalive {
		client = httplongclinet.NewClient(chanID, request)
		startTime := time.Now()
		resp, err = client.Do(req)
		requestTime = uint64(helper.DiffNano(startTime))
		if err != nil {
			logErr.Println("请求失败:", err)
			return
		}
		return
	} else {
		req.Close = true
		tr := &http.Transport{}
		if request.HTTP2 {
			// 使用真实证书 验证证书 模拟真实请求
			tr = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
			}
			if err = http2.ConfigureTransport(tr); err != nil {
				return
			}
		} else {
			// 跳过证书验证
			tr = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}

		client = &http.Client{
			Transport: tr,
			Timeout:   timeout,
		}
	}

	startTime := time.Now()
	resp, err = client.Do(req)
	requestTime = uint64(helper.DiffNano(startTime))

	// TODO fortest
	// cuiresp, _ := io.ReadAll(resp.Body)
	// fmt.Println("cui:", string(cuiresp))

	if err != nil {
		logErr.Println("请求失败:", err)

		return
	}
	return
}
