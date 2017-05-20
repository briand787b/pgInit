package pgInit

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	path   string
	ip     string
	port   string
	dbName string
}

// Default values, can be overridden
var (
	dbConfigPath = "configuration/DBCredentials.json"
	ipAddress    = "127.0.0.1"
	portNumber   = "5432"
)

// Creates a new DBConfig set to default values
func NewDBConfig(databaseName string) *DBConfig {
	return &DBConfig{
		path:   dbConfigPath,
		ip:     ipAddress,
		port:   portNumber,
		dbName: databaseName,
	}
}

// This is the most efficient way to connect to the db
// with default credentials
func ConnectDefault(databaseName string) (*sql.DB, error) {
	config := NewDBConfig(databaseName)
	return config.Connect()
}

func (dbc *DBConfig) SetDBConfigPath(path string) {
	dbc.path = path
}

func (dbc *DBConfig) SetIPAddr(addr string) error {
	if net.ParseIP(addr) == nil {
		return errors.New("malformed IP address")
	}

	dbc.ip = addr
	return nil
}

func (dbc *DBConfig) SetPort(port int) error {
	if port < 0 {
		return errors.New("negative port number")
	}

	if port > 65535 {
		return errors.New("port number exceeds maximum value")
	}

	dbc.port = string(port)
	return nil
}

func (dbc *DBConfig) SetDatabaseName(name string) {
	dbc.dbName = name
}

func (dbc *DBConfig) Connect() (*sql.DB, error) {
	file, err := ioutil.ReadFile(dbc.path)
	if err != nil {
		return nil, err
	}

	var dbCredentials struct {
		Username string
		Password string
	}

	err = json.Unmarshal(file, &dbCredentials)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		dbCredentials.Username,
		dbCredentials.Password,
		dbc.ip,
		dbc.port,
		dbc.dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}