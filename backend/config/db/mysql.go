// mysql.go 初始化 MySQL 数据库连接，提供全局 *sql.DB 实例
package db

import (
	"ElainaWeb/config"
	"database/sql"
)

var DBPool *sql.DB // DBPool 全局数据库连接池

func InitDB(dbConfig *config.DbConfig) error {
	dsn := dbConfig.GetDSN()
	db, err := sql.Open(dbConfig.SqlName, dsn)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)

	DBPool = db

	return nil
}
