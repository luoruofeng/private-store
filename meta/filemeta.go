package meta

//FileMeta : file entity
type FileMeta struct{
	FileSha1 string
	FileName string
	Location string
	UploadAt string
	FileSize int64
}

var fileMetas map[string]FileMeta

func init(){
	fileMetas = make(map[string]FileMeta)
}

//Update : add or update file to fileMetas
func Update(fileMeta FileMeta){
	fileMetas[fileMeta.FileSha1] = fileMeta
}

//Get : get file from fileMetas by Sha1
func Get(sha1 string)FileMeta{
	return fileMetas[sha1]
}

func Delete(sha1 string){
	delete(fileMetas,sha1)
}