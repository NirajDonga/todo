package user

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"CreatedAt"`
}

type RegisterUser struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogInUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
