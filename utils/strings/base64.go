/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/2/1 23:21:17
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/2/1 23:21:37
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
package strings

import (
	"bytes"
	"encoding/base64"
	"io"
)

func DecodeBase64(i string) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(i)))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
