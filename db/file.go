package db

import (
	"database/sql"
	"fmt"
	mydb "github.com/luoruofeng/private-store/db/mysql"
)

func OnFileUploadFinished(filehash string, filename string, fileaddr string, filesize int64) bool {
	if stmt, err := mydb.DBConn().Prepare("insert into tbl_file(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,1)"); err != nil {
		fmt.Printf("Failed to prepare statement,err:%s\n", err.Error())
		return false
	} else {
		defer stmt.Close()
		exec, err := stmt.Exec(filehash, filename, filesize, fileaddr)
		if err != nil {
			fmt.Printf("failed to exec.err:%s\n", err.Error())
			return false
		}
		if result, err := exec.RowsAffected(); err != nil {
			fmt.Printf("failed to get rowAffected.err:%s\n", err.Error())
			return false
		} else {
			if result <= 0 {
				fmt.Printf("File With hash:%s has been uploaded before \n", filehash)
				return false
			} else {
				fmt.Printf("Success upload.File With hash:%s\n", filehash)
				return true
			}
		}
		return false
	}
}

// TableFile : 文件表结构体
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMeta : 从mysql获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file " +
			"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(
		&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应记录， 返回参数及错误均为nil
			return nil, nil
		}
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, nil
}
