package main

import (
	"fmt"
	"github.com/luoruofeng/private-store/handler"
	"net/http"
)

func main(){
	http.HandleFunc("/file/upload",handler.UploadHandler)
	http.HandleFunc("/file/upload_success",handler.UploadSuccuess)
	http.HandleFunc("/file/meta",handler.GetFileMeta)//访问路径：http://localhost:8080/file/meta?filehash=sha1vlaue linux计算文件的sha1：sha1sum ./index.js
	http.HandleFunc("/file/download",handler.Download)//访问路径：http://localhost:8080/file/meta?download=sha1vlaue
	http.HandleFunc("/file/rename",handler.Rename)
	http.HandleFunc("/file/delete",handler.Delete)

	if err := http.ListenAndServe(":8080",nil);err !=nil{
		fmt.Printf("error:%s\n",err.Error())
	}

}
