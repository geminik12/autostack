/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/1/25 00:03:04
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/1/26 01:04:07
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
// Package server provides a simple server.
package server

import (
	"context"
	"time"

	"github.com/geminik12/autostack/log"
)

// Server 定义所有服务器类型的接口.
type Server interface {
	// RunOrDie 运行服务器，如果运行失败会退出程序（OrDie的含义所在）.
	RunOrDie()
	// GracefulStop 方法用来优雅关停服务器。关停服务器时需要处理 context 的超时时间.
	GracefulStop(ctx context.Context)
}

// Serve starts the server and blocks until the context is canceled.
// It ensures the server is gracefully shut down when the context is done.
func Serve(ctx context.Context, srv Server) error {
	go srv.RunOrDie()

	// Block until the context is canceled or terminated.
	<-ctx.Done()

	// Shutdown the server gracefully.
	log.Infof("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully stop the server.
	srv.GracefulStop(ctx)

	log.Infof("Server exited successfully.")

	return nil
}
