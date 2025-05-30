package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Util787/junTask/entities"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetAllUsers(limit, offset int, name, surname, patronymic, gender string) ([]entities.User, error) {
	query := `SELECT * FROM users WHERE 1=1`
	params := map[string]interface{}{}

	if name != "" {
		query += " AND name ILIKE :name"
		params["name"] = "%" + name + "%"
	}
	if surname != "" {
		query += " AND surname ILIKE :surname"
		params["surname"] = "%" + surname + "%"
	}
	if patronymic != "" {
		query += " AND patronymic ILIKE :patronymic"
		params["patronymic"] = "%" + patronymic + "%"
	}
	if gender != "" {
		query += " AND gender = :gender"
		params["gender"] = gender
	}

	if limit > 0 {
		query += " LIMIT :limit"
		params["limit"] = limit
	}
	if offset > 0 {
		query += " OFFSET :offset"
		params["offset"] = offset
	}

	var users []entities.User
	namedStmt, err := u.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	err = namedStmt.Select(&users, params)
	if err != nil {
		return nil, err
	}

	return users, nil
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

func (u *UserRepository) UpdateUser(id int32, params entities.UpdateUserParams) error {
	fields := []string{}
	args := map[string]interface{}{
		"id":         id,
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

func (u *UserRepository) DeleteUser(id int32) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := u.db.Exec(query, id)
	return err
}
