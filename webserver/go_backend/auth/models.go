package auth

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Lang     string `json:"lang"`
}

type SignupUserRequest struct {
	Username string `json:"username"`
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type SignupAdminRequest struct {
	DevType   interface{} `json:"devType"`
	Username  string `json:"username"`
	Account   string `json:"account"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	AdminPass string `json:"adminPass"`
	Lang      string `json:"lang"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	NewPass  string `json:"newPass"`
	Role     string `json:"role"`
	API      string `json:"api"`
}

type ResetPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type VerifyAdminRequest struct {
	Password string `json:"password" binding:"required"`
}
