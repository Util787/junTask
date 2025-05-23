package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/Util787/junTask/entities"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetAllUsers(limit, offset int, name, surname, patronymic, gender string) ([]entities.User, error) {
	var users []entities.User
	query := `SELECT * FROM users WHERE 
		(name = '' OR name = $1) AND 
		(surname = '' OR surname = $2) AND 
		(patronymic = '' OR patronymic = $3) AND 
		(gender = '' OR gender = $4) 
		LIMIT $5 OFFSET $6`
	err := u.db.Select(&users, query, name, surname, patronymic, gender, limit, offset)
	return users, err
}

func (u *UserRepository) CreateUser(params entities.User) (entities.User, error) {
	query := `INSERT INTO users (name, surname, patronymic, age, gender, nationality, created_at, updated_at) 
	          VALUES (:name, :surname, :patronymic, :age, :gender, :nationality, :created_at, :updated_at)
	          RETURNING id`

	params.Created_at = time.Now()
	params.Updated_at = time.Now()

	stmt, err := u.db.PrepareNamed(query)
	if err != nil {
		return entities.User{}, err
	}
	defer stmt.Close()

	err = stmt.Get(&params.Id, params)
	if err != nil {
		return entities.User{}, err
	}

	return params, nil
}

func (u *UserRepository) ExistByFullName(params entities.FullName) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM users 
		WHERE name = $1 AND surname = $2 AND patronymic = $3)`

	err := u.db.Get(&exists, query, params.Name, params.Surname, params.Patronymic)
	return exists, err
}

func (u *UserRepository) ExistById(id int32) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM users WHERE id = $1)`

	err := u.db.Get(&exists, query, id)
	return exists, err
}

func (u *UserRepository) GetUserById(id int32) (entities.User, error) {
	var user entities.User
	query := `SELECT * FROM users WHERE id = $1`

	err := u.db.Get(&user, query, id)
	return user, err
}

func (u *UserRepository) UpdateUser(params entities.UpdateUserParams) error {
	fields := []string{}
	args := map[string]interface{}{
		"id":         params.Id,
		"updated_at": time.Now(),
	}

	if params.Name != "" {
		fields = append(fields, "name = :name")
		args["name"] = params.Name
	}
	if params.Surname != "" {
		fields = append(fields, "surname = :surname")
		args["surname"] = params.Surname
	}
	if params.Patronymic != "" {
		fields = append(fields, "patronymic = :patronymic")
		args["patronymic"] = params.Patronymic
	}
	if params.Age != 0 {
		fields = append(fields, "age = :age")
		args["age"] = params.Age
	}
	if params.Gender != "" {
		fields = append(fields, "gender = :gender")
		args["gender"] = params.Gender
	}
	if params.Nationality != "" {
		fields = append(fields, "nationality = :nationality")
		args["nationality"] = params.Nationality
	}

	fields = append(fields, "updated_at = :updated_at")

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = :id", strings.Join(fields, ", "))

	_, err := u.db.NamedExec(query, args)
	return err
}

// func (u *UserRepository) UpdateUser(params entities.UpdateUserParams) error {
// 	query := `UPDATE users SET
// 		updated_at = :updated_at,
// 		name = :name,
// 		surname = :surname,
// 		patronymic = :patronymic,
// 		age = :age,
// 		gender = :gender,
// 		nationality = :nationality
// 		WHERE id = :id`

// 	params.Updated_at = time.Now()

// 	_, err := u.db.NamedExec(query, params)
// 	return err
// }

func (u *UserRepository) DeleteUser(id int32) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := u.db.Exec(query, id)
	return err
}
