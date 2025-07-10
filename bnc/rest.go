package bnc

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Req[T any] struct {
	Method             string
	Header             http.Header
	BaseURL            string
	Path               string
	Params             T
	APIPvk             string
	RespInMicroseconds bool
	retry              int
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
		if err == io.EOF && reqInfo.retry < 5 {
			reqInfo.retry++
			slog.Error("bnc: request EOF, retry", "retry", reqInfo.retry, "url", url, "error", err)
			time.Sleep(time.Second)
			return Request[ReqParams, RespData](reqInfo)
		}
		return
	}
	if reqInfo.RespInMicroseconds {
		req.Header.Set("X-MBX-TIME-UNIT", "MICROSECOND")
	}
	req.Header = reqInfo.Header
	rawResp, err := http.DefaultClient.Do(req)
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
			slog.Error("bnc: request -1021, retry", "retry", reqInfo.retry, "url", url, "error", err)
			time.Sleep(time.Second)
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
