package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"nazartaraniuk/alertsProject/internal/adapter/handler"
	"nazartaraniuk/alertsProject/internal/config"
	"nazartaraniuk/alertsProject/internal/domain"
)

func doReq(t *testing.T, e *echo.Echo, h echo.HandlerFunc, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(http.MethodPost, "/login", &buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	return rec
}

func mustHash(t *testing.T, pwd string, cost int) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	return string(h)
}


func TestLogin_BadJSON(t *testing.T) {
	e := echo.New()
	secret := "testsecret"

	cfg := config.Config{
		Server: config.ServerConfig{
			JWTSecret: secret,
		},
	}

	svc := &mockUserService{
		LoginUserFunc: func(domain.LoginReq) (*domain.User, error) {
			return &domain.User{}, nil
		},
	}

	h := handler.LoginHandler(cfg, svc)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{bad"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	_ = h(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestLogin_NoRows_Unauthorized(t *testing.T) {
	e := echo.New()
	secret := "testsecret"
	cfg := config.Config{Server: config.ServerConfig{JWTSecret: secret}}

	svc := &mockUserService{
		LoginUserFunc: func(r domain.LoginReq) (*domain.User, error) {
			return &domain.User{}, sql.ErrNoRows
		},
	}

	h := handler.LoginHandler(cfg, svc)
	rec := doReq(t, e, h, domain.LoginReq{Email: "a@b.com", Password: "pass"})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d; body=%s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}
}

func TestLogin_ServiceError_500(t *testing.T) {
	e := echo.New()
	secret := "testsecret"
	cfg := config.Config{Server: config.ServerConfig{JWTSecret: secret}}

	svc := &mockUserService{
		LoginUserFunc: func(r domain.LoginReq) (*domain.User, error) {
			return &domain.User{}, errors.New("db down")
		},
	}

	h := handler.LoginHandler(cfg, svc)
	rec := doReq(t, e, h, domain.LoginReq{Email: "a@b.com", Password: "pass"})
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestLogin_InvalidPassword_401(t *testing.T) {
	e := echo.New()
	secret := "testsecret"
	cfg := config.Config{Server: config.ServerConfig{JWTSecret: secret}}

	hashed := mustHash(t, "goodpass", bcrypt.DefaultCost)
	user := domain.User{
		ID:           7,
		Email:        "a@b.com",
		Username:     "nazar",
		PasswordHash: hashed,
		CreatedAt:    time.Now(),
	}

	svc := &mockUserService{
		LoginUserFunc: func(r domain.LoginReq) (*domain.User, error) {
			return &user, nil
		},
	}

	h := handler.LoginHandler(cfg, svc)
	rec := doReq(t, e, h, domain.LoginReq{Email: "a@b.com", Password: "wrong"})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d, body=%s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}
	if got := rec.Body.String(); got == "" {
		t.Fatalf("expected error body")
	}
}

func TestLogin_Success_200_JWT(t *testing.T) {
	e := echo.New()
	secret := "testsecret"
	cfg := config.Config{Server: config.ServerConfig{JWTSecret: secret}}

	hashed := mustHash(t, "pass", bcrypt.DefaultCost)
	user := domain.User{
		ID:           42,
		Email:        "a@b.com",
		Username:     "nazar",
		PasswordHash: hashed,
		CreatedAt:    time.Now(),
	}

	svc := &mockUserService{
		LoginUserFunc: func(r domain.LoginReq) (*domain.User, error) {
			if r.Email != "a@b.com" {
				t.Fatalf("unexpected email: %s", r.Email)
			}
			return &user, nil
		},
	}

	h := handler.LoginHandler(cfg, svc)
	rec := doReq(t, e, h, domain.LoginReq{Email: "a@b.com", Password: "pass"})

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var tok domain.TokenResp
	if err := json.Unmarshal(rec.Body.Bytes(), &tok); err != nil {
		t.Fatalf("bad json: %v, body=%s", err, rec.Body.String())
	}
	if tok.AccessToken == "" || tok.ExpiresAt == 0 {
		t.Fatalf("empty token response: %+v", tok)
	}

	parsed, err := jwt.Parse(tok.AccessToken, func(tk *jwt.Token) (any, error) {
		if tk.Method != jwt.SigningMethodHS256 {
			t.Fatalf("unexpected signing method: %v", tk.Method)
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatalf("invalid token: %v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("claims cast failed")
	}

	if sub, _ := claims["sub"].(string); sub != strconv.FormatInt(user.ID, 10) {
		t.Fatalf("unexpected sub: %v", claims["sub"])
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("exp missing")
	}
	exp := time.Unix(int64(expFloat), 0)
	if time.Until(exp) <= 10*time.Minute || time.Until(exp) > 16*time.Minute {
		t.Fatalf("unexpected exp: %v (until %v)", exp, time.Until(exp))
	}
}
