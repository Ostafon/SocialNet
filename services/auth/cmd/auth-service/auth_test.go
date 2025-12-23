package main

import (
	"fmt"
	"os"
	"socialnet/services/auth/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "socialnet/services/auth/gen"
	"socialnet/services/auth/internal/model"
	"socialnet/services/auth/internal/repos"
)

var testSvc *service.AuthService
var testDB *gorm.DB

// =========================================================
//
//	TestMain
//
// =========================================================
func TestMain(m *testing.M) {

	// Используем переменную окружения, если есть
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://ostafon:0000@localhost:5432/test?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("❌ cannot connect to TEST PostgreSQL: %v", err))
	}
	testDB = db

	fmt.Println("🔵 Running AUTH migrations on TEST DB...")

	// --- удаляем только таблицы AUTH сервиса ---
	_ = db.Migrator().DropTable(
		&model.User{},
		&model.RefreshToken{},
		&model.PasswordReset{},
	)

	// --- выполняем миграции ---
	if err := db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.PasswordReset{},
	); err != nil {
		panic(fmt.Sprintf("❌ migration error: %v", err))
	}

	// создаём repo/service
	repo := repos.NewAuthRepo(db)
	testSvc = service.NewAuthService(repo)

	// запускаем все тесты
	code := m.Run()
	os.Exit(code)
}

// =========================================================
//           Helper — очищаем таблицы перед тестом
// =========================================================

func resetTables(t *testing.T) {
	err := testDB.Exec(`
		TRUNCATE users, refresh_tokens, password_resets 
		RESTART IDENTITY CASCADE;
	`).Error
	if err != nil {
		t.Fatalf("❌ reset error: %v", err)
	}
}

// =========================================================
//                       TESTS
// =========================================================

// --------------------- Register OK -----------------------

func TestRegister_Success(t *testing.T) {
	resetTables(t)

	req := &pb.RegisterRequest{
		Email:    "reg_ok@example.com",
		Password: "StrongPass123!",
		Username: "reguser",
	}

	access, refresh, err := testSvc.Register(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
}

// ------------- Register Duplicate Email -------------------

func TestRegister_DuplicateEmail(t *testing.T) {
	resetTables(t)

	_, _, err := testSvc.Register(&pb.RegisterRequest{
		Email:    "dup@example.com",
		Password: "12345678Aa!",
		Username: "dupuser",
	})
	assert.NoError(t, err)

	_, _, err2 := testSvc.Register(&pb.RegisterRequest{
		Email:    "dup@example.com",
		Password: "12345678Aa!",
		Username: "dupuser",
	})

	assert.Error(t, err2)
}

// ----------------------- Login OK -------------------------

func TestLogin_Success(t *testing.T) {
	resetTables(t)

	email := "login_ok@example.com"
	pass := "Pass123456!"

	_, _, err := testSvc.Register(&pb.RegisterRequest{
		Email:    email,
		Password: pass,
		Username: "testlogin",
	})
	assert.NoError(t, err)

	access, refresh, err := testSvc.Login(&pb.LoginRequest{
		Email:    email,
		Password: pass,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
}

// ------------- Login Wrong Password -----------------------

func TestLogin_WrongPassword(t *testing.T) {
	resetTables(t)

	email := "wrongp@example.com"

	_, _, err := testSvc.Register(&pb.RegisterRequest{
		Email:    email,
		Password: "CorrectPass123!",
		Username: "wrongp",
	})
	assert.NoError(t, err)

	_, _, err = testSvc.Login(&pb.LoginRequest{
		Email:    email,
		Password: "WrongP",
	})

	assert.Error(t, err)
}

// -------------------- Refresh Token OK ---------------------

func TestRefreshToken_Success(t *testing.T) {
	resetTables(t)

	email := "ref@example.com"
	pass := "Pass123456!"

	_, _, err := testSvc.Register(&pb.RegisterRequest{
		Email:    email,
		Password: pass,
		Username: "refuser",
	})
	assert.NoError(t, err)

	_, refresh, err := testSvc.Login(&pb.LoginRequest{
		Email:    email,
		Password: pass,
	})
	assert.NoError(t, err)

	newAccess, err := testSvc.RefreshToken(&pb.RefreshRequest{
		RefreshToken: refresh,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccess)
}

// -------------------- Update Password OK --------------------

func TestUpdatePassword_Success(t *testing.T) {
	resetTables(t)

	email := "upd@example.com"
	old := "OldPass123!"
	newp := "NewPass321!"

	_, _, err := testSvc.Register(&pb.RegisterRequest{
		Email:    email,
		Password: old,
		Username: "upduser",
	})
	assert.NoError(t, err)

	user, err := testSvc.Repo.GetUserByEmail(email)
	assert.NoError(t, err)

	_, err = testSvc.UpdatePassword(&pb.UpdatePasswordRequest{
		Id:              fmt.Sprint(user.ID),
		CurrentPassword: old,
		NewPassword:     newp,
	})
	assert.NoError(t, err)

	_, _, err = testSvc.Login(&pb.LoginRequest{
		Email:    email,
		Password: newp,
	})
	assert.NoError(t, err)
}

// -------------------- GetProfile OK -------------------------

func TestGetProfile_Success(t *testing.T) {
	resetTables(t)

	email := "prof@example.com"
	pass := "Pass123456!"

	_, _, err := testSvc.Register(&pb.RegisterRequest{
		Email:    email,
		Password: pass,
		Username: "profuser",
	})
	assert.NoError(t, err)

	user, err := testSvc.Repo.GetUserByEmail(email)
	assert.NoError(t, err)

	profile, err := testSvc.GetProfile(user.ID)
	assert.NoError(t, err)

	assert.Equal(t, email, profile.Email)
	assert.Equal(t, "profuser", profile.Username)
}
