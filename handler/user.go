package handler

import (
	"encoding/json"
	store_const "github.com/luoruofeng/private-store/const"
	dblayer "github.com/luoruofeng/private-store/db"
	"github.com/luoruofeng/private-store/util"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ResponseBody struct {
	Code int `json:"code"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		html, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(html)
		return
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if len(username) < 3 || len(password) < 5 {
			w.Write([]byte("Invalid parameter"))
			return
		}

		enc_password := util.Sha1([]byte(password + store_const.PWD_SALT))
		if dblayer.UserSignup(username, enc_password) {
			resultBody := &ResponseBody{Code: 10000}
			if json, err := json.Marshal(resultBody); json != nil && err == nil {
				w.Write(json)
			}

		} else {
			w.Write([]byte("Failed"))
		}
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		file, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(file)
		return
	}

	//1.check info
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encpassword := util.Sha1([]byte(password + store_const.PWD_SALT))
	pwdchecked := dblayer.UserSignin(username, encpassword)
	if !pwdchecked {
		w.Write([]byte("Failed"))
		return
	}

	//2.generate token
	token := genToken(username)
	updateToken := dblayer.UpdateToken(username, token)
	if !updateToken {
		w.Write([]byte("Failed"))
		return
	}
	//3.redirect
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())

}

func genToken(username string) string {
	t := time.Now().Unix()
	ts := string(strconv.FormatInt(t, 16))
	return string(util.MD5([]byte(username+ts+store_const.TOKEN_SALT))) + ts[:8]
}

// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}

// UserInfoHandler ： 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//	token := r.Form.Get("token")

	// // 2. 验证token是否有效
	// isValidToken := IsTokenValid(token)
	// if !isValidToken {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	// 3. 查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}
