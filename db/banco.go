package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Driver de conexão com o mysql
)

// Conectar abre conexão com o banco de dados
func Conectar() (*sql.DB, error) {
	stringConexao := "root:admin@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", stringConexao)
	if err != nil {
		return nil, err
	}

	return db, nil
}
