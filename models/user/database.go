package user

import (
	"database/sql"
	"fmt"
	"oui/models/postgresql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

func GetSQLUserToken(Email, Password string) (token UserToken, err error) {

	var (
		id    sql.NullInt64
		email sql.NullString
	)

	query := "select id, email from account " +
		"where email=$1 and password=crypt($2, password)"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	err = sqlCo.QueryRow(postgresql.SQLCtx, query, Password).Scan(
		&id,
		&email,
	)

	if err == pgx.ErrNoRows {
		golog.Infof("Tentative de connexion infructueuse pour l'utilisateur : %v. Email inconnu ou utilisateur non actif.", Email)
		return
	} else if err != nil {
		golog.Error(query, err)
		return
	}

	token = UserToken{
		ID:        id.Int64,
		Email:     email.String,
		CreatedAt: time.Now(),
	}

	return
}

func CreateAccount(firstname, lastname, email, password string) (id int64) {

	query := "insert into account (id, firstname, lastname, email, password) " +
		"VALUES ($1::bigint,$2,$3,$4,crypt($5, gen_salt('bf')))"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return -1
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	id = time.Now().UnixNano()

	_, err = sqlCo.Exec(postgresql.SQLCtx, query, id, firstname, lastname, email, password)
	if err != nil {
		return -1
	}
	return id
}

func SecurityCheck(id int64, password string) (checked bool) {
	query := "select id from account where id=$1 and password=crypt($2, password)"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var id_ int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, id, password).Scan(&id_)
	if err == nil {
		return id_ == id
	}
	return
}

func CheckEmailAvailability(email string) (available bool) {
	query := "select id from account where email=$1"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var id int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, email).Scan(&id)
	golog.Debug(err, id)
	if err == pgx.ErrNoRows {
		return true
	}
	return
}

func GetUserById(UserId int64) (u User, err error) {

	var (
		email, firstname, lastname, taxiToken, resetToken, password sql.NullString
	)

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	// Requête GetUserById
	var query = "SELECT email, firstname, lastname, taxi_token, reset_token, password" +
		"FROM account " +
		"WHERE id=$1 "

	err = sqlCo.QueryRow(postgresql.SQLCtx, query, UserId).Scan(
		&email,
		&firstname,
		&lastname,
		&taxiToken,
		&resetToken,
		&password,
	)

	if err != nil {
		return
	}

	u = User{
		ID:         UserId,
		Firstname:  firstname.String,
		Lastname:   lastname.String,
		Email:      email.String,
		Password:   password.String,
		TaxiToken:  taxiToken.String,
		ResetToken: resetToken.String,
	}

	return
}

func GetUserByTaxiToken(taxiToken string) (u User, err error) {

	var (
		id                                     sql.NullInt64
		email, firstname, lastname, resetToken sql.NullString
	)

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var query = "SELECT id, email, firstname, lastname, reset_token " +
		"FROM account " +
		"WHERE taxi_token=$1 "

	err = sqlCo.QueryRow(postgresql.SQLCtx, query, taxiToken).Scan(
		&id,
		&email,
		&firstname,
		&lastname,
		&resetToken,
	)

	if err != nil {
		return
	}

	u = User{
		ID:         id.Int64,
		Firstname:  firstname.String,
		Lastname:   lastname.String,
		Email:      email.String,
		TaxiToken:  taxiToken,
		ResetToken: resetToken.String,
	}

	return
}

func GetAllUsers() (users []User) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}

	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT id, email, firstname, lastname, reset_token, password FROM account"

	rows, err := sqlCo.Query(postgresql.SQLCtx, query)
	if err != nil {
		golog.Errorf("execution query '%s':\n%s", query, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id                                                sql.NullInt64
			email, firstname, lastname, reset_token, password sql.NullString
		)

		err = rows.Scan(
			&id,
			&email,
			&firstname,
			&lastname,
			&reset_token,
			&password,
		)

		if err == pgx.ErrNoRows {
			return
		} else if err != nil {
			golog.Errorf("psql request '%v' failed with error : %v", query, err)
			return
		}

		users = append(users, User{
			ID:         id.Int64,
			Firstname:  firstname.String,
			Lastname:   lastname.String,
			Email:      email.String,
			Password:   password.String,
			ResetToken: reset_token.String,
		})
	}

	return
}

func GetUserByEmail(userEmail string) (u User, msg string) {

	if userEmail == "" {
		msg = "Empty email"
		return
	}

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		msg = "Internal server error"
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT id, firstname, lastname, password, taxi_token, reset_token " +
		"FROM account " +
		"where email=$1"

	var (
		id                                                   sql.NullInt64
		firstname, lastname, taxiToken, resetToken, password sql.NullString
	)
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, userEmail).Scan(
		&id,
		&firstname,
		&lastname,
		&password,
		&taxiToken,
		&resetToken,
	)

	if err == pgx.ErrNoRows {
		msg = "Username not found"
		return
	} else if err != nil {
		golog.Errorf("psql request '%v' failed with error : %v", query, err)
		msg = "Internal server error"
		return
	}

	u = User{
		ID:         id.Int64,
		Firstname:  firstname.String,
		Lastname:   lastname.String,
		Email:      userEmail,
		Password:   password.String,
		TaxiToken:  taxiToken.String,
		ResetToken: resetToken.String,
	}

	return
}

func UpdateUser(user User, password bool) (ok bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var query string
	var args []interface{}
	if password {
		query = "UPDATE account set (firstname, lastname, email, password) = ($1,$2,$3,crypt($4, gen_salt('bf'))) " +
			"WHERE id=$5"
		args = []interface{}{
			user.Firstname,
			user.Lastname,
			user.Email,
			user.Password,
			user.ID,
		}
	} else {
		query = "UPDATE account set (firstname, lastname, email) = ($1,$2,$3) " +
			"WHERE id=$4"
		args = []interface{}{
			user.Firstname,
			user.Lastname,
			user.Email,
			user.ID,
		}
	}

	cmd, err := sqlCo.Exec(postgresql.SQLCtx, query, args...)
	ok = cmd.RowsAffected() == 1 && err == nil

	return
}

func GenResetToken(email string) (token string) {

	if email == "" {
		return
	}

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	tempUUID := uuid.New().String()

	query := "UPDATE account " +
		"SET reset_token=$1 " +
		"WHERE email=$2"

	updateSQLCmd, err := sqlCo.Exec(postgresql.SQLCtx, query, tempUUID, email)
	if rowsAffected := updateSQLCmd.RowsAffected(); err != nil {
		golog.Error(query, err)
	} else if rowsAffected == 0 {
		golog.Infof("Tentative de génération de token infructueuse. L'utilisateur %v n'existe pas ou n'est pas actif !", email)
	} else if rowsAffected == 1 {
		token = tempUUID
	}

	return
}

func GenTaxiToken(userId int64) (token string) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	tempUUID := uuid.New().String()

	query := "UPDATE account " +
		"SET taxi_token=$1 " +
		"WHERE id=$2"

	updateSQLCmd, err := sqlCo.Exec(postgresql.SQLCtx, query, tempUUID, userId)
	if rowsAffected := updateSQLCmd.RowsAffected(); err != nil {
		golog.Error(query, err)
	} else if rowsAffected == 0 {
		golog.Infof("Tentative de génération de token CariTaxi infructueuse. L'utilisateur %d n'existe pas ou n'est pas actif !", userId)
	} else if rowsAffected == 1 {
		token = tempUUID
	}

	return
}

func DefinePasswordWithResetToken(ResetToken, Password string) (success bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	// DefinePasswordWithResetToken: UPDATE PASSWORD
	query := "UPDATE account " +
		"SET password=crypt($1, gen_salt('bf')), reset_token=$3 " +
		"WHERE reset_token=$2"

	cmd, err := sqlCo.Exec(postgresql.SQLCtx, query, Password, ResetToken, sql.NullString{})
	if rowsAffected := cmd.RowsAffected(); err != nil {
		golog.Error(query, err)
	} else if rowsAffected == 0 {
		golog.Warnf("Tentative de réinitialisation de mot de passe infructueuse. Le token %v n'existe pas !", ResetToken)
	} else {
		success = true
	}

	return
}

func GetUserPermissions(account int64) (permissions []string) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}

	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT permission FROM account_permission WHERE account=$1"
	rows, err := sqlCo.Query(postgresql.SQLCtx, query, account)
	if err != nil {
		golog.Errorf("execution query '%s':\n%s", query, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err == pgx.ErrNoRows {
			return
		} else if err != nil {
			golog.Errorf("psql scan '%v' failed with error : %v", query, err)
			return
		}
		permissions = append(permissions, permission)
	}
	return
}

func UpdateUserPermissions(user int64, permissions []string) (msg string) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return "Internal server error"
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	sqlCo.Exec(postgresql.SQLCtx, "DELETE FROM account_permission WHERE account=$1", user)

	if len(permissions) > 0 {
		var values []string
		var args []interface{}

		for i, v := range permissions {
			values = append(values, fmt.Sprintf("($%d,$%d)", 2*i+1, 2*(i+1)))
			args = append(args, user, v)
		}

		query := "INSERT INTO account_permission(account,permission) VALUES " + strings.Join(values, ",")

		_, err := sqlCo.Exec(postgresql.SQLCtx, query, args...)
		if err != nil {
			return "Unable to insert user permissions"
		}
	}
	return
}
