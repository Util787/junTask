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

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) GetAllUsers(pageSize, page int, name, surname, patronymic, gender string) ([]entities.User, int, error) {
	totalCountBuilder := sq.Select("COUNT(*)").From("users").Where("1=1").PlaceholderFormat(sq.Dollar)
	usersBuilder := sq.Select("*").From("users").Where("1=1").PlaceholderFormat(sq.Dollar)

	if name != "" {
		nameLike := sq.ILike{"name": name + "%"}
		usersBuilder = usersBuilder.Where(nameLike)
		totalCountBuilder = totalCountBuilder.Where(nameLike)
	}
	if surname != "" {
		surnameLike := sq.ILike{"surname": surname + "%"}
		usersBuilder = usersBuilder.Where(surnameLike)
		totalCountBuilder = totalCountBuilder.Where(surnameLike)
	}
	if patronymic != "" {
		patronymicLike := sq.ILike{"patronymic": patronymic + "%"}
		usersBuilder = usersBuilder.Where(patronymicLike)
		totalCountBuilder = totalCountBuilder.Where(patronymicLike)
	}
	if gender != "" {
		genderEq := sq.Eq{"gender": gender}
		usersBuilder = usersBuilder.Where(genderEq)
		totalCountBuilder = totalCountBuilder.Where(genderEq)
	}

	offset := (page - 1) * pageSize
	usersBuilder = usersBuilder.Limit(uint64(pageSize)).Offset(uint64(offset))

	usersQuery, usersArgs, err := usersBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	totalCountQuery, totalCountArgs, err := totalCountBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	tx, err := u.db.Beginx()
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	var users []entities.User
	err = tx.Select(&users, usersQuery, usersArgs...)
	if err != nil {
		return nil, 0, err
	}

	var totalCount int
	err = tx.Get(&totalCount, totalCountQuery, totalCountArgs...)
	if err != nil {
		return nil, 0, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (u *userRepository) CreateUser(params entities.User) (entities.User, error) {
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

func (u *userRepository) ExistByFullName(params entities.FullName) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM users 
		WHERE name = $1 AND surname = $2 AND patronymic = $3)`

	err := u.db.Get(&exists, query, params.Name, params.Surname, params.Patronymic)
	return exists, err
}

func (u *userRepository) ExistById(id int32) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`

	err := u.db.Get(&exists, query, id)
	return exists, err
}

func (u *userRepository) GetUserById(id int32) (entities.User, error) {
	var user entities.User
	query := `SELECT * FROM users WHERE id = $1`

	err := u.db.Get(&user, query, id)
	return user, err
}

func (u *userRepository) UpdateUser(id int32, params entities.UpdateUserParams) error {
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

	_, err = u.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) DeleteUser(id int32) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := u.db.Exec(query, id)
	return err
}
