package orm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
)

var db  *gorm.DB

func TestNewMySQL(t *testing.T) {
	db = NewMySQL(&Config{
		DSN: "user:pass@tcp(localhost:3306)/dbname?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=true&loc=Local&charset=utf8",
		Active: 5,
		Idle: 2,
		IdleTimeout: 3600,
	})
	db.DB().Ping()

	rows, err := db.Raw("select * from users limit 10").Rows()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(fmt.Sprintf("%+v", rows))
}