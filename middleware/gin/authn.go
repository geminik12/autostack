/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/1/26 23:38:19
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/1/26 23:59:49
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
package gin

import (
	"context"

	"github.com/geminik12/autostack/contextx"
	"github.com/geminik12/autostack/core"
	"github.com/geminik12/autostack/errorsx"
	"github.com/geminik12/autostack/log"
	"github.com/geminik12/autostack/model"
	"github.com/geminik12/autostack/token"
	"github.com/gin-gonic/gin"
)

// UserRetriever 用于根据用户名获取用户的接口.
type UserRetriever interface {
	// GetUser 根据用户ID获取用户信息
	GetUser(ctx context.Context, userID string) (*model.UserM, error)
}

// AuthnMiddleware 是一个认证中间件，用于从 gin.Context 中提取 token 并验证 token 是否合法.
func AuthnMiddleware(retriever UserRetriever) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析 JWT Token
		userID, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, nil, errorsx.ErrTokenInvalid.WithMessage("%s", err.Error()))
			c.Abort()
			return
		}

		log.Debugw("Token parsing successful", "userID", userID)

		user, err := retriever.GetUser(c, userID)
		if err != nil {
			core.WriteResponse(c, nil, errorsx.ErrUnauthenticated.WithMessage("%s", err.Error()))
			c.Abort()
			return
		}

		ctx := contextx.WithUserID(c.Request.Context(), user.UserID)
		ctx = contextx.WithUsername(ctx, user.Username)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
