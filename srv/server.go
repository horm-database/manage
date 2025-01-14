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

// Package srv is the Go implementation of server, which is designed to be high-performance,
// everything-pluggable and easy for testing.
package srv

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/horm-database/common/log/logger"
	"github.com/horm-database/common/types"
	"github.com/horm-database/manage/srv/codec"
	"github.com/horm-database/manage/srv/naming"
	"github.com/horm-database/manage/srv/transport"
	"github.com/horm-database/manage/srv/transport/web"
	"go.uber.org/automaxprocs/maxprocs"
)

// Description 服务描述信息
type Description struct {
	Name  string
	Funcs []Func
}

// Server is a server.
// One process, one server. A server may offer one or more services.
type Server struct {
	services map[string]Service

	failedServices sync.Map
	signalCh       chan os.Signal
	closeOnce      sync.Once
}

// NewServer 新建服务
func NewServer(serverDesc *Description) *Server {
	cfg, err := loadConfig(confFile)
	if err != nil {
		panic("load config fail: " + err.Error())
	}

	logger.CreateDefaultLogger(cfg.Log)

	// go maxprocs for docker
	maxprocs.Set(maxprocs.Logger(logger.DefaultLogger.Debugf))

	s := &Server{}

	codec.InitGlobalContext(cfg.Env, cfg.Machine, cfg.Server.Name)

	if cfg.Server.WebPort > 0 {
		webServiceName := "web." + cfg.Server.Name
		s.addService(webServiceName, newService(webServiceName, "web", cfg))
	}

	// 注册 service
	for _, srv := range s.services {
		err = srv.Register(serverDesc.Funcs)
		if err != nil {
			panic(fmt.Sprintf("register service error:%v", err))
		}
	}

	return s
}

func (s *Server) Close() {
	s.closeOnce.Do(func() {
		SetClosing()

		var wg sync.WaitGroup

		for name, srv := range s.services {
			if _, ok := s.failedServices.Load(name); ok {
				continue
			}

			wg.Add(1)
			go func(service Service) {
				defer wg.Done()

				c := make(chan struct{}, 1)
				go service.Close(c)

				select {
				case <-c:
				}
			}(srv)
		}

		// wait all service close
		wg.Wait()
	})
}

// addService adds a service to server.
func (s *Server) addService(serviceName string, service Service) {
	if s.services == nil {
		s.services = make(map[string]Service)
	}
	s.services[serviceName] = service
}

func newService(name, protocol string, cfg *config) Service {
	//配置参数
	opts := &Options{
		Protocol:         protocol,
		ServiceName:      name,
		Env:              cfg.Env,
		Machine:          cfg.Machine,
		Timeout:          types.GetMillisecond(cfg.Server.Timeout),
		MaxCloseWaitTime: types.GetMillisecond(cfg.Server.MaxCloseWaitTime),
		TransportOptions: transport.Options{
			ServiceName:  name,
			Protocol:     protocol,
			Network:      "tcp",
			EventLoopNum: cfg.Server.EventLoopNum,
			IdleTimeout:  types.GetMillisecond(cfg.Server.IdleTime),
		},
	}

	opts.CloseWaitTime = types.GetMillisecond(cfg.Server.CloseWaitTime)
	if opts.CloseWaitTime > maxCloseWaitTime { // 注销名字服务之后等待时间最多为 10s
		opts.CloseWaitTime = maxCloseWaitTime
	}

	if protocol == "web" {
		opts.Codec = web.DefaultServerCodec
		opts.Transport = web.DefaultWebTransport
		opts.Address = net.JoinHostPort(cfg.LocalIP, strconv.Itoa(int(cfg.Server.WebPort)))
		opts.TransportOptions.Address = opts.Address

		// TLS 证书配置
		opts.TransportOptions.TLSCertFile = cfg.Server.TLSCert
		opts.TransportOptions.TLSKeyFile = cfg.Server.TLSKey
		opts.TransportOptions.CACertFile = cfg.Server.CACert
	}

	if cfg.Register != nil && cfg.Register.Enable != false {
		// 名字服务注册器
		reg, err := naming.Add(protocol, opts.ServiceName, opts.Address, cfg.Register)
		if err != nil {
			panic("setup polaris config fail: " + err.Error())
		}
		opts.Registry = reg
	}

	return New(opts)
}

var (
	closing     bool
	closeLocker = new(sync.RWMutex)
)

func SetClosing() {
	closeLocker.Lock()
	closing = true
	closeLocker.Unlock()
}

func IsClosing() bool {
	closeLocker.RLock()
	defer closeLocker.RUnlock()
	return closing
}
