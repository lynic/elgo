package elgo

import (
	"context"
	"fmt"
	"os"
)

const InternalCall = "internalCall"

type AuthContext struct {
	UserID   string
	Username string
	Roles    []string
}

type AuthResponse struct {
	AuthUserPass *AuthContext
}

func (a *AuthContext) IsAdmin() bool {
	return StrIn(a.Roles, "Admin")
}

type AuthHandler struct {
	Policy map[string][]string
}

func (r *AuthHandler) Init() error {
	// policyPath := os.Getenv()
	policy, _ := r.LoadPolicy()
	// if err != nil {
	// 	return err
	// }
	r.Policy = policy
	return nil
}

func (r *AuthHandler) LoadPolicy() (map[string][]string, error) {
	rbacFilePath := os.Getenv("RBACFile")
	if rbacFilePath == "" {
		rbacFilePath = "./rbac.json"
	}
	// var policy *map[string][]string
	policy := make(map[string][]string)
	// ppolicy := &policy
	err := JsonFromFile(rbacFilePath, &policy)
	if err != nil {
		fmt.Printf("Failed to load policy from %s:%s\n", rbacFilePath, err.Error())
		return policy, err
	}
	return policy, nil
}

func (r *AuthHandler) AuthCtx(ctx context.Context) (*AuthContext, error) {
	authCtx, exists := ctx.Value(AuthContextKey).(*AuthContext)
	if exists == false {
		return nil, fmt.Errorf("no authentication context")
	}
	return authCtx, nil
}

func (r *AuthHandler) CheckRule(ctx context.Context, method string) bool {
	if _, exists := ctx.Value(InternalCall).(bool); exists {
		return true
	}
	authCtx, err := r.AuthCtx(ctx)
	if err != nil {
		// no authentication wrapper
		return true
	}
	if authCtx.IsAdmin() {
		return true
	}
	// no rbac rules for method
	// if _, exists := r.Policy[method]; exists == false {
	// 	if authCtx.IsAdmin() {
	// 		return true
	// 	}
	// 	return false
	// }
	// check rbac rules for method
	for _, role := range authCtx.Roles {
		if StrIn(r.Policy[method], role) {
			return true
		}
	}
	return false
}

func (r *AuthHandler) IsSelf(ctx context.Context, username string) bool {
	authCtx, err := r.AuthCtx(ctx)
	if err != nil {
		// no auth wrapper
		return true
	}
	if authCtx.IsAdmin() {
		return true
	}
	if authCtx.Username == username {
		return true
	}
	return false
}

func (r *AuthHandler) IsAdmin(ctx context.Context) bool {
	authCtx, err := r.AuthCtx(ctx)
	if err != nil {
		// no auth wrapper
		return true
	}
	return authCtx.IsAdmin()
}
