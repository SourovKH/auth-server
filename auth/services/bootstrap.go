package auth_services

import (
	"context"
	"log"
	"os"
	"time"

	auth_constants "lem-be/auth/constants"
	auth_models "lem-be/auth/models"
	auth_utils "lem-be/auth/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitSuperuser checks for an existing super_admin and creates one if it doesn't exist
func InitSuperuser(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := db.Collection("users")

	// Check if any super_admin exists
	var existingSuperAdmin auth_models.User
	err := usersCollection.FindOne(ctx, bson.M{"role": auth_constants.RoleSuperAdmin}).Decode(&existingSuperAdmin)
	if err == nil {
		log.Println("Superuser already exists.")
		return nil
	}

	if err != mongo.ErrNoDocuments {
		return err
	}

	// No super_admin found, create one from environment variables
	email := os.Getenv("SUPERUSER_EMAIL")
	password := os.Getenv("SUPERUSER_PASSWORD")

	if email == "" || password == "" {
		log.Println("WARNING: SUPERUSER_EMAIL or SUPERUSER_PASSWORD not set. Skipping superuser bootstrap.")
		return nil
	}

	hashedPassword, err := auth_utils.HashPassword(password)
	if err != nil {
		return err
	}

	superuser := auth_models.User{
		Email:     email,
		Password:  hashedPassword,
		Role:      auth_constants.RoleSuperAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = usersCollection.InsertOne(ctx, superuser)
	if err != nil {
		return err
	}

	log.Printf("Successfully bootstrapped superuser: %s\n", email)
	return nil
}
