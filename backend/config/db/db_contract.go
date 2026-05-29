// db_contract.go 定义数据库连接池接口，用于依赖注入和测试 mock
package db

import (
	"context"
	"database/sql"
)

// DBTX 封装 *sql.DB 的常用操作，用于依赖注入和测试 mock。
// *sql.DB 天然满足此接口，无需额外适配。
type DBTX interface {
	// QueryRow 执行查询并返回单行结果
	QueryRow(query string, args ...any) *sql.Row
	// Query 执行查询并返回多行结果
	Query(query string, args ...any) (*sql.Rows, error)
	// Exec 执行写入/更新/删除语句
	Exec(query string, args ...any) (sql.Result, error)
	// Ping 检测数据库连接是否存活
	Ping() error
	// PingContext 带上下文的连接检测
	PingContext(ctx context.Context) error
	// Begin 开启事务
	Begin() (*sql.Tx, error)
	// BeginTx 带上下文和选项开启事务
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	// SetMaxIdleConns 设置最大空闲连接数
	SetMaxIdleConns(n int)
	// SetMaxOpenConns 设置最大打开连接数
	SetMaxOpenConns(n int)
}
