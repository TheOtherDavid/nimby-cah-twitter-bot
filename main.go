package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	//Construct string
	message := "We aren't NIMBYs, we're just concerned about "
	//Import list of cards
	cards := getCardsList("cards.txt")
	//Get the right card
	cardNumber, _ := strconv.Atoi(os.Getenv("CARD_NUMBER"))
	cardMessage := cards[cardNumber]
	fullMessage := message + cardMessage
	fmt.Println(fullMessage)
	//Send message to Twitter API
	SendTweet(fullMessage)

	//Update env variable
	newCardNumber := cardNumber + 1
	os.Setenv("CARD_NUMBER", strconv.Itoa(newCardNumber))
}

func getCardsList(fileName string) []string {
	fileBytes, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sliceData := strings.Split(string(fileBytes), "\r\n")

	return sliceData
}
