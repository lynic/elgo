package elgo

import (
	"context"

	"github.com/machinebox/graphql"
)

type GClient struct {
	Url      string
	AuthUser string
	AuthPass string
}

func (c *GClient) Request(ctx context.Context, query string, resp interface{}) error {
	client := graphql.NewClient(c.Url)
	req := graphql.NewRequest(query)
	if c.AuthPass != "" {
		req.Header.Add("Authorization", BasicAuth(c.AuthUser, c.AuthPass))
	}
	err := client.Run(ctx, req, resp)
	return err
}
