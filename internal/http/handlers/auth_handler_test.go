package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)


func TestAuthHandler(t *testing.T) {

	t.Run("should return token", func(t *testing.T) {
		handler := NewAuthHandler(&mockAuthSvc{})

		payload := registerDTO{
			FullName: "Juan Diego Bermudez",
			Age: 21,
			Gender: "Masculino",
			Email: "d@gmail.com",
			Password: "12345678",
		}
		bodyContent, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(bodyContent))
		if err != nil{
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Post("/v1/auth/register", handler.registerUser)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated{
			t.Errorf("Expected status code %d, got %d %s", http.StatusCreated, rr.Code, rr.Body)
		}
	})
}

type mockAuthSvc struct{
}

func (s *mockAuthSvc) RegisterUser(
	fullName string,
	age int,
	gender string,
	email string,
	password string,
) (string, error){
	return "12345678", nil
}

func (s *mockAuthSvc)LoginUser(email string, password string) (string, error){
	return "12345678", nil
}

func (s *mockAuthSvc)RefreshAccessToken(refreshToken string) (string, error){
	return "12345678", nil
}