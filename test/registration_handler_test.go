package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"nazartaraniuk/alertsProject/internal/adapter/handler"
	"nazartaraniuk/alertsProject/internal/domain"
	"nazartaraniuk/alertsProject/internal/repository"
)

func makeRequest(t *testing.T, e *echo.Echo, h echo.HandlerFunc, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	return rec
}

func TestRegistrationHandler_Success(t *testing.T) {
	e := echo.New()
	mockSvc := &mockUserService{
		RegisterUserFunc: func(req domain.RegisterReq) (int64, error) {
			if req.Email != "a@b.com" || req.Username != "nazar" || req.Password != "pass" {
				t.Fatalf("unexpected req: %+v", req)
			}
			return 42, nil
		},
	}
	h := handler.RegistrationHandler(mockServiceAdapter{mockSvc})

	rec := makeRequest(t, e, h, http.MethodPost, "/register", domain.RegisterReq{
		Email:    "a@b.com",
		Username: "nazar",
		Password: "pass",
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var resp domain.UserCreatedResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("bad json: %v, body=%s", err, rec.Body.String())
	}
	if resp.ID != 42 || resp.Email != "a@b.com" || resp.Username != "nazar" {
		t.Fatalf("unexpected resp: %+v", resp)
	}
}

func TestRegistrationHandler_BadRequest_InvalidJSON(t *testing.T) {
	e := echo.New()
	mockSvc := &mockUserService{
		RegisterUserFunc: func(req domain.RegisterReq) (int64, error) { return 0, nil },
	}
	h := handler.RegistrationHandler(mockServiceAdapter{mockSvc})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("{bad json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = h(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestRegistrationHandler_BadRequest_EmptyFields(t *testing.T) {
	e := echo.New()
	mockSvc := &mockUserService{
		RegisterUserFunc: func(req domain.RegisterReq) (int64, error) { return 0, nil },
	}
	h := handler.RegistrationHandler(mockServiceAdapter{mockSvc})

	rec := makeRequest(t, e, h, http.MethodPost, "/register", domain.RegisterReq{
		Email:    "",
		Username: "nazar",
		Password: "",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestRegistrationHandler_Conflict_EmailExists(t *testing.T) {
	e := echo.New()
	mockSvc := &mockUserService{
		RegisterUserFunc: func(req domain.RegisterReq) (int64, error) {
			return 0, repository.ErrEmailAlreadyExists
		},
	}
	h := handler.RegistrationHandler(mockServiceAdapter{mockSvc})

	rec := makeRequest(t, e, h, http.MethodPost, "/register", domain.RegisterReq{
		Email:    "a@b.com",
		Username: "nazar",
		Password: "pass",
	})
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusConflict, rec.Body.String())
	}
	if got := rec.Body.String(); got == "" {
		t.Fatalf("expected error message in body")
	}
}

func TestRegistrationHandler_InternalError(t *testing.T) {
	e := echo.New()
	mockSvc := &mockUserService{
		RegisterUserFunc: func(req domain.RegisterReq) (int64, error) {
			return 0, repository.ErrCannotCreateUser
		},
	}
	h := handler.RegistrationHandler(mockServiceAdapter{mockSvc})

	rec := makeRequest(t, e, h, http.MethodPost, "/register", domain.RegisterReq{
		Email:    "a@b.com",
		Username: "nazar",
		Password: "pass",
	})
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

type userServiceInterface interface {
	RegisterUser(domain.RegisterReq) (int64, error)
}

type mockServiceAdapter struct{ m userServiceInterface }

func (a mockServiceAdapter) LoginUser(req domain.LoginReq) (*domain.User, error) {
	panic("unimplemented")
}

func (a mockServiceAdapter) RegisterUser(req domain.RegisterReq) (int64, error) {
	return a.m.RegisterUser(req)
}
