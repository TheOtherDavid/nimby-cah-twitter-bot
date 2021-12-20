package main

import (
	"encoding/csv"
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
	cardNumber, _ := strconv.Atoi(getLastCardNumberFromFile("last_card.csv"))
	cardMessage := cards[cardNumber]
	fullMessage := message + cardMessage
	fmt.Println(fullMessage)
	//Send message to Twitter API
	SendTweet(fullMessage)

	//Update env variable
	newCardNumber := cardNumber + 1
	writeLastCardNumberToFile("last_card.csv", strconv.Itoa(newCardNumber))

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

func getLastCardNumberFromFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()

	if err != nil {
		fmt.Println(err)
	}

	lastCard := records[0][0]
	return lastCard
}

func writeLastCardNumberToFile(filename string, lastCard string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	lastCardArray := []string{lastCard}

	err = writer.Write(lastCardArray)

	if err != nil {
		fmt.Println(err)
	}
}
