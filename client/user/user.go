package user

type User struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	History History
}

type History struct {
	BrandName []string
}
