package db

import "github.com/ahmedkhaeld/banking-app/db/models"

func RunMigrations() error {
	if err := AddUUIDExtension(); err != nil {
		return err
	}

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Entry{},
		&models.Transfer{},
	); err != nil {
		return err
	}

	return nil
}

func AddUUIDExtension() error {
	if err := DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return err
	}
	return nil
}
