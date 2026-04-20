package auth

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRow struct {
	Account  string `json:"account"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	API      string `json:"api"`
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hashed, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func ensureUserTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user (
		account TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL,
		role TEXT NOT NULL,
		api TEXT NOT NULL
	)`)
	return err
}

func getUserByAccount(db *sql.DB, account string) (*UserRow, error) {
	row := db.QueryRow("SELECT account, username, password, email, role, api FROM user WHERE account=?", account)
	u := &UserRow{}
	if err := row.Scan(&u.Account, &u.Username, &u.Password, &u.Email, &u.Role, &u.API); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func getAllUsers(db *sql.DB) ([]UserRow, error) {
	rows, err := db.Query("SELECT account, username, password, email, role, api FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserRow
	for rows.Next() {
		var u UserRow
		if err := rows.Scan(&u.Account, &u.Username, &u.Password, &u.Email, &u.Role, &u.API); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func insertUser(db *sql.DB, account, username, password, email, role string) error {
	return insertUserWithAPI(db, account, username, password, email, role, "No")
}

func insertUserWithAPI(db *sql.DB, account, username, password, email, role, api string) error {
	hashed, err := hashPassword(password)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"INSERT INTO user (account, username, password, email, role, api) VALUES (?, ?, ?, ?, ?, ?)",
		account, username, hashed, email, role, api,
	)
	return err
}

func updateUserRole(db *sql.DB, account, role string) error {
	res, err := db.Exec("UPDATE user SET role=? WHERE account=?", role, account)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func deleteUser(db *sql.DB, account string) error {
	_, err := db.Exec("DELETE FROM user WHERE account=?", account)
	return err
}

func userTableExists(db *sql.DB) bool {
	var name string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='user'").Scan(&name)
	return err == nil
}

func userCount(db *sql.DB) int {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	return count
}
