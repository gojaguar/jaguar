package database

import (
	"github.com/go-jaguar/jaguar/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupConnectionSQL sets up a Database connection to an SQL database using the Gorm
// library. Depending on the given config.Database's engine, it will connect to either
// a MySQL or a Postgres database.
func SetupConnectionSQL(cfg config.Database) (*gorm.DB, error) {
	dialect := dialector(cfg.Engine)
	return gorm.Open(dialect(cfg.ToDNS()))
}

func dialector(eng string) func(dsn string) gorm.Dialector {
	switch eng {
	case config.EngineMySQL:
		return mysql.Open
	case config.EnginePostgres:
		return postgres.Open
	default:
		return mysql.Open
	}
}
