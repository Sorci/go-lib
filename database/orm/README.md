#ORM
A mysql orm is based on gorm.

##Base Usage:
```go
    orm := NewMySQL(&Config{
		DSN: "user:pass@tcp(localhost:3306)/dbname?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=true&loc=Local&charset=utf8",
		Active: 5,
		Idle: 2,
		IdleTimeout: 3600,
	})

    orm.DB().Ping()

    rows, err := orm.Raw("select * from users limit 10").Rows()
```