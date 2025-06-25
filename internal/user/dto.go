package user

// createUserRequest represents the request body for creating a new user.
// @Description Request payload for creating a new user.
type CreateUserRequest struct {
	// Username for the new user. Must be alphanumeric and unique.
	// required: true
	Username string `json:"username" binding:"required,alphanum"`
	// Password for the new user. Minimum 6 characters.
	// required: true
	Password string `json:"password" binding:"required,min=6"`
	// Full name of the user.
	// required: true
	FullName string `json:"full_name" binding:"required"`
	// Email address of the user. Must be a valid email format and unique.
	// required: true
	Email string `json:"email" binding:"required,email"`
}

// LoginUserRequest represents the request body for user login.
// @Description Request payload for logging in a user.
type LoginUserRequest struct {
	// Username of the user
	// required: true
	Username string `json:"username" binding:"required,alphanum"`
	// Password of the user
	// required: true
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	// Access token for the user
	// required: true
	AccessToken string `json:"access_token"`
	// User information
	// required: true
	User struct {
		// ID of the user
		// required: true
		ID string `json:"id"`
		// Username of the user
		// required: true
		Username string `json:"username"`
		// Full name of the user
		// required: true
		FullName string `json:"full_name"`
		// Email of the user
		// required: true
		Email string `json:"email"`
	} `json:"user"`
}
