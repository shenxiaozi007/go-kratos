package biz

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// mockUserRepo 内存实现的 UserRepo，用于单测。
type mockUserRepo struct {
	byUsername map[string]*User
	byID      map[uint64]*User
	nextID    uint64
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		byUsername: make(map[string]*User),
		byID:      make(map[uint64]*User),
		nextID:    1,
	}
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	u, ok := m.byUsername[username]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *u
	return &cp, nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uint64) (*User, error) {
	u, ok := m.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *u
	return &cp, nil
}

func (m *mockUserRepo) Create(ctx context.Context, u *User) (*User, error) {
	u.ID = m.nextID
	m.nextID++
	m.byUsername[u.Username] = u
	m.byID[u.ID] = u
	cp := *u
	return &cp, nil
}

func (m *mockUserRepo) Count(ctx context.Context) (int64, error) {
	return int64(len(m.byUsername)), nil
}

func TestAuthUsecase_Register(t *testing.T) {
	ctx := context.Background()

	t.Run("first_user_gets_admin_role", func(t *testing.T) {
		repo := newMockUserRepo()
		uc := NewAuthUsecase(repo)
		hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)

		u, err := uc.Register(ctx, "admin1", string(hash))
		if err != nil {
			t.Fatalf("Register: %v", err)
		}
		if u.Role != "admin" {
			t.Errorf("first user role = %s, want admin", u.Role)
		}
		if u.Username != "admin1" {
			t.Errorf("Username = %s, want admin1", u.Username)
		}
	})

	t.Run("second_user_gets_operator_role", func(t *testing.T) {
		repo := newMockUserRepo()
		uc := NewAuthUsecase(repo)
		hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
		_, _ = uc.Register(ctx, "first", string(hash))

		u, err := uc.Register(ctx, "second", string(hash))
		if err != nil {
			t.Fatalf("Register: %v", err)
		}
		if u.Role != "operator" {
			t.Errorf("second user role = %s, want operator", u.Role)
		}
	})

	t.Run("duplicate_username_returns_ErrUserExists", func(t *testing.T) {
		repo := newMockUserRepo()
		uc := NewAuthUsecase(repo)
		hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
		_, _ = uc.Register(ctx, "dup", string(hash))

		_, err := uc.Register(ctx, "dup", string(hash))
		if err == nil {
			t.Fatal("expected error for duplicate username")
		}
		if !errors.Is(err, ErrUserExists) {
			t.Errorf("got err %v, want ErrUserExists", err)
		}
	})
}

func TestAuthUsecase_Login(t *testing.T) {
	ctx := context.Background()
	plain := "mypassword"
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("success", func(t *testing.T) {
		repo := newMockUserRepo()
		_, _ = repo.Create(ctx, &User{Username: "alice", PasswordHash: string(hash), Role: "admin"})
		uc := NewAuthUsecase(repo)

		u, err := uc.Login(ctx, "alice", plain)
		if err != nil {
			t.Fatalf("Login: %v", err)
		}
		if u.Username != "alice" || u.Role != "admin" {
			t.Errorf("got Username=%s Role=%s", u.Username, u.Role)
		}
	})

	t.Run("user_not_found_returns_ErrAuthUserNotFound", func(t *testing.T) {
		repo := newMockUserRepo()
		uc := NewAuthUsecase(repo)

		_, err := uc.Login(ctx, "nobody", plain)
		if err == nil {
			t.Fatal("expected error for unknown user")
		}
		if !errors.Is(err, ErrAuthUserNotFound) {
			t.Errorf("got err %v, want ErrAuthUserNotFound", err)
		}
	})

	t.Run("wrong_password_returns_ErrBadPassword", func(t *testing.T) {
		repo := newMockUserRepo()
		_, _ = repo.Create(ctx, &User{Username: "alice", PasswordHash: string(hash), Role: "admin"})
		uc := NewAuthUsecase(repo)

		_, err := uc.Login(ctx, "alice", "wrongpassword")
		if err == nil {
			t.Fatal("expected error for wrong password")
		}
		if !errors.Is(err, ErrBadPassword) {
			t.Errorf("got err %v, want ErrBadPassword", err)
		}
	})
}

func TestAuthUsecase_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		repo := newMockUserRepo()
		created, _ := repo.Create(ctx, &User{Username: "getbyid", PasswordHash: "hash", Role: "operator"})
		uc := NewAuthUsecase(repo)

		u, err := uc.GetByID(ctx, created.ID)
		if err != nil {
			t.Fatalf("GetByID: %v", err)
		}
		if u.ID != created.ID || u.Username != "getbyid" {
			t.Errorf("got ID=%d Username=%s", u.ID, u.Username)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		repo := newMockUserRepo()
		uc := NewAuthUsecase(repo)

		_, err := uc.GetByID(ctx, 999)
		if err == nil {
			t.Fatal("expected error for missing id")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("got err %v", err)
		}
	})
}
