package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

// Config mysql config.
type Config struct {
	DSN         string         // data source name.
	Active      int            // pool
	Idle        int            // pool
	IdleTimeout time.Duration  // connect max life time.
}

/**
 * @title NewMysql
 * @description new mysql client
 * @param c *Config
 * @return *gorm.DB
 **/
func NewMySQL(c *Config) *gorm.DB {
	db, err := gorm.Open("mysql", c.DSN)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(c.Idle)
	db.DB().SetMaxOpenConns(c.Active)
	db.DB().SetConnMaxLifetime(c.IdleTimeout / time.Second)
	return db
}