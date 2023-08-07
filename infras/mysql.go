package infras

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/evermos/boilerplate-go/shared/failure"

	"github.com/evermos/boilerplate-go/configs"
	// use MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	maxIdleConnection = 10
	maxOpenConnection = 10
)

// Block contains a transaction block
type Block func(db *sqlx.Tx, c chan error)

// MySQLConn wraps a pair of read/write MySQL connections.
type MySQLConn struct {
	Read  *sqlx.DB
	Write *sqlx.DB
}

// ProvideMySQLConn is the provider for MySQLConn.
func ProvideMySQLConn(config *configs.Config) *MySQLConn {
	return &MySQLConn{
		Read:  CreateMySQLReadConn(*config),
		Write: CreateMySQLWriteConn(*config),
	}
}

// CreateMySQLWriteConn creates a database connection for write access.
func CreateMySQLWriteConn(config configs.Config) *sqlx.DB {
	return CreateDBConnection(
		"write",
		config.DB.MySQL.Write.Username,
		config.DB.MySQL.Write.Password,
		config.DB.MySQL.Write.Host,
		config.DB.MySQL.Write.Port,
		config.DB.MySQL.Write.Name,
		config.DB.MySQL.Write.Timezone)

}

// CreateMySQLReadConn creates a database connection for read access.
func CreateMySQLReadConn(config configs.Config) *sqlx.DB {
	return CreateDBConnection(
		"read",
		config.DB.MySQL.Read.Username,
		config.DB.MySQL.Read.Password,
		config.DB.MySQL.Read.Host,
		config.DB.MySQL.Read.Port,
		config.DB.MySQL.Read.Name,
		config.DB.MySQL.Read.Timezone)

}

// CreateDBConnection creates a database connection.
func CreateDBConnection(name, username, password, host, port, dbName, timeZone string) *sqlx.DB {
	descriptor := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=%s&parseTime=true",
		username,
		password,
		host,
		port,
		dbName,
		url.QueryEscape(timeZone))
	db, err := sqlx.Connect("mysql", descriptor)
	if err != nil {
		log.
			Fatal().
			Err(err).
			Str("name", name).
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Failed connecting to database")
	} else {
		log.
			Info().
			Str("name", name).
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Connected to database")
	}
	db.SetMaxIdleConns(maxIdleConnection)
	db.SetMaxOpenConns(maxOpenConnection)

	return db
}

// OpenMock opens a database connection for mocking purposes.
func OpenMock(db *sql.DB) *MySQLConn {
	conn := sqlx.NewDb(db, "mysql")
	return &MySQLConn{
		Write: conn,
		Read:  conn,
	}
}

// WithTransaction performs queries with transaction
func (m *MySQLConn) WithTransaction(block Block) (err error) {
	e := make(chan error)
	tx, err := m.Write.Beginx()
	if err != nil {
		return
	}
	go block(tx, e)
	err = <-e
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			err = failure.InternalError(errTx)
		}
		return
	}
	err = tx.Commit()
	return
}
