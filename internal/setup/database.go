package setup

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" //
	"github.com/spf13/viper"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewDatabase(v *viper.Viper) (*sql.DB, error) {

	dcfg := &DbConfig{
		Host:     v.GetString("posgres.host"),
		Port:     v.GetString("posgres.port"),
		User:     v.GetString("posgres.user"),
		Password: v.GetString("posgres.password"),
		Database: v.GetString("posgres.database"),
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", dcfg.User, dcfg.Password, dcfg.Database, dcfg.Port)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
		return nil, err
	}
	return db, nil
}
