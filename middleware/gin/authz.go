/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/1/26 23:38:11
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/1/26 23:53:28
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
package gin

import (
	"github.com/geminik12/autostack/contextx"
	"github.com/geminik12/autostack/core"
	"github.com/geminik12/autostack/errorsx"
	"github.com/geminik12/autostack/log"
	"github.com/gin-gonic/gin"
)

// Authorizer 用于定义授权接口的实现.
type Authorizer interface {
	Authorize(subject, object, action string) (bool, error)
}

// AuthzMiddleware 是一个 Gin 中间件，用于进行请求授权.
func AuthzMiddleware(authorizer Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := contextx.UserID(c.Request.Context())
		object := c.Request.URL.Path
		action := c.Request.Method

		// 记录授权上下文信息
		log.Debugw("Build authorize context", "subject", subject, "object", object, "action", action)

		// 调用授权接口进行验证
		if allowed, err := authorizer.Authorize(subject, object, action); err != nil || !allowed {
			core.WriteResponse(c, nil, errorsx.ErrPermissionDenied.WithMessage(
				"access denied: subject=%s, object=%s, action=%s, reason=%v",
				subject,
				object,
				action,
				err,
			))
			c.Abort()
			return
		}

		c.Next() // 继续处理请求
	}
}
