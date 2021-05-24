package handler

import (
	"encoding/json"
	"fmt"
	store_const "github.com/luoruofeng/private-store/const"
	"github.com/luoruofeng/private-store/meta"
	"github.com/luoruofeng/private-store/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(res http.ResponseWriter,req *http.Request)  {
	if req.Method == http.MethodGet{
		byteContent,err := ioutil.ReadFile("./static/test.html")
		if err != nil{
			io.WriteString(res,err.Error())
			return
		}
		io.WriteString(res,string(byteContent))
	}else if req.Method == http.MethodPost{
		fmt.Printf("content-type:%s\n",req.Header.Get("Content-Type"))
		file,fileHead,err := req.FormFile("file")
		if err != nil{
			fmt.Printf("failed to get data %s \n",err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: fileHead.Filename,
			Location: store_const.STORE_ROOT+fileHead.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),

		}

		newFile,err := os.Create(fileMeta.Location)
		if err != nil{
			fmt.Printf("failed to create file error:%s\n",err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize,err = io.Copy(newFile,file)
		if err != nil{
			fmt.Printf("failed to write file error:%s\n",err.Error())
			return
		}

		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.Update(fileMeta)

		http.Redirect(res,req,"/file/upload_success",http.StatusFound)
	}
}


func UploadSuccuess(res http.ResponseWriter,req *http.Request){
	res.Write([]byte("upload success"))
}


// GetFileMeta : get file meta
func GetFileMeta(res http.ResponseWriter,req *http.Request){
	req.ParseForm()
	filehash := req.Form["filehash"][0]
	fMeta := meta.Get(filehash)
	data,err := json.Marshal(fMeta)
	if err != nil{
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Write(data)
}

func Download(rep http.ResponseWriter,req *http.Request){
	req.ParseForm()
	sha1 := req.Form.Get("filehash")
	fmeta := meta.Get(sha1)
	file,err := os.Open(fmeta.Location)
	if err != nil{
		rep.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data,err := ioutil.ReadAll(file)
	if err !=nil{
		rep.WriteHeader(http.StatusInternalServerError)
		return
	}

	rep.Header().Set("Content-Type","application/octect-stream")
	rep.Header().Set("content-disposition","attachment;filename=\""+fmeta.FileName+"\"")
	rep.Write(data)
}

func Rename(rep http.ResponseWriter , req *http.Request){
	req.ParseForm()
	opType := req.Form.Get("optype")
	filehash := req.Form.Get("filehash")
	filename := req.Form.Get("filename")

	if opType == "0"{
		rep.WriteHeader(http.StatusForbidden)
		return
	}

	if req.Method != http.MethodPost{
		rep.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmeta := meta.Get(filehash)
	fmt.Printf("rename old:%s new:%s\n",fmeta.Location,store_const.STORE_ROOT+filename)
	os.Rename(fmeta.Location,store_const.STORE_ROOT+filename)
	fmeta.FileName = filename
	fmeta.Location = store_const.STORE_ROOT+filename
	meta.Update(fmeta)



	if data,err := json.Marshal(fmeta) ; err != nil{
		rep.WriteHeader(http.StatusInternalServerError)
		return
	}else{
		rep.WriteHeader(http.StatusOK)
		rep.Write(data)
	}
}

func Delete(rep http.ResponseWriter,req *http.Request){
	req.ParseForm()
	fileHash := req.Form.Get("filehash")

	fmeta := meta.Get(fileHash)
	meta.Delete(fmeta.FileSha1)
	os.Remove(fmeta.Location)

	if data,err := json.Marshal(fmeta);err != nil{
		rep.WriteHeader(http.StatusInternalServerError)
		return
	}else {
		rep.WriteHeader(http.StatusOK)
		rep.Write(data)
	}
}