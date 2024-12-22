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

// Package web provides support for http protocol by default,
// provides rpc server with http protocol, and provides rpc database
// for calling http protocol.
package web

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/horm-database/common/codec"
	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/log"
	"github.com/horm-database/common/snowflake"
	cc "github.com/horm-database/manage/srv/codec"
	"github.com/horm-database/manage/srv/transport"

	"github.com/kavu/go_reuseport"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// DefaultWebTransport default server http client.
var DefaultWebTransport = NewWebTransport()

// transportWeb web client layer.
type transportWeb struct {
	Server *http.Server
}

// NewWebTransport create a new web transport.
func NewWebTransport() transport.Transport {
	return &transportWeb{}
}

// Serve starts listening and serve
func (t *transportWeb) Serve(ctx context.Context, opts *transport.Options) error {
	if opts.Handler == nil {
		return errors.New("http server client handler empty")
	}

	// All server http server client only register this http.Handler.
	serveFunc := func(w http.ResponseWriter, r *http.Request) {
		// Generates new empty general message structure body and save it to ctx.
		webCtx, msg := codec.NewMessage(r.Context())
		defer codec.RecycleMessage(msg)

		fc := frameCodec{Request: r, Response: w}
		msg.WithFrameCodec(&fc)
		msg.WithSpanID(snowflake.GenerateID())

		r = r.WithContext(webCtx)

		// Records LocalAddr and RemoteAddr to Context.
		localAddr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
		if ok {
			msg.WithLocalAddr(localAddr)
		}

		remoteAddr, _ := net.ResolveTCPAddr("tcp", r.RemoteAddr)
		msg.WithRemoteAddr(remoteAddr)

		_, err := opts.Handler.Handle(webCtx, []byte{})
		if err != nil {
			log.Error(cc.GCtx, errs.ErrSystem, "web server handle error: ", err)
			return
		}
	}

	s, err := t.newWebServer(serveFunc, opts)
	if err != nil {
		return err
	}

	t.configureWebServer(s, opts)

	if err := t.serveWeb(ctx, s, opts); err != nil {
		return err
	}

	return nil
}

func (t *transportWeb) serveWeb(ctx context.Context, s *http.Server, opts *transport.Options) (err error) {
	ln := opts.Listener
	if ln == nil {
		ln, err = reuseport.Listen(opts.Network, s.Addr)
		if err != nil {
			return fmt.Errorf("http reuseport listen error:%v", err)
		}
	}

	if len(opts.TLSKeyFile) != 0 && len(opts.TLSCertFile) != 0 {
		go func() {
			err := s.ServeTLS(
				tcpKeepAliveListener{ln.(*net.TCPListener)},
				opts.TLSCertFile,
				opts.TLSKeyFile,
			)

			if err != http.ErrServerClosed {
				log.Error(cc.GCtx, errs.ErrSystem, "serve TLS failed: ", err)
			}
		}()
	} else {
		go func() {
			_ = s.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
		}()
	}

	go func() {
		<-ctx.Done()
		_ = s.Shutdown(context.TODO())
	}()

	return nil
}

// configureWebServer sets properties of web server.
func (t *transportWeb) configureWebServer(svr *http.Server, opts *transport.Options) {
	if t.Server != nil {
		svr.ReadTimeout = t.Server.ReadTimeout
		svr.ReadHeaderTimeout = t.Server.ReadHeaderTimeout
		svr.WriteTimeout = t.Server.WriteTimeout
		svr.MaxHeaderBytes = t.Server.MaxHeaderBytes
		svr.IdleTimeout = t.Server.IdleTimeout
		svr.ConnState = t.Server.ConnState
		svr.ErrorLog = t.Server.ErrorLog
	}

	svr.IdleTimeout = opts.IdleTimeout
}

// newWebServer creates web server.
func (t *transportWeb) newWebServer(serveFunc func(w http.ResponseWriter, r *http.Request),
	opts *transport.Options) (*http.Server, error) {
	s := &http.Server{
		Addr:    opts.Address,
		Handler: http.HandlerFunc(serveFunc),
	}
	// Enable h2c without tls.
	if opts.EnableH2C {
		h2s := &http2.Server{}
		s.Handler = h2c.NewHandler(http.HandlerFunc(serveFunc), h2s)
		return s, nil
	}
	if len(opts.CACertFile) != 0 { // Enable two-way authentication to verify database certificate.
		s.TLSConfig = &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
		certPool, err := getCertPool(opts.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("http server get ca cert file error:%v", err)
		}
		s.TLSConfig.ClientCAs = certPool
	}
	return s, nil
}

// getCertPool gets certificate information.
func getCertPool(caCertFile string) (*x509.CertPool, error) {
	// "root" means to use the root ca certificate installed on the machine to verify database,
	// if there is not "root", means to use the input ca certificate to verify database.
	if caCertFile != "root" {
		ca, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return nil, err
		}
		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM(ca)
		if !ok {
			return nil, errors.New("appendCertsFromPEM fail")
		}
		return pool, nil
	}
	return nil, nil
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts new request.
func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	_ = tc.SetKeepAlive(true)
	_ = tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
