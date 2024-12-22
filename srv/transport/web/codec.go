// Copyright (c) 2024 The horm-database Authors. All rights reserved.
// This file Author:  CaoHao <18500482693@163.com> .
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	cc "github.com/horm-database/common/codec"
	"github.com/horm-database/common/consts"
	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/json"
	"github.com/horm-database/common/types"
	"github.com/horm-database/common/util"
	"github.com/horm-database/manage/srv/codec"
	"github.com/horm-database/manage/srv/transport/web/head"
)

var contentTypeSerializationType = map[string]int{
	"application/json": codec.SerializationTypeJSON,
	"application/xml":  codec.SerializationTypeXML,
}

var serializationTypeContentType = map[int]string{
	codec.SerializationTypeJSON: "application/json",
	codec.SerializationTypeXML:  "application/xml",
}

var (
	// DefaultServerCodec is the default http server codec.
	DefaultServerCodec = &ServerCodec{}
)

// ServerCodec http server side decoder.
type ServerCodec struct{}

// frameCodec frame codec
type frameCodec struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// Decode implements codec.Codec, decode http head
// http server has filled all the body of request into ctx, and reqBuf here is empty.
func (sc *ServerCodec) Decode(msg *cc.Msg, _ []byte) ([]byte, error) {
	fc := msg.FrameCodec().(*frameCodec)
	if fc == nil {
		return nil, errors.New("server decode missing frame codec in context")
	}

	reqBody, err := sc.getReqBody(fc, msg)
	if err != nil {
		return nil, err
	}

	if err := sc.setReqHeader(fc, msg); err != nil {
		return nil, err
	}

	return reqBody, nil
}

// Encode implements codec.Codec, encode http head.
// The buffer of the returned packet has been written to the
// response writer in head, no need to return respBuf.
func (sc *ServerCodec) Encode(msg *cc.Msg, respBody []byte) (b []byte, err error) {
	fc := msg.FrameCodec().(*frameCodec)
	if fc == nil {
		return nil, errors.New("server encode missing frame codec in context")
	}

	fc.Response.Header().Add("X-Content-Type-Options", "nosniff")
	ct := fc.Response.Header().Get("Content-Type")
	if ct == "" {
		ct = fc.Request.Header.Get("Content-Type")
		if fc.Request.Method == http.MethodGet || ct == "" {
			ct = "application/json"
		}
		fc.Response.Header().Add("Content-Type", ct)
	}

	if strings.Contains(ct, serializationTypeContentType[codec.SerializationTypeXML]) {
		fc.Response.Header().Set("Content-Type", "application/xml")
	}

	if e := msg.ServerRespError(); e != nil {
		err = errHandler(fc.Response, e)
		return
	}

	err = respHandler(fc.Response, respBody)
	return
}

func (sc *ServerCodec) setReqHeader(fc *frameCodec, msg *cc.Msg) error {
	reqHeader := &head.WebReqHeader{}
	msg.WithServerReqHead(reqHeader)

	reqHeader.RequestType = consts.RequestTypeWeb
	reqHeader.Callee = msg.CallRPCName()

	if v := fc.Request.Header.Get(head.Version); v != "" {
		reqHeader.Version = v
	}
	if v := fc.Request.Header.Get(head.RequestID); v != "" {
		i, _ := strconv.ParseUint(v, 10, 64)
		reqHeader.RequestId = i
		msg.WithRequestID(i)
	}
	if v := fc.Request.Header.Get(head.Timestamp); v != "" {
		reqHeader.Timestamp, _ = strconv.ParseUint(v, 10, 64)
	}
	if v := fc.Request.Header.Get(head.Timeout); v != "" {
		i, _ := strconv.Atoi(v)
		reqHeader.Timeout = uint32(i)
		msg.WithRequestTimeout(time.Millisecond * time.Duration(i))
	}
	if v := fc.Request.Header.Get(head.UserID); v != "" {
		reqHeader.Userid, _ = strconv.ParseUint(v, 10, 64)
	}
	if v := fc.Request.Header.Get(head.WorkspaceID); v != "" {
		i, _ := strconv.Atoi(v)
		reqHeader.WorkspaceId = uint32(i)
	}
	if v := fc.Request.Header.Get(head.Caller); v != "" {
		reqHeader.Caller = v
		msg.WithCallerServiceName(v)
	}
	if v := fc.Request.Header.Get(head.AuthRand); v != "" {
		i, _ := strconv.Atoi(v)
		reqHeader.AuthRand = uint32(i)
	}
	if v := fc.Request.Header.Get(head.Sign); v != "" {
		reqHeader.Sign = v
	}

	reqHeader.Ip = util.GetIpFromAddr(msg.RemoteAddr())

	respHeader := head.WebRespHeader{
		Version:   reqHeader.Version,
		RequestId: reqHeader.RequestId,
	}

	msg.WithServerRespHead(&respHeader)
	return nil
}

func (sc *ServerCodec) getReqBody(fc *frameCodec, msg *cc.Msg) ([]byte, error) {
	urlPath := fc.Request.URL.Path
	if urlPath[0] == '/' {
		urlPath = urlPath[1:]
	}

	msg.WithCallRPCName(urlPath)

	var reqBody []byte
	if fc.Request.Method == http.MethodGet {
		reqBody = types.StringToBytes(fc.Request.URL.RawQuery)
	} else {
		var exist bool
		ct := fc.Request.Header.Get("Content-Type")
		for contentType, serializationType := range contentTypeSerializationType {
			if strings.Contains(ct, contentType) {
				msg.WithSerializationType(serializationType)
				exist = true
				break
			}
		}

		if exist {
			var err error
			reqBody, err = getBody(fc.Request)
			if err != nil {
				return nil, err
			}
		}
	}
	return reqBody, nil
}

func getBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("body readAll: %w", err)
	}
	return body, nil
}

func errHandler(w http.ResponseWriter, e *errs.Error) error {
	result := map[string]interface{}{}
	result["code"] = e.Code
	result["msg"] = e.Msg

	if _, err := w.Write(json.Marshal(result)); err != nil {
		return fmt.Errorf("web write response error: %s", err.Error())
	}

	return nil
}

func respHandler(w http.ResponseWriter, respBody []byte) error {
	var resp []byte

	if len(respBody) == 0 {
		resp = []byte(`{"code":0,"msg":"success"}`)
	} else {
		resp = []byte(`{"code":0,"msg":"success","data":`)
		resp = append(resp, respBody...)
		resp = append(resp, "}"...)
	}

	if _, err := w.Write(resp); err != nil {
		return fmt.Errorf("web write response error: %s", err.Error())
	}
	return nil
}
