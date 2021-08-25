package utils

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)


type Database struct {
	pg    *pg.DB
	pgorm *gorm.DB
}

func (db *Database) GetPg() *pg.DB {
	return db.pg
}

func (db *Database) GetGorm() *gorm.DB {
	return db.pgorm
}

func NewDataBase() *Database {
	return &Database{
		pg:    newPgDB(),
		pgorm: newGormDB(),
	}
}

func newPgDB() *pg.DB {
	return connectPg(
		&pg.Options{
			User:         Conf.Databases.User,
			Password:     Conf.Databases.Password,
			Database:     Conf.Databases.Database,
			Addr:         fmt.Sprintf("%s:%d", Conf.Databases.Host, Conf.Databases.Port),
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 5,
			IdleTimeout:  time.Second * 120,
			PoolSize:     Conf.Databases.PoolSize,
		}, Conf.Databases.Slow)
}

func newGormDB() *gorm.DB {
	instance, err := gorm.Open("postgres", connStr())
	if err != nil {
		return nil
	}
	instance.DB().SetConnMaxLifetime(time.Minute * 5)
	instance.DB().SetMaxIdleConns(10)
	instance.DB().SetMaxOpenConns(Conf.Databases.PoolSize)
	instance.LogMode(true)
	return instance
}

func connectPg(opt *pg.Options, slow int) *pg.DB {
	db := pg.Connect(opt)
	var n string
	_, err := db.QueryOne(pg.Scan(&n), "select now() ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("connect pg %s %s success on %s", db.String(), opt.Database, n)

	return db
}

func connStr() string {
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%d",
		Conf.Databases.Host, Conf.Databases.User, Conf.Databases.Database, Conf.Databases.Password, Conf.Databases.Port)
	return connStr
}
