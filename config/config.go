package config

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=1234 dbname=shop port=5432 sslmode=disable TimeZone=Asia/Almaty"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("ошибка подключения к БД: %w", err))
	}
	DB = db

	migrationURL := "postgres://postgres:1234@localhost:5432/shop?sslmode=disable"

	runMigrations(migrationURL)
}

func runMigrations(url string) {
	m, err := migrate.New("file://migrations", url)
	if err != nil {
		fmt.Printf("Ошибка инициализации миграций: %v\n", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Printf("Ошибка применения миграций: %v\n", err)
	} else if err == migrate.ErrNoChange {
		fmt.Println("Миграции не требуются (база уже актуальна).")
	} else {
		fmt.Println("Миграции успешно применены!")
	}
}
