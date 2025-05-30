package repository

import (
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
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
	builder := sq.Select("*").
		From("users").
		Where("1=1").
		PlaceholderFormat(sq.Dollar)

	if name != "" {
		builder = builder.Where(sq.ILike{"name": "%" + name + "%"})
	}
	if surname != "" {
		builder = builder.Where(sq.ILike{"surname": "%" + surname + "%"})
	}
	if patronymic != "" {
		builder = builder.Where(sq.ILike{"patronymic": "%" + patronymic + "%"})
	}
	if gender != "" {
		builder = builder.Where(sq.Eq{"gender": gender})
	}
	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	}
	if offset > 0 {
		builder = builder.Offset(uint64(offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var users []entities.User
	err = u.db.Select(&users, query, args...)
	return users, err
}

func (u *UserRepository) CreateUser(params entities.User) (entities.User, error) {
	params.Created_at = time.Now()
	params.Updated_at = time.Now()

	builder := sq.Insert("users").
		Columns("name", "surname", "patronymic", "age", "gender", "nationality", "created_at", "updated_at").
		Values(params.Name, params.Surname, params.Patronymic, params.Age, params.Gender, params.Nationality, params.Created_at, params.Updated_at).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return entities.User{}, err
	}

	err = u.db.Get(&params.Id, query, args...)
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
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`

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
	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	builder := sq.Update("users").Where(sq.Eq{"id": id}).Set("updated_at", time.Now()).PlaceholderFormat(sq.Dollar)

	if params.Name != nil {
		builder = builder.Set("name", *params.Name)
	}
	if params.Surname != nil {
		builder = builder.Set("surname", *params.Surname)
	}
	if params.Patronymic != nil {
		builder = builder.Set("patronymic", *params.Patronymic)
	}
	if params.Age != nil {
		builder = builder.Set("age", *params.Age)
	}
	if params.Gender != nil {
		builder = builder.Set("gender", *params.Gender)
	}
	if params.Nationality != nil {
		builder = builder.Set("nationality", *params.Nationality)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *UserRepository) DeleteUser(id int32) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := u.db.Exec(query, id)
	return err
}
