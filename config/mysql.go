package config

const (
	MYSQL_HOST     = "127.0.0.1"
	MYSQL_PORT     = "3306"
	MYSQL_USER     = "root"
	MYSQL_PASSWORD = "root"
	MYSQL_DB       = "ranking"
	MYSQL_CHARSET  = "utf8mb4"

	MYSQLDB = MYSQL_USER + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_HOST + ":" + MYSQL_PORT + ")/" + MYSQL_DB + "?charset=" + MYSQL_CHARSET + "&parseTime=True&loc=Local"

	//MYSQLDB = "root:root@tcp(127.0.0.1:3306)/ranking?charset=utf8mb4&parseTime=True&loc=Local"
)
