/*
 * @Author: panlq01@mingyuanyun.com
 * @Date: 2020-10-12 16:29:01
 * @Description: Some desc
 * @LastEditors: panlq01@mingyuanyun.com
 * @LastEditTime: 2020-10-13 18:42:32
 */

package util

import (
	"go-gin-ex/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()

	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}

	return result
}
