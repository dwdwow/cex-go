package bnc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Req[T any] struct {
	Method  string
	Header  http.Header
	BaseURL string
	Path    string
	Params  T
	APIPvk  string
	retry   int
}

type Resp[T any] struct {
	StatusCode int    `json:"-"`
	Status     string `json:"-"`
	Code       int64  `json:"code"`
	Msg        string `json:"msg"`
	Data       T      `json:"-"`
}

func Request[ReqParams, RespData any](reqInfo Req[ReqParams]) (resp Resp[RespData], err error) {
	url := reqInfo.BaseURL + reqInfo.Path + "?" + structToSignedQuery(reqInfo.Params, reqInfo.APIPvk)
	req, err := http.NewRequest(reqInfo.Method, url, nil)
	if err != nil {
		return
	}
	req.Header = reqInfo.Header
	client := &http.Client{}
	rawResp, err := client.Do(req)
	if err != nil {
		return
	}
	defer rawResp.Body.Close()
	body, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return
	}
	statusCode := rawResp.StatusCode
	defer func() {
		resp.StatusCode = statusCode
		resp.Status = rawResp.Status
	}()
	if statusCode != http.StatusOK {
		if err = json.Unmarshal(body, &resp); err != nil {
			err = fmt.Errorf("bnc: request %d %s %w", statusCode, string(body), err)
			return
		}
		// sometimes, binance return -1021, which means timestamp for this request is outside of the recvWindow,
		// so retry with a new timestamp
		if resp.Code == -1021 && reqInfo.retry < 5 {
			reqInfo.retry++
			return Request[ReqParams, RespData](reqInfo)
		}
		err = fmt.Errorf("bnc: request %d %d %s", statusCode, resp.Code, resp.Msg)
		return
	} else {
		data := new(RespData)
		if err = json.Unmarshal(body, data); err != nil {
			err = fmt.Errorf("bnc: request %d %s %w", statusCode, string(body), err)
			return
		}
		resp.Data = *data
	}
	return
}
