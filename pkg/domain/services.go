package domain

type AuthSvc interface {
	RegisterUser(
		fullName string,
		age int,
		gender string,
		email string,
		password string,
	) (string, error)
	LoginUser(email string, password string) (string, error)
	RefreshAccessToken(refreshToken string) (string, error)
}

type UserSvc interface {
}