package auth

import (
	"context"
	"testing"
)

func TestWithUser_FromContext(t *testing.T) {
	ctx := context.Background()
	u := &UserInfo{UserID: 1, Username: "test", Role: "admin"}

	t.Run("nil_user", func(t *testing.T) {
		got := WithUser(ctx, nil)
		if got == nil {
			t.Fatal("WithUser(ctx, nil) must not return nil context")
		}
		out := FromContext(got)
		if out != nil {
			t.Errorf("FromContext(WithUser(ctx,nil)) = %v, want nil", out)
		}
	})

	t.Run("roundtrip", func(t *testing.T) {
		got := WithUser(ctx, u)
		out := FromContext(got)
		if out == nil {
			t.Fatal("FromContext: got nil")
		}
		if out.UserID != u.UserID || out.Username != u.Username || out.Role != u.Role {
			t.Errorf("got UserID=%d Username=%s Role=%s", out.UserID, out.Username, out.Role)
		}
	})

	t.Run("empty_ctx_returns_nil", func(t *testing.T) {
		out := FromContext(ctx)
		if out != nil {
			t.Errorf("FromContext(background) = %v, want nil", out)
		}
	})
}
