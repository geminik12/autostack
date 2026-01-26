/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/1/26 23:56:03
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/1/26 23:58:36
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
package model

import "time"

const TableNameUserM = "user"

// UserM mapped from table <user>
type UserM struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID    string    `gorm:"column:userID;not null;uniqueIndex:idx_user_userID;comment:用户唯一 ID" json:"userID"`       // 用户唯一 ID
	Username  string    `gorm:"column:username;not null;uniqueIndex:idx_user_username;comment:用户名（唯一）" json:"username"` // 用户名（唯一）
	Password  string    `gorm:"column:password;not null;comment:用户密码（加密后）" json:"password"`                             // 用户密码（加密后）
	Nickname  string    `gorm:"column:nickname;not null;comment:用户昵称" json:"nickname"`                                  // 用户昵称
	Email     string    `gorm:"column:email;not null;comment:用户电子邮箱地址" json:"email"`                                    // 用户电子邮箱地址
	Phone     string    `gorm:"column:phone;not null;uniqueIndex:idx_user_phone;comment:用户手机号" json:"phone"`            // 用户手机号
	CreatedAt time.Time `gorm:"column:createdAt;not null;default:current_timestamp;comment:用户创建时间" json:"createdAt"`    // 用户创建时间
	UpdatedAt time.Time `gorm:"column:updatedAt;not null;default:current_timestamp;comment:用户最后修改时间" json:"updatedAt"`  // 用户最后修改时间
}

// TableName UserM's table name
func (*UserM) TableName() string {
	return TableNameUserM
}
