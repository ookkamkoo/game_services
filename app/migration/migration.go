package migration

import (
	"game_services/app/database"
	"game_services/app/models"
)

func RunMigration() {
	// migrations
	DB := database.DB

	// Drop tables
	err := DB.Migrator().DropTable(
		&models.Reports{},
		&models.Pg100Transactions{},
		&models.GplayTransactions{},
	)

	if err != nil {
		// Handle error
		panic(err) // Example: panic if migration fails
	}

	// Auto-migrate tables
	err = DB.AutoMigrate(
		&models.Reports{},
		&models.GplayTransactions{},
	)
	if err != nil {
		// Handle error
		panic(err) // Example: panic if migration fails
	}
}
