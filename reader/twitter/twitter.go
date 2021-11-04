package twitter

import (
	"fmt"
	"net/url"
	"path"
	"strconv"

	"git.icyphox.sh/forlater/navani/reader"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Creds struct {
	AccessToken       string
	AccessTokenSecret string
	APIKey            string
	APISecret         string
}

func GetClient(creds *Creds) *twitter.Client {
	config := oauth1.NewConfig(creds.APIKey, creds.APISecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return client
}

func parseSnowflake(url string) (int64, error) {
	strid := path.Base(url)
	id, err := strconv.ParseInt(strid, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("snowflake: %w\n", err)
	}

	return id, nil
}

func GetTweet(client *twitter.Client, url string) (*twitter.Tweet, error) {
	id, err := parseSnowflake(url)
	if err != nil {
		return nil, err
	}

	tweet, _, err := client.Statuses.Show(id, nil)
	return tweet, err
}

func MakeTweetArticle(tweet *twitter.Tweet, u *url.URL) *reader.Article {
	article := reader.Article{}
	article.Content = tweet.Text
	article.Title = fmt.Sprintf("Tweet: %s @%s", tweet.User.Name, tweet.User.ScreenName)
	article.Byline = fmt.Sprintf("@%s tweeted on %s", tweet.User.ScreenName, tweet.CreatedAt)
	article.URL = u
	return &article
}
