/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/1/27 00:39:34
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/1/27 00:39:51
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
package gin

import (
	"github.com/geminik12/autostack/contextx"
	"github.com/geminik12/autostack/known"
	"github.com/geminik12/autostack/log"
	"github.com/gin-gonic/gin"
)

// AuthnBypasswMiddleware 是一个认证中间件，
// 用于从 gin.Context 的 Header 中提取用户 ID，模拟所有请求认证通过。
func AuthnBypasswMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header中提取用户ID，假设请求头名称为 "X-User-ID"
		userID := "user-000001" // 默认用户ID
		if val := c.GetHeader(known.XUserID); val != "" {
			userID = val
		}

		log.Debugw("Simulated authentication successful", "userID", userID)

		// 将用户ID和用户名注入到上下文中
		ctx := contextx.WithUserID(c.Request.Context(), userID)
		c.Request = c.Request.WithContext(ctx)

		// 继续后续的操作
		c.Next()
	}
}
