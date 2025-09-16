package store

import (
	"fmt"
	"gorm.io/driver/postgres"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var PDB *gorm.DB // PDB for PostgreSQL DB

func InitDB() error {
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *gorm.DB {
	return DB
}

// --- 【新增】PostgreSQL 全局实例 ---

// InitPDB initializes the PostgreSQL database connection.
func InitPDB() error {
	host := viper.GetString("database_pg.host")
	port := viper.GetInt("database_pg.port")
	user := viper.GetString("database_pg.user")
	password := viper.GetString("database_pg.password")
	dbname := viper.GetString("database_pg.dbname")
	sslmode := viper.GetString("database_pg.sslmode")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode)

	var err error
	PDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

// GetPDB returns the global PostgreSQL database instance.
func GetPDB() *gorm.DB {
	return PDB
}
