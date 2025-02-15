package main

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/db"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Insert 3 users into database
func main() {
	db.Init()
	defer db.Close()

	usersRepo := repositories.UsersRepository{Db: db.Pool}
	users := []models.User{}

	// users = append(users, *models.CreateUser("abc123", "asdf"))
	// users = append(users, *models.CreateUser("xcvnjk1", "czxijklvn"))
	// users = append(users, *models.CreateUser("890sdf", "loxkcvb"))

	// err := usersRepo.InsertManyUsers(users)
	users = append(users, *models.CreateUser("abc123", "asdf", 1))
	users = append(users, *models.CreateUser("xcvnjk1", "czxijklvn", 2))
	users = append(users, *models.CreateUser("890sdf", "loxkcvb", 2))

	for _, user := range users {
		err := usersRepo.InsertOneUser(&user)
		if err != nil {
			utils.Logger.Error().Err(err)
		}
	}
}
