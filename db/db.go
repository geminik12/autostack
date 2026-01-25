/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/1/25 11:48:52
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/1/25 11:49:39
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
package db

import (
	"github.com/google/wire"
	redis "github.com/redis/go-redis/v9"
)

// ProviderSet is db providers.
var ProviderSet = wire.NewSet(
	NewMySQL,
	NewRedis,
	wire.Bind(new(redis.UniversalClient), new(*redis.Client)), // 正确绑定接口和实现
)

