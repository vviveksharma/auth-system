// This service creates the roles in the table
package initsetup

import (
	"log"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

func InitRoles() {
	AdminId := "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	UserId := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	GuestId := "550e8400-e29b-41d4-a716-446655440000"
	ModeratorId := "1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed"

	db := db.DB
	err := CreateAdminEntry(db, AdminId)
	if err != nil {
		log.Fatal("error while creating the admin role entry", err)
	}
	err = CreateUserEntry(db, UserId)
	if err != nil {
		log.Fatal("error while creating the admin role entry", err)
	}
	err = CreateGuestEntry(db, GuestId)
	if err != nil {
		log.Fatal("error while creating the admin role entry", err)
	}
	err = CreateModEntry(db, ModeratorId)
	if err != nil {
		log.Fatal("error while creating the admin role entry", err)
	}
	log.Println("roles created successfully!!")
}


func CreateAdminEntry(db *gorm.DB, AdminId string) error {
	transaction := db.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	newEntry := transaction.Create(&models.DBRoles{
		Role:   "Admin",
		RoleId: uuid.MustParse(AdminId),
	})
	if newEntry.Error != nil {
		return newEntry.Error
	}
	transaction.Commit()
	return nil
}

func CreateUserEntry(db *gorm.DB, UserId string) error {
	transaction := db.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	newEntry := transaction.Create(&models.DBRoles{
		Role:   "User",
		RoleId: uuid.MustParse(UserId),
	})
	if newEntry.Error != nil {
		return newEntry.Error
	}
	transaction.Commit()
	return nil
}

func CreateGuestEntry(db *gorm.DB, GuesId string) error {
	transaction := db.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	newEntry := transaction.Create(&models.DBRoles{
		Role:   "Guest",
		RoleId: uuid.MustParse(GuesId),
	})
	if newEntry.Error != nil {
		return newEntry.Error
	}
	transaction.Commit()
	return nil
}

func CreateModEntry(db *gorm.DB, MoId string) error {
	transaction := db.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	newEntry := transaction.Create(&models.DBRoles{
		Role:   "Moderator",
		RoleId: uuid.MustParse(MoId),
	})
	if newEntry.Error != nil {
		return newEntry.Error
	}
	transaction.Commit()
	return nil
}
