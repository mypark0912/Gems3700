package cfgroute

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func openMaintenanceDB(configDir string) (*sql.DB, error) {
	dbPath := filepath.Join(configDir, "maintenance.db")
	isNew := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		isNew = true
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	if isNew {
		db.Exec(`CREATE TABLE IF NOT EXISTS maintenance (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT, context TEXT, mtype INTEGER, utype TEXT,
			f_version TEXT, a_version TEXT, w_version TEXT, c_version TEXT, smart_version TEXT,
			date TEXT
		)`)
	}
	return db, nil
}

func openLogDB(configDir string) (*sql.DB, error) {
	dbPath := filepath.Join(configDir, "log.db")
	isNew := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		isNew = true
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	if isNew {
		db.Exec(`CREATE TABLE IF NOT EXISTS log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			logdate TEXT, account TEXT, userRole TEXT, action TEXT
		)`)
	}
	return db, nil
}

type Post struct {
	ID           int    `json:"id,omitempty"`
	Title        string `json:"title"`
	Context      string `json:"context"`
	Mtype        int    `json:"mtype"`
	Utype        string `json:"utype"`
	FVersion     string `json:"f_version"`
	AVersion     string `json:"a_version"`
	WVersion     string `json:"w_version"`
	CVersion     string `json:"c_version"`
	SmartVersion string `json:"smart_version"`
	Date         string `json:"date"`
}

func savePost(db *sql.DB, post Post, mode, idx int) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	if mode == 0 {
		_, err := db.Exec(
			`INSERT INTO maintenance (title, context, mtype, utype, f_version, a_version, w_version, c_version, smart_version, date)
			 VALUES (?,?,?,?,?,?,?,?,?,?)`,
			post.Title, post.Context, post.Mtype, post.Utype,
			post.FVersion, post.AVersion, post.WVersion, post.CVersion, post.SmartVersion, now)
		return err
	}
	_, err := db.Exec(
		`UPDATE maintenance SET title=?, context=?, mtype=?, utype=?, f_version=?, a_version=?, w_version=?, c_version=?, smart_version=?, date=? WHERE id=?`,
		post.Title, post.Context, post.Mtype, post.Utype,
		post.FVersion, post.AVersion, post.WVersion, post.CVersion, post.SmartVersion, now, idx)
	return err
}

func getAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT id, title, context, mtype, utype, f_version, a_version, w_version, c_version, smart_version, date FROM maintenance")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.Title, &p.Context, &p.Mtype, &p.Utype,
			&p.FVersion, &p.AVersion, &p.WVersion, &p.CVersion, &p.SmartVersion, &p.Date)
		posts = append(posts, p)
	}
	return posts, nil
}

func getLastPost(db *sql.DB) (*Post, error) {
	row := db.QueryRow("SELECT id, title, context, mtype, utype, f_version, a_version, w_version, c_version, smart_version, date FROM maintenance ORDER BY id DESC LIMIT 1")
	var p Post
	err := row.Scan(&p.ID, &p.Title, &p.Context, &p.Mtype, &p.Utype,
		&p.FVersion, &p.AVersion, &p.WVersion, &p.CVersion, &p.SmartVersion, &p.Date)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

type LogEntry struct {
	ID       int    `json:"id"`
	LogDate  string `json:"logdate"`
	Account  string `json:"account"`
	UserRole string `json:"userRole"`
	Action   string `json:"action"`
}

func getLogsPaginated(db *sql.DB, page, pageSize int) ([]LogEntry, int, error) {
	var total int
	db.QueryRow("SELECT COUNT(*) FROM log").Scan(&total)

	offset := (page - 1) * pageSize
	rows, err := db.Query("SELECT id, logdate, account, userRole, action FROM log ORDER BY id DESC LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var l LogEntry
		rows.Scan(&l.ID, &l.LogDate, &l.Account, &l.UserRole, &l.Action)
		logs = append(logs, l)
	}
	return logs, total, nil
}

func insertLog(db *sql.DB, account, role, action string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec("INSERT INTO log (logdate, account, userRole, action) VALUES (?,?,?,?)",
		now, account, role, action)
	return err
}

func updateLog(db *sql.DB, account, role, action string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	// Check today's entry
	today := time.Now().Format("2006-01-02")
	var id int
	err := db.QueryRow("SELECT id FROM log WHERE logdate LIKE ? AND account=? AND action=?",
		today+"%", account, action).Scan(&id)
	if err == nil {
		// Update existing
		_, err = db.Exec("UPDATE log SET logdate=? WHERE id=?", now, id)
		return err
	}
	// Insert new
	return insertLog(db, account, role, action)
}

func deletePost(db *sql.DB, idx string) error {
	_, err := db.Exec("DELETE FROM maintenance WHERE id=?", idx)
	return err
}

func getPostCount(db *sql.DB) int {
	var c int
	db.QueryRow("SELECT COUNT(*) FROM maintenance").Scan(&c)
	return c
}

func init() {
	// Ensure table names don't have issues
	_ = fmt.Sprintf
}
