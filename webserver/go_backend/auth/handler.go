package auth

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"serverGO/crypto"
	"serverGO/infra"
)

type Handler struct {
	db        *sql.DB
	redis     *infra.RedisState
	cipher    *crypto.AESCipher
	configDir string
	baseDir   string
}

func RegisterRoutes(rg *gin.RouterGroup, deps *infra.Dependencies) {
	db, err := infra.EnsureUserDB(deps.Config.ConfigDir)
	if err != nil {
		log.Fatalf("auth: open user.db: %v", err)
	}

	h := &Handler{
		db:        db,
		redis:     deps.Redis,
		cipher:    deps.Crypto,
		configDir: deps.Config.ConfigDir,
		baseDir:   deps.Config.BaseDir,
	}

	rg.GET("/checkInstall", h.checkInstall)
	rg.GET("/checkDBMS", h.checkDBMS)
	rg.POST("/joinAdmin", h.joinAdmin)
	rg.GET("/checkSession", h.checkSession)
	rg.GET("/getUser", h.getUser)
	rg.GET("/getUserList", h.getUserList)
	rg.POST("/resetPassword", h.resetPassword)
	rg.GET("/saveUser/:account/:role", h.saveUser)
	rg.GET("/removeUser/:username", h.removeUser)
	rg.POST("/updateProfile", h.updateProfile)
	rg.POST("/join", h.join)
	rg.GET("/checkAccount", h.checkAccount)
	rg.POST("/checkLogins", h.checkLogins)
	rg.GET("/checkRemote/:user", h.checkRemote)
	rg.POST("/verify-admin", h.verifyAdmin)
	rg.GET("/logout", h.logout)
}

func (h *Handler) checkInstall(c *gin.Context) {
	ensureUserTable(h.db)

	hasUsers := false
	if userTableExists(h.db) {
		hasUsers = userCount(h.db) > 0
	}

	calPath := filepath.Join(h.configDir, "calibration.csv")
	calibrated := fileExists(calPath)

	// InfluxDB is no longer used — treat DBMS as always initialized when users exist.
	if hasUsers {
		c.JSON(http.StatusOK, gin.H{"result": 2, "calibration": calibrated})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": 0})
}

func (h *Handler) checkDBMS(c *gin.Context) {
	// InfluxDB removed — SQLite-only backend is always ready.
	c.JSON(http.StatusOK, gin.H{"result": 1})
}

func (h *Handler) joinAdmin(c *gin.Context) {
	var req SignupAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "msg": err.Error()})
		return
	}

	if !h.cipher.CheckAdmin(req.AdminPass) {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "msg": "Admin Password is Wrong"})
		return
	}

	devType := 0
	switch v := req.DevType.(type) {
	case float64:
		devType = int(v)
	case string:
		fmt.Sscanf(v, "%d", &devType)
	}
	mode := fmt.Sprintf("device%d", devType)
	diag := "No"
	if devType > 0 && devType <= 4 {
		diag = "Yes"
	}
	if devType > 4 {
		mode = "server"
	}

	// Admin user (role=3) with API flag matching diagnosis usage.
	if err := insertUserWithAPI(h.db, req.Account, req.Username, req.Password, req.Email, "3", diag); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "msg": err.Error()})
		return
	}

	// Default client accounts (role 2 = admin, role 0 = guest).
	insertUser(h.db, "client_admin", "client_admin", "1234", "ntek@nteksys.com", "2")
	insertUser(h.db, "client_guest", "client_guest", "1234", "ntek@nteksys.com", "0")

	// Save first-install maintenance post with current versions.
	saveFirstInstallPost(h.configDir, h.baseDir, mode)

	// Copy default.json → setup.json and configure.
	setupPath := filepath.Join(h.configDir, "setup.json")
	defaultPath := filepath.Join(h.configDir, "default.json")
	mac := getMACAddressClean()
	ip := getLocalIP()

	if data, err := os.ReadFile(defaultPath); err == nil {
		var setting map[string]interface{}
		if json.Unmarshal(data, &setting) == nil {
			setting["mode"] = mode
			setting["lang"] = req.Lang
			if gen, ok := setting["General"].(map[string]interface{}); ok {
				if di, ok := gen["deviceInfo"].(map[string]interface{}); ok {
					di["mac_address"] = mac
					di["serial_number"] = mac
				}
				if ip != "0.0.0.0" {
					if tcp, ok := gen["tcpip"].(map[string]interface{}); ok {
						tcp["ip_address"] = ip
					}
				}
			}
			out, _ := json.MarshalIndent(setting, "", "  ")
			os.WriteFile(setupPath, out, 0644)

			ctx := context.Background()
			h.redis.Client0.HSet(ctx, "System", "setup", string(out))
			h.redis.Client0.HSet(ctx, "System", "mode", mode)

			// Trigger service save/restart for A35 / Core.
			if isServiceActiveAuth("sv500A35") {
				h.redis.Client0.HSet(ctx, "Service", "save", 1)
			} else {
				_ = exec.Command("sudo", "systemctl", "start", "sv500A35").Run()
			}
			if isServiceActiveAuth("core") {
				h.redis.Client0.HSet(ctx, "Service", "restart", 1)
			} else {
				_ = exec.Command("sudo", "systemctl", "start", "core").Run()
			}

			// Flush meter DB to discard stale readings tied to previous setup.
			h.redis.Client1.FlushDB(ctx)
		}
	}

	// Initialize bearing DB from NTEKBearingDB.csv.
	if res := initBearingDBFromCSV(h.configDir); !res.success {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "msg": res.msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

func (h *Handler) checkSession(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"loggedIn": false, "username": nil, "userRole": nil, "mode": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"loggedIn": true,
		"username": user,
		"userRole": session.Get("userRole"),
		"mode":     session.Get("devMode"),
	})
}

func (h *Handler) getUser(c *gin.Context) {
	session := sessions.Default(c)
	account, ok := session.Get("user").(string)
	if !ok || account == "" {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}
	u, err := getUserByAccount(h.db, account)
	if err != nil || u == nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"passOK": "1",
		"data":   gin.H{"username": u.Username, "email": u.Email, "api": u.API},
	})
}

func (h *Handler) getUserList(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("user") == nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}
	users, err := getAllUsers(h.db)
	if err != nil || len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"passOK": "1", "data": users})
}

func (h *Handler) resetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}

	u, err := getUserByAccount(h.db, req.Username)
	if err != nil || u == nil || u.Email != req.Email {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}

	hashed, _ := hashPassword(req.Password)
	h.db.Exec("UPDATE user SET password=? WHERE account=? AND email=?", hashed, req.Username, req.Email)
	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

func (h *Handler) saveUser(c *gin.Context) {
	account := c.Param("account")
	role := c.Param("role")
	if err := updateUserRole(h.db, account, role); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

func (h *Handler) removeUser(c *gin.Context) {
	username := c.Param("username")
	if username == "admin" {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}
	deleteUser(h.db, username)
	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

func (h *Handler) updateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}

	session := sessions.Default(c)
	account, _ := session.Get("user").(string)
	if account == "" {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}

	u, err := getUserByAccount(h.db, account)
	if err != nil || u == nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}

	if u.Role == "3" {
		// Admin: update without password check
		h.db.Exec("UPDATE user SET username=?, email=?, api=? WHERE account=?",
			req.Username, req.Email, req.API, account)
		c.JSON(http.StatusOK, gin.H{"passOK": "1"})
		return
	}

	if !checkPassword(u.Password, req.Password) {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}

	hashed, _ := hashPassword(req.NewPass)
	h.db.Exec("UPDATE user SET username=?, password=?, email=?, api=? WHERE account=?",
		req.Username, hashed, req.Email, req.API, account)
	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

func (h *Handler) join(c *gin.Context) {
	var req SignupUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}
	if req.Account == "" {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}
	if err := insertUser(h.db, req.Account, req.Username, req.Password, req.Email, req.Role); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

func (h *Handler) checkAccount(c *gin.Context) {
	if !userTableExists(h.db) {
		c.JSON(http.StatusOK, gin.H{"result": false})
		return
	}

	admin, _ := getUserByAccount(h.db, "client_admin")
	guest, _ := getUserByAccount(h.db, "client_guest")

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"exist": gin.H{
			"ntek":  true,
			"admin": admin != nil,
			"guest": guest != nil,
		},
	})
}

func (h *Handler) checkLogins(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Account == "" {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}

	ctx := context.Background()
	u, err := getUserByAccount(h.db, req.Account)
	if err != nil || u == nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "5"})
		return
	}

	if !checkPassword(u.Password, req.Password) {
		c.JSON(http.StatusOK, gin.H{"passOK": "4"})
		return
	}

	session := sessions.Default(c)
	session.Set("user", req.Account)
	session.Set("lang", req.Lang)
	session.Set("userRole", u.Role)

	// Get mode from Redis
	mode := ""
	if setupJSON, err := h.redis.Client0.HGet(ctx, "System", "setup").Result(); err == nil {
		var setup map[string]interface{}
		if json.Unmarshal([]byte(setupJSON), &setup) == nil {
			mode, _ = setup["mode"].(string)
		}
	}

	if storedMode, err := h.redis.Client0.HGet(ctx, "System", "mode").Result(); err == nil {
		if storedMode != mode && mode != "" {
			h.redis.Client0.HSet(ctx, "System", "mode", mode)
		}
	}

	session.Set("devMode", mode)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"passOK": "1",
		"data":   gin.H{"lang": req.Lang, "userRole": u.Role},
		"mode":   mode,
	})
}

func (h *Handler) checkRemote(c *gin.Context) {
	user := c.Param("user")
	if user == "" {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	ctx := context.Background()
	mode := "device0"
	if m, err := h.redis.Client0.HGet(ctx, "System", "mode").Result(); err == nil {
		mode = m
	}

	session := sessions.Default(c)

	if user == "ntek" {
		session.Set("user", user)
		session.Set("userRole", 3)
		session.Set("devMode", mode)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"mode":    mode,
			"data":    gin.H{"account": user, "userRole": 3, "mode": mode},
		})
		return
	}

	u, err := getUserByAccount(h.db, user)
	if err != nil || u == nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	session.Set("user", user)
	session.Set("userRole", 2)
	session.Set("devMode", mode)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"mode":    mode,
		"data":    gin.H{"account": user, "userRole": 2, "mode": mode},
	})
}

func (h *Handler) verifyAdmin(c *gin.Context) {
	session := sessions.Default(c)
	user, _ := session.Get("user").(string)
	userRole := fmt.Sprintf("%v", session.Get("userRole"))

	if user == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Not logged in"})
		return
	}
	if userRole != "2" && userRole != "3" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Permission denied"})
		return
	}

	var req VerifyAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}

	u, err := getUserByAccount(h.db, user)
	if err != nil || u == nil || !checkPassword(u.Password, req.Password) {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// helpers

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// getMACAddressClean returns MAC lowercased with all separators stripped — matches
// Python get_mac_address. Prefers sw0ep / end1 before falling back to any non-loopback.
func getMACAddressClean() string {
	clean := func(addr string) string {
		m := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(addr, ":", ""), "-", ""))
		if m == "" || m == "000000000000" {
			return ""
		}
		return m
	}

	ifaces, err := net.Interfaces()
	if err == nil {
		for _, name := range []string{"sw0ep", "end1"} {
			for _, i := range ifaces {
				if i.Name == name && i.HardwareAddr != nil {
					if m := clean(i.HardwareAddr.String()); m != "" {
						return m
					}
				}
			}
		}
		for _, i := range ifaces {
			if i.Flags&net.FlagLoopback == 0 && i.HardwareAddr != nil {
				if m := clean(i.HardwareAddr.String()); m != "" {
					return m
				}
			}
		}
	}
	return "000000000000"
}

// getLocalIP returns the first IPv4 address of sw0ep/end1, falling back to any non-loopback.
func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "0.0.0.0"
	}
	pickIP := func(i net.Interface) string {
		addrs, err := i.Addrs()
		if err != nil {
			return ""
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if ip != "127.0.0.1" {
					return ip
				}
			}
		}
		return ""
	}
	for _, name := range []string{"sw0ep", "end1"} {
		for _, i := range ifaces {
			if i.Name == name {
				if ip := pickIP(i); ip != "" {
					return ip
				}
			}
		}
	}
	for _, i := range ifaces {
		if i.Flags&net.FlagLoopback == 0 {
			if ip := pickIP(i); ip != "" {
				return ip
			}
		}
	}
	return "0.0.0.0"
}

// isServiceActiveAuth runs `systemctl is-active <name>`.
func isServiceActiveAuth(name string) bool {
	out, err := exec.Command("systemctl", "is-active", name).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "active"
}

// readVersionInfo parses /home/root/versionInfo.txt (KEY=VALUE per line).
func readVersionInfo() map[string]string {
	data, err := os.ReadFile("/home/root/versionInfo.txt")
	if err != nil {
		return nil
	}
	out := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), "=", 2)
		if len(parts) == 2 {
			out[parts[0]] = parts[1]
		}
	}
	return out
}

// getLatestBuildVersion returns the last .md filename (without ext) in release_notes/ko.
func getLatestBuildVersion(baseDir string) string {
	dir := filepath.Join(baseDir, "release_notes", "ko")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var versions []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		versions = append(versions, strings.TrimSuffix(e.Name(), ".md"))
	}
	if len(versions) == 0 {
		return ""
	}
	// sorted ascending → last is latest.
	return versions[len(versions)-1]
}

// saveFirstInstallPost writes a first-install maintenance.db row matching Python getVersionSave.
func saveFirstInstallPost(configDir, baseDir, mode string) {
	dbPath := filepath.Join(configDir, "maintenance.db")
	os.MkdirAll(configDir, 0755)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	db.Exec(`CREATE TABLE IF NOT EXISTS maintenance (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL, context TEXT NOT NULL,
		mtype INT NOT NULL, utype TEXT NOT NULL,
		f_version TEXT, a_version TEXT, w_version TEXT, c_version TEXT,
		smart_version TEXT, build_version TEXT, date TEXT
	)`)

	versions := readVersionInfo()
	build := getLatestBuildVersion(baseDir)

	getv := func(key, def string) string {
		if v, ok := versions[key]; ok && v != "" {
			return v
		}
		return def
	}

	utype := "fw,a35,web,core"
	smart := ""
	if mode != "device0" {
		utype = "fw,a35,web,core,smartsystem"
		smart = getv("smartsystem", "1.0.0")
	}

	today := time.Now().Format("2006-01-02")
	db.Exec(
		`INSERT INTO maintenance (title, context, mtype, utype, f_version, a_version, w_version, c_version, smart_version, build_version, date) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		"Fist Installation", "SV-500 Installed", 0, utype,
		getv("fw", "1.0.0"), getv("a35", "1.0.0"), getv("web", "1.0.0"), getv("core", "1.0.0"),
		smart, build, today,
	)
}

type bearingInitResult struct {
	success  bool
	msg      string
	inserted int
	skipped  int
}

func parseBearingCSV(r *os.File) ([]map[string]string, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 1 {
		return nil, nil
	}
	header := records[0]
	out := make([]map[string]string, 0, len(records)-1)
	for _, rec := range records[1:] {
		m := make(map[string]string, len(header))
		for i, h := range header {
			if i < len(rec) {
				m[h] = rec[i]
			}
		}
		out = append(out, m)
	}
	return out, nil
}

func parseFloatAuth(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// initBearingDBFromCSV mirrors Python init_bearing_db_from_csv — seeds bearings.db from NTEKBearingDB.csv.
func initBearingDBFromCSV(configDir string) bearingInitResult {
	dbPath := filepath.Join(configDir, "bearings.db")
	csvPath := filepath.Join(configDir, "NTEKBearingDB.csv")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return bearingInitResult{msg: err.Error()}
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS bearings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT UNIQUE NOT NULL,
		BPFO REAL NOT NULL, BPFI REAL NOT NULL,
		BSF REAL NOT NULL, FTF REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return bearingInitResult{msg: err.Error()}
	}

	if _, err := os.Stat(csvPath); err != nil {
		// No seed CSV — treat as soft success so joinAdmin doesn't fail.
		return bearingInitResult{success: true, msg: "NTEKBearingDB.csv not found, table ready"}
	}

	f, err := os.Open(csvPath)
	if err != nil {
		return bearingInitResult{msg: err.Error()}
	}
	defer f.Close()

	rows, err := parseBearingCSV(f)
	if err != nil {
		return bearingInitResult{msg: err.Error()}
	}

	stmt, err := db.Prepare(`INSERT INTO bearings (Name, BPFO, BPFI, BSF, FTF) VALUES (?,?,?,?,?)`)
	if err != nil {
		return bearingInitResult{msg: err.Error()}
	}
	defer stmt.Close()

	res := bearingInitResult{success: true}
	for _, r := range rows {
		if _, err := stmt.Exec(r["Name"], parseFloatAuth(r["BPFO"]), parseFloatAuth(r["BPFI"]), parseFloatAuth(r["BSF"]), parseFloatAuth(r["FTF"])); err != nil {
			res.skipped++
			continue
		}
		res.inserted++
	}
	res.msg = fmt.Sprintf("inserted=%d skipped=%d", res.inserted, res.skipped)
	return res
}
