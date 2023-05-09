package database

import (
	"errors"
	"fmt"
	"github.com/gojaguar/jaguar/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrInvalidDialect = errors.New("an invalid or unsupported dialect was provided")
)

// SetupConnectionSQL sets up a Database connection to an SQL database using the Gorm
// library. Depending on the given config.Database's engine, it will connect to either
// a MySQL or a Postgres database.
func SetupConnectionSQL(cfg config.Database) (*gorm.DB, error) {
	dialect := dialector(cfg.Engine)
	if dialect == nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDialect, cfg.Engine)
	}
	return gorm.Open(dialect(cfg.ToDNS()))
}

func dialector(eng string) func(dsn string) gorm.Dialector {
	switch eng {
	case config.EngineMySQL:
		return mysql.Open
	case config.EnginePostgres:
		return postgres.Open
	case config.EngineSQLite:
		return sqlite.Open
	default:
		return noOp
	}
}

func noOp(_ string) gorm.Dialector {
	return nil
}
