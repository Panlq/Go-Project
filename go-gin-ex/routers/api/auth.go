/*
 * @Author: panlq01@mingyuanyun.com
 * @Date: 2020-10-13 09:52:24
 * @Description: Some desc
 * @LastEditors: panlq01@mingyuanyun.com
 * @LastEditTime: 2020-10-13 15:43:41
 */
package api

import (
	"go-gin-ex/models"
	"go-gin-ex/pkg/e"
	"go-gin-ex/pkg/logging"
	"go-gin-ex/pkg/util"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				code = e.SUCCESS
			}

		} else {
			for _, err := range valid.Errors {
				logging.Info(err.Key, err.Message)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}
