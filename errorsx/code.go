/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/1/26 23:43:11
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/1/26 23:51:23
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
package errorsx

import (
	"net/http"
)

var (
	// OK 代表请求成功.
	OK = &ErrorX{Code: http.StatusOK, Message: ""}

	// ErrInternal 表示所有未知的服务器端错误.
	ErrInternal = &ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError", Message: "Internal server error."}

	// ErrNotFound 表示资源未找到.
	ErrNotFound = &ErrorX{Code: http.StatusNotFound, Reason: "NotFound", Message: "Resource not found."}

	// ErrBind 表示请求体绑定错误.
	ErrBind = &ErrorX{Code: http.StatusBadRequest, Reason: "BindError", Message: "Error occurred while binding the request body to the struct."}

	// ErrInvalidArgument 表示参数验证失败.
	ErrInvalidArgument = &ErrorX{Code: http.StatusBadRequest, Reason: "InvalidArgument", Message: "Argument verification failed."}

	// ErrUnauthenticated 表示认证失败.
	ErrUnauthenticated = &ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated", Message: "Unauthenticated."}

	// ErrSignToken 表示签发 JWT Token 时出错.
	ErrSignToken = &ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated.SignToken", Message: "Error occurred while signing the JSON web token."}

	// ErrTokenInvalid 表示 JWT Token 格式无效.
	ErrTokenInvalid = &ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated.TokenInvalid", Message: "Token was invalid."}

	// ErrPermissionDenied 表示请求没有权限.
	ErrPermissionDenied = &ErrorX{Code: http.StatusForbidden, Reason: "PermissionDenied", Message: "Permission denied. Access to the requested resource is forbidden."}

	// ErrOperationFailed 表示操作失败.
	ErrOperationFailed = &ErrorX{Code: http.StatusConflict, Reason: "OperationFailed", Message: "The requested operation has failed. Please try again later."}
)
