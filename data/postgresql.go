package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// PostgreSQLClient implementation of the Source interface for PostgreSQL
type PostgreSQLClient struct {
	db           *bun.DB
	ServerConfig map[string]interface{}
}

// NewPostgreSQLClient Generates a new PostgreSQL client from the information in the configuration file
func NewPostgreSQLClient(config map[string]interface{}) (*PostgreSQLClient, error) {

	pgParams, isOfCorrectType := config["postgresql_parameters"].(map[interface{}]interface{})
	if !isOfCorrectType {
		return nil, errors.New("`postgresql_parameters` item in config file of the wrong type, see README")
	}

	var parameters string

	ssl := fmt.Sprintf("%v", pgParams["ssl"])
	if ssl == "enabled" {
		parameters += "?sslmode=disable"
	}

	dsn := fmt.Sprintf("postgres://%v:@%v:%v/%v%s", pgParams["user"], pgParams["url"], pgParams["database_port"], pgParams["database_name"], parameters)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	return &PostgreSQLClient{
		db: db,
	}, nil
}

//GetDataSource return server data source
func (pgc *PostgreSQLClient) GetDataSource() string {
	return "PostgreSQL"
}

//GetConfig return server configuration
func (pgc *PostgreSQLClient) GetConfig() (map[string]interface{}, error) {
	return pgc.ServerConfig, nil
}

// GetData Implementation of the Source interface for Celestrak
func (pgc *PostgreSQLClient) GetData() ([]Satellite, error) {
	return []Satellite{}, nil
}
