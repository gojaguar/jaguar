package config

import "fmt"

// Config is used to configure an application metadata.
type Config struct {
	Environment string `env:"ENVIRONMENT" envDefault:"staging"`
	Name        string `env:"APPLICATION_NAME,required"`
	Port        int    `env:"APPLICATION_PORT" envDefault:"3030"`
}

// Database contains the information needed to establish connection with a database. It usually describes a
// config file structure (JSON/YAML) or the environment variables that should be read.
//
//	Developers use this data type to configure specific database services:
//
//	type Application struct {
//		UserDB Database `json:"user_db" envPrefix:"USER_DB_"`
//		DataDB Database `json:"data_db" envPrefix:"DATA_DB_"`
//	}
//
// Read more information about environment variables loader from caarlos0's library: https://github.com/caarlos0/env
type Database struct {
	Engine   string `json:"engine" env:"ENGINE" envDefault:"mysql"`
	Host     string `json:"host" env:"HOST,notEmpty"`
	User     string `json:"user" env:"USER,notEmpty"`
	Password string `json:"password" env:"PASSWORD,notEmpty"`
	Port     uint   `json:"port" env:"PORT" envDefault:"3306"`
	Name     string `json:"name" env:"NAME,notEmpty"`
	Charset  string `json:"charset" env:"CHARSET" envDefault:"utf8mb4"`
}

// ToDNS converts the current database config to a Data Source Name string, usually used to connect to a database.
func (d Database) ToDNS() string {
	switch d.Engine {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", d.User, d.Password, d.Host, d.Port, d.Name, d.Charset)
	default:
		return ""
	}
}

const (
	// EngineMySQL contains a string used to identify a MySQL database engine.
	EngineMySQL = "mysql"

	// EnginePostgres contains a string used to identify a Postgres database engine.
	EnginePostgres = "postgres"
)
