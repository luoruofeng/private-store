package main

import (
	"fmt"
	"github.com/luoruofeng/private-store/handler"
	"net/http"
)

func main() {

	//处理静态资源映射
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload_success", handler.HTTPInterceptor(handler.UploadSuccuess))
	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMeta))  //访问路径：http://localhost:8080/file/meta?filehash=sha1vlaue linux计算文件的sha1：sha1sum ./index.js
	http.HandleFunc("/file/download", handler.HTTPInterceptor(handler.Download)) //访问路径：http://localhost:8080/file/meta?download=sha1vlaue
	http.HandleFunc("/file/rename", handler.HTTPInterceptor(handler.Rename))
	http.HandleFunc("/file/delete", handler.HTTPInterceptor(handler.Delete))

	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}

}
