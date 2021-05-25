package store_const

const (
	STORE_ROOT = "C:\\Users\\luoruofeng\\Desktop\\"
	MYSQL_NAME = "root"
	MYSQL_PW   = "root"
	DB_NAME    = "fileserver"
	DB_ADDR    = MYSQL_NAME + ":" + MYSQL_PW + "@tcp(127.0.0.1:3306)/" + DB_NAME + "?charset=utf8"
)
