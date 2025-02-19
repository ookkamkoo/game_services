package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "backend/app/migration"
	// "backend/app/models"
)

var (
	DB   *gorm.DB
	once sync.Once // ป้องกันการเรียกซ้ำ
)

func PG_Connect() error {
	once.Do(func() {
		// ✅ โหลดค่า ENV
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: No .env file found, using system environment")
		}

		// ✅ ดึงค่าจาก .env
		host := os.Getenv("PG_HOST")
		user := os.Getenv("PG_USER")
		password := os.Getenv("PG_PASSWORD")
		dbname := os.Getenv("PG_DBNAME")
		port := os.Getenv("PG_PORT")
		sslmode := "require"
		timezone := "Asia/Bangkok"

		// ✅ ตรวจสอบ ENV
		if host == "" || user == "" || password == "" || dbname == "" || port == "" {
			log.Fatal("❌ Database configuration is missing in .env file")
		}

		// ✅ สร้าง DSN
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			host, user, password, dbname, port, sslmode, timezone,
		)

		// ✅ เชื่อมต่อ Database
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("❌ Failed to connect to database: ", err)
		}

		// ✅ ตรวจสอบว่า Database พร้อมใช้งาน
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal("❌ Failed to get database connection: ", err)
		}

		// ✅ กำหนด Connection Pool
		// sqlDB.SetMaxOpenConns(100)    // จำนวน Connection สูงสุดที่เปิดพร้อมกัน
		// sqlDB.SetMaxIdleConns(5)      // จำนวน Connection ที่เปิดทิ้งไว้
		// sqlDB.SetConnMaxIdleTime(300) // เวลาสูงสุดที่ Connection อยู่ในสถานะ Idle
		// sqlDB.SetConnMaxLifetime(900) // เวลาสูงสุดของ Connection ก่อนถูกปิด

		sqlDB.SetMaxIdleConns(5)                  // ✅ จำกัดจำนวน Connection ที่ว่างไว้สูงสุด 10
		sqlDB.SetMaxOpenConns(100)                // ✅ จำกัดจำนวน Connection ที่เปิดพร้อมกันสูงสุด 100
		sqlDB.SetConnMaxIdleTime(2 * time.Minute) // time idle
		sqlDB.SetConnMaxLifetime(5 * time.Minute) // ✅ Connection แต่ละตัวสามารถใช้ได้นานสุด 1 ชั่วโมง

		log.Println("✅ Connected to Database")
		DB = db
	})
	return nil
}
