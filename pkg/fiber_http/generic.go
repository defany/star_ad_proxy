package httpwrap

import (
	"github.com/defany/goblin/errfmt"
	"github.com/gofiber/fiber/v3/client"
)

type Response[B any] struct {
	*client.Response
	Body B
}

func Post[B any](c *client.Client, url string, cfg ...client.Config) (Response[B], error) {
	rawResp, err := c.Post(url, cfg...)
	if err != nil {
		return Response[B]{}, errfmt.WithSource(err)
	}

	var body B

	if err := rawResp.JSON(&body); err != nil {
		return Response[B]{}, errfmt.WithSource(err)
	}

	return Response[B]{
		Response: rawResp,
		Body:     body,
	}, nil
}
