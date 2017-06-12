package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-resty/resty"
	"github.com/urfave/cli"
)

type Body struct {
	Response struct {
		Posts []struct {
			Type   string `json:"type"`
			Photos []struct {
				OriginalSize struct {
					Url string `json:"url"`
				} `json:"original_size"`
			} `json:"photos"`
		} `json:"posts"`
	} `json:"response"`
}

func main() {
	app := cli.NewApp()
	app.Name = "random tumblr"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "api-key",
		},
		cli.StringFlag{
			Name: "blog",
		},
		cli.StringFlag{
			Name: "tag",
		},
	}

	app.Action = func(c *cli.Context) error {
		path := fmt.Sprintf("https://api.tumblr.com/v2/blog/%s/posts/photo", c.String("blog"))

		resp, err := resty.R().
			SetQueryParams(map[string]string{
				"api_key": c.String("api-key"),
				"tag":     c.String("tag"),
				"limit":   "50",
			}).
			Get(path)

		if err != nil {
			return err
		}
		var body = Body{}
		err = json.Unmarshal(resp.Body(), &body)
		if err != nil {
			return err
		}
		posts := body.Response.Posts
		postLen := len(posts)
		rand.Seed(time.Now().UnixNano())
		idx := rand.Intn(postLen - 1)
		fmt.Print(posts[idx].Photos[0].OriginalSize.Url)
		return nil
	}

	app.Run(os.Args)
}
