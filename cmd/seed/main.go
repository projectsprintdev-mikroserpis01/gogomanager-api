package main

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/infra/database"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/infra/env"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/bcrypt"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/log"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/uuid"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/validator"
)

const SeedersFilePath = "data/seeders/"
const SeedersDevPath = SeedersFilePath + "dev/"
const SeedersProdPath = SeedersFilePath + "prod/"

func main() {
	psqlDB := database.NewPgsqlConn()
	defer psqlDB.Close()

	var path string
	if env.AppEnv.AppEnv == "production" {
		path = SeedersProdPath
	} else {
		path = SeedersDevPath
	}

	validator := validator.Validator
	uuid := uuid.UUID
	bcrypt := bcrypt.Bcrypt

	seedUsers(path, psqlDB, validator, uuid, bcrypt)
}

func seedUsers(path string, db *sqlx.DB, validator validator.ValidatorInterface, uuid uuid.UUIDInterface, bcrypt bcrypt.BcryptInterface) {
	path += "users.csv"

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err,
		}, "[seed][seedUsers] Error opening file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err,
		}, "[seed][seedUsers] Error reading file")
	}

	for _, record := range records {
		log.Info(log.LogInfo{
			"record": record,
		}, "[seed][seedUsers] Inserting record")

		roleID, err := strconv.Atoi(record[3])
		if err != nil {
			log.Fatal(log.LogInfo{
				"error": err,
			}, "[seed][seedUsers] Error converting role id")
		}

		req := dto.CreateUserRequest{
			Name:     record[0],
			Email:    record[1],
			Password: record[2],
			RoleID:   roleID,
		}

		valErr := validator.Validate(req)
		if valErr != nil {
			log.Fatal(log.LogInfo{
				"error": valErr,
			}, "[seed][seedUsers] Error validating user")
		}

		id, err := uuid.NewV7()
		if err != nil {
			log.Fatal(log.LogInfo{
				"error": err,
			}, "[seed][seedUsers] Error generating uuid")
		}

		hashedPassword, err := bcrypt.Hash(req.Password)
		if err != nil {
			log.Fatal(log.LogInfo{
				"error": err,
			}, "[seed][seedUsers] Error hashing password")
		}

		user := &entity.User{
			ID:       id,
			Name:     req.Name,
			Email:    req.Email,
			Password: hashedPassword,
			RoleID:   req.RoleID,
		}

		_, err = db.NamedExec(
			`INSERT INTO users (id, name, email, password, role_id) VALUES (:id, :name, :email, :password, :role_id)`,
			user,
		)

		if err != nil {
			log.Fatal(log.LogInfo{
				"error": err,
			}, "[seed][seedUsers] Error inserting user")
		}
	}
}
