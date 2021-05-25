package meta

import "github.com/luoruofeng/private-store/db"

//FileMeta : file entity
type FileMeta struct {
	FileSha1 string
	FileName string
	Location string
	UploadAt string
	FileSize int64
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

//Update : add or update file to fileMetas
func Update(fileMeta FileMeta) {
	fileMetas[fileMeta.FileSha1] = fileMeta
}

func UploadDB(fileMeta FileMeta) bool {
	return db.OnFileUploadFinished(fileMeta.FileSha1, fileMeta.FileName, fileMeta.Location, fileMeta.FileSize)
}

//Get : get file from fileMetas by Sha1
func Get(sha1 string) FileMeta {
	return fileMetas[sha1]
}

func Delete(sha1 string) {
	delete(fileMetas, sha1)
}

//GetFileMetaDB:从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tfile, err := db.GetFileMeta(fileSha1)
	if tfile == nil || err != nil {
		return nil, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return &fmeta, nil
}
