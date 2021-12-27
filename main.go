package main

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"net/http"
	"os"
	"strconv"
	"strings"
)

func handleRequest() (string, error) {
	GenerateNIMBYTweet()
	return "Success", nil
}

func main() {
	fmt.Println("Executing program on env " + os.Getenv("ENV"))
	if os.Getenv("ENV") == "aws" {
		lambda.Start(handleRequest)
	} else {
		GenerateNIMBYTweet()
	}
}

func GenerateNIMBYTweet() {
	//Construct string
	message := "We aren't NIMBYs, we're just concerned about "
	//Import list of cards
	cards := getCardsList("cards.txt")
	//Get the right card
	filePath := os.Getenv("FILE_PATH")
	//cardNumber, _ := strconv.Atoi(getLastCardNumberFromFile("last_card.csv"))
	cardNumber, _ := strconv.Atoi(getLastCardNumberFromS3(filePath + "last_card.csv"))
	cardMessage := cards[cardNumber]
	fullMessage := message + cardMessage
	fmt.Println(fullMessage)
	//Send message to Twitter API
	SendTweet(fullMessage)

	//Update env variable
	newCardNumber := cardNumber + 1

	writeLastCardNumberToS3(filePath+"last_card.csv", strconv.Itoa(newCardNumber))
	fmt.Println("File successfully updated")

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
		os.Exit(1)
	}

	lastCard := records[0][0]
	return lastCard
}

func writeLastCardNumberToFile(filename string, lastCard string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	lastCardArray := []string{lastCard}

	err = writer.Write(lastCardArray)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getLastCardNumberFromS3(filename string) string {
	bucket := os.Getenv("AWS_BUCKET")
	file, err := DownloadFromS3Bucket(filename, bucket)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	records, err := csv.NewReader(&file).ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lastCard := records[0][0]
	fmt.Println("Retrieved", file.Name(), "from AWS")

	return lastCard
}

func DownloadFromS3Bucket(filename string, bucket string) (os.File, error) {
	fmt.Println("Creating file", filename)

	file, err := os.Create(filename)
	fmt.Println("File", file.Name(), "created")

	if err != nil {
		fmt.Println(err)
	}

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		})
	if err != nil {
		fmt.Println(err)
		return os.File{}, err
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return *file, nil
}

func writeLastCardNumberToS3(filename string, lastCard string) {
	bucket := os.Getenv("AWS_BUCKET")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	lastCardArray := []string{lastCard}

	err = writer.Write(lastCardArray)
	writer.Flush()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	UploadToS3Bucket(filename, bucket)
	defer file.Close()

}

func UploadToS3Bucket(filename string, bucket string) (os.File, error) {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	file.Read(fileBuffer)

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(filename),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Uploaded", file.Name())
	return *file, nil
}
