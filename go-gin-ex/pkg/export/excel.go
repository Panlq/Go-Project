/*
 * @Author: panlq01@mingyuanyun.com
 * @Date: 2020-10-13 19:03:50
 * @Description: Some desc
 * @LastEditors: panlq01@mingyuanyun.com
 * @LastEditTime: 2020-10-13 19:06:26
 */
package export

import (
	"go-gin-ex/pkg/setting"
)

func GetExcelFullUrl(name string) string {
	return setting.AppSetting.StaticFilePrefixUrl + "/" + GetExcelPath() + name
}

func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}
