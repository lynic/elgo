package elgo

import (
	"context"
	"fmt"
	"os"

	"github.com/machinebox/graphql"
)

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
	err := JsonFromFile(rbacFilePath, policy)
	if err != nil {
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
	authCtx, err := r.AuthCtx(ctx)
	if err != nil {
		// no authentication wrapper
		return true
	}
	// no rbac rules for method
	if _, exists := r.Policy[method]; exists == false {
		if authCtx.IsAdmin() {
			return true
		}
		return false
	}
	// check rbac rules for method
	for _, role := range authCtx.Roles {
		if StrIn(r.Policy[method], role) {
			return true
		}
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

func AuthRequest(authorization string) (*AuthContext, error) {
	authURL := os.Getenv("AUTHURL")
	if authURL == "" {
		return nil, fmt.Errorf("No AUTHURL")
	}
	username, password, err := DecodeAuthHeader(authorization)
	if err != nil {
		return nil, err
	}
	fmt.Println(username, password)
	godUser := os.Getenv("GODUSER")
	if godUser != "" && godUser == username {
		ret := &AuthContext{
			UserID:   godUser,
			Username: godUser,
			Roles:    []string{"Admin"},
		}
		return ret, nil
	}
	resp := &AuthResponse{}
	// authUser := &AuthUser{}
	cli := graphql.NewClient(authURL)
	query := fmt.Sprintf(`
		query {
			authUserPass(username: "%s", password: "%s") {
				userID
				username
				roles
			}
		}`, username, password)
	req := graphql.NewRequest(query)
	// authorization := fmt.Sprintf("Basic %s", GenAuthorization(username, password))
	req.Header.Set("Authorization", authorization)
	ctx := context.Background()
	err = cli.Run(ctx, req, resp)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", resp.AuthUserPass)
	return resp.AuthUserPass, nil
}
