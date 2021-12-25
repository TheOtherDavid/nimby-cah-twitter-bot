package main

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/oauth1"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

//Struct to parse tweet
type Tweet struct {
	Id    int64
	IdStr string `json:"id_str"`
	User  User
	Text  string
}

//Struct to parse user
type User struct {
	Id     int64
	IdStr  string `json:"id_str"`
	Name   string
	Handle string `json:"screen_name"`
}

func CreateClient() *http.Client {
	//Verify environment variables
	if os.Getenv("CONSUMER_KEY") == "" {
		fmt.Println("Consumer Key is not specified")
		os.Exit(1)
	}
	if os.Getenv("CONSUMER_SECRET") == "" {
		fmt.Println("Consumer Secret is not specified")
		os.Exit(1)
	}
	if os.Getenv("TOKEN") == "" {
		fmt.Println("Token is not specified")
		os.Exit(1)
	}
	if os.Getenv("TOKEN_SECRET") == "" {
		fmt.Println("Token Secret is not specified")
		os.Exit(1)
	}
	//Create oauth client with consumer keys and access token
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TOKEN"), os.Getenv("TOKEN_SECRET"))

	return config.Client(oauth1.NoContext, token)
}

func SendTweet(tweet string) (*Tweet, error) {
	//Initialize tweet object to store response in
	var responseTweet Tweet
	//Add params
	params := url.Values{}
	params.Set("status", tweet)
	//Grab client and post
	client := CreateClient()
	resp, err := client.PostForm("https://api.twitter.com/1.1/statuses/update.json", params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	//Decode response and send out
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	err = json.Unmarshal(body, &responseTweet)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &responseTweet, nil
}
