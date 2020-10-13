/*
 * @Author: panlq01@mingyuanyun.com
 * @Date: 2020-10-13 09:52:24
 * @Description: Some desc
 * @LastEditors: panlq01@mingyuanyun.com
 * @LastEditTime: 2020-10-13 15:54:17
 */
package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	// var auth Auth
	// db.Select("id").Where(Auth{
	// 	Username: username,
	// 	Password: password,
	// }).First(&auth)

	// if auth.ID > 0 {
	// 	return true
	// }

	// return false
	return true
}
