package app

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) RegisterUser(
	fullName string,
	age int,
	gender string,
	email string,
	password string,
) (string, error) {
	return "", nil
}

func (s *AuthService) LoginUser(email string, password string) (string, error) {
	return "", nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (string, error) {
	return "", nil
}