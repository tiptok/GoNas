package global

import (
	"database/sql"

	"log"

	_ "github.com/alexbrainman/odbc"
	//"github.com/denisenkom/go-mssqldb"
)

var DBConn *sql.DB

func DBInstance() *sql.DB {
	var err error
	if DBConn == nil {
		DBConn, err = sql.Open("odbc", Param.DBConnString)
		if err != nil {
			log.Println("Get DBInstance Err:%v", err)
			DBConn.Close()
		}
		if DBConn != nil {
			log.Println("Init DB %v", DBConn)
		}
	}
	return DBConn
}

// var DBConn *sql.DB

// func DBInstance() *sql.DB {
// 	var err error
// 	if DBConn == nil {
// 		DBConn, err = sql.Open("sqlserver", Param.DBConnString)
// 		if err != nil {
// 			log.Println("Get DBInstance Err:%v", err)
// 			DBConn.Close()
// 		}
// 		DBConn.SetConnMaxLifetime(time.Duration(10) * time.Minute)
// 		log.Println(DBConn.Stats())
// 		if DBConn != nil {
// 			log.Println("Init DB %v", DBConn)
// 		}
// 	}
// 	return DBConn
// }
