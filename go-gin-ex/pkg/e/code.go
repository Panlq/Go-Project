/*
 * @Author: panlq01@mingyuanyun.com
 * @Date: 2020-10-12 16:26:11
 * @Description: Some desc
 * @LastEditors: panlq01@mingyuanyun.com
 * @LastEditTime: 2020-10-13 18:23:38
 */

package e

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_EXIST_TAG         = 10001
	ERROR_NOT_EXIST_TAG     = 10002
	ERROR_NOT_EXIST_ARTICLE = 10003

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004

	// 保存图片失败
	ERROR_UPLOAD_SAVE_IMAGE_FAIL = 30001
	// 检查图片失败
	ERROR_UPLOAD_CHECK_IMAGE_FAIL = 30002
	// 校验图片错误，图片格式或大小有问题
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT = 30003
)
