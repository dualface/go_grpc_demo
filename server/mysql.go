package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"go_grpc_demo/config"
	"log"
)

type (
	MySQLConnect struct {
		DB   *sql.DB
		Conn *sql.Conn
		CTX  context.Context
	}
)

var mysqlConn *MySQLConnect

func InitMySQL(cfg *config.GlobalConfig, ctx context.Context) (*MySQLConnect, error) {
	if mysqlConn != nil {
		return nil, fmt.Errorf("MySQLConnect already initialized")
	}

	myc := mysql.NewConfig()
	myc.User = cfg.MySQL.Username
	myc.Passwd = cfg.MySQL.Password
	myc.Addr = cfg.MySQL.Addr
	myc.DBName = cfg.MySQL.DB
	myc.Collation = cfg.MySQL.Collation
	db, err := sql.Open("mysql", myc.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("[MySQL]: open failed, %s", err)
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("[MySQL]: connect failed, %s", err)
	}

	mysqlConn = &MySQLConnect{DB: db, Conn: conn, CTX: ctx}
	return mysqlConn, nil
}

func GetMySQL() *MySQLConnect {
	if mysqlConn == nil {
		log.Panicln("[MySQL]: not connected")
	}
	return mysqlConn
}
