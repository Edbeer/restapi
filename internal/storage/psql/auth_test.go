package psql

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Edbeer/restapi/internal/entity"
	"github.com/Edbeer/restapi/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestPsql_Register(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authStorage := NewAuthStorage(sqlxDB)

	t.Run("Regoster", func(t *testing.T) {

		role := "admin"
		columns := []string{
			"first_name",
			"last_name",
			"email",
			"password",
			"role"}
		rows := sqlmock.NewRows(columns).AddRow(
			"Pavel",
			"Volkov",
			"edbeermtn@gmail.com",
			"d2345678",
			"admin",
		)

		user := &entity.User{
			FirstName: "Pavel",
			LastName:  "Volkov",
			Email:     "edbeermtn@gmail.com",
			Password:  "d2345678",
			Role:      &role,
		}

		mock.ExpectQuery(createUserQuery).WithArgs(
			&user.FirstName, &user.LastName, &user.Email,
			&user.Password, &user.Role, &user.Avatar,
			&user.PhoneNumber, &user.Address, &user.City,
			&user.Country, &user.Postcode).WillReturnRows(rows)

		createdUser, err := authStorage.Register(context.Background(), user)
		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, user)
	})
}

func TestPsql_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authStorage := NewAuthStorage(sqlxDB)

	t.Run("Update", func(t *testing.T) {
		role := "admin"
		columns := []string{
			"first_name",
			"last_name",
			"email",
			"password",
			"role"}
		rows := sqlmock.NewRows(columns).AddRow(
			"Pavel",
			"Volkov",
			"edbeermtn@gmail.com",
			"d2345678",
			"admin",
		)

		user := &entity.User{
			FirstName: "Pavel",
			LastName:  "Volkov",
			Email:     "edbeermtn@gmail.com",
			Password:  "d2345678",
			Role:      &role,
		}

		mock.ExpectQuery(updateUserQuery).WithArgs(
			&user.FirstName, &user.LastName, &user.Email,
			&user.Role, &user.Avatar, &user.PhoneNumber,
			&user.Address, &user.City, &user.Country,
			&user.Postcode, &user.ID,
		).WillReturnRows(rows)

		updatedUser, err := authStorage.Update(context.Background(), user)
		require.NoError(t, err)
		require.NotNil(t, updatedUser)
		require.Equal(t, user, updatedUser)
	})
}

func TestPsql_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authStorage := NewAuthStorage(sqlxDB)

	t.Run("Delete", func(t *testing.T) {
		uid := uuid.New()

		mock.ExpectExec(deleteUserQuery).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 1))

		err := authStorage.Delete(context.Background(), uid)
		require.NoError(t, err)
	})

	t.Run("Delete no rows", func(t *testing.T) {
		uid := uuid.New()

		mock.ExpectExec(deleteUserQuery).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 0))

		err := authStorage.Delete(context.Background(), uid)
		require.NotNil(t, err)
	})
}

func TestPsql_GetUserByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authStorage := NewAuthStorage(sqlxDB)

	t.Run("GetUserByID", func(t *testing.T) {
		uid := uuid.New()

		columns := []string{
			"user_id",
			"first_name",
			"last_name",
			"email",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			uid,
			"Pavel",
			"Volkov",
			"edbeermtn@gmail.com",
		)

		testUser := &entity.User{
			ID:        uid,
			FirstName: "Pavel",
			LastName:  "Volkov",
			Email:     "edbeermtn@gmail.com",
		}

		mock.ExpectQuery(getUserByID).WithArgs(uid).WillReturnRows(rows)

		user, err := authStorage.GetUserByID(context.Background(), uid)
		require.NoError(t, err)
		require.Equal(t, user.FirstName, testUser.FirstName)
		fmt.Printf("test user: %s \n", testUser.FirstName)
		fmt.Printf("user: %s \n", user.FirstName)
	})
}

func TestPsql_FindUserByName(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authStorage := NewAuthStorage(sqlxDB)

	t.Run("FindUserByName", func(t *testing.T) {
		uid := uuid.New()
		userName := "Pavel"

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		columns := []string{
			"user_id",
			"first_name",
			"last_name",
			"email",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			uid,
			"Pavel",
			"Volkov",
			"edbeermtn@gmail.com",
		)

		mock.ExpectQuery(getTotalCount).WillReturnRows(totalCountRows)
		mock.ExpectQuery(findUsersByName).WithArgs("", 0, 10).WillReturnRows(rows)

		userList, err := authStorage.FindUsersByName(context.Background(), userName, &utils.PaginationQuery{
			Size:    10,
			Page:    0,
			OrderBy: "",
		})

		require.NoError(t, err)
		require.NotNil(t, userList)
	})
}

func TestPsql_GetUsers(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authStorage := NewAuthStorage(sqlxDB)

	t.Run("GetUsers", func(t *testing.T) {
		uid := uuid.New()

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		columns := []string{
			"user_id",
			"first_name",
			"last_name",
			"email",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			uid,
			"Pavel",
			"Volkov",
			"edbeermtn@gmail.com",
		)

		mock.ExpectQuery(getTotal).WillReturnRows(totalCountRows)
		mock.ExpectQuery(getUsers).WithArgs("", 0, 10).WillReturnRows(rows)

		userList, err := authStorage.GetUsers(context.Background(), &utils.PaginationQuery{
			Size:    10,
			Page:    0,
			OrderBy: "",
		})

		require.NoError(t, err)
		require.NotNil(t, userList)
	})
}

func TestPsql_GetUserByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authStorage := NewAuthStorage(sqlxDB)

	t.Run("GetUserByEmail", func(t *testing.T) {
		uid := uuid.New()

		columns := []string{
			"user_id",
			"first_name",
			"last_name",
			"email",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			uid,
			"Pavel",
			"Volkov",
			"edbeermtn@gmail.com",
		)

		testUser := &entity.User{
			ID:        uid,
			FirstName: "Pavel",
			LastName:  "Volkov",
			Email:     "edbeermtn@gmail.com",
		}

		mock.ExpectQuery(findUserByEmail).WithArgs(&testUser.Email).WillReturnRows(rows)

		user, err := authStorage.FindUserByEmail(context.Background(), testUser)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, user.FirstName, testUser.FirstName)
	})
}
