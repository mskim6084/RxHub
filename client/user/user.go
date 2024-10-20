package user

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

var LoggedInUsers []User
var CreatedUsers = make(map[int]User)
var UserCount = 0

type History struct {
	BrandName []string
}
