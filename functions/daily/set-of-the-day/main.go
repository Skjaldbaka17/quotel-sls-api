package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

func (requestHandler *RequestHandler) insertEnglishQOD(today string) error {
	//------------------- Start ENGLISH QOD -------------------//

	var err error
	var quote structs.QuoteDBModel
	for i := 0; i < 10; i++ {
		//Note we are sampling from table "topicsview" but not quotes for more chance of a good quote
		err = requestHandler.Db.Raw("select * from topicsview tablesample system(0.1)").Limit(1).First(&quote).Error
		if err != nil {
			log.Fatalf("Got error when getting random english quote: %s", err)
		}
		if !quote.IsIcelandic {
			break
		}
		quote = structs.QuoteDBModel{}
	}
	createQOD := quote.ConvertToQODDBModel(today)
	err = requestHandler.Db.Table("qods").Create(&createQOD).Error
	if err != nil {
		log.Fatalf("Got error when creating english QOD: %s", err)
	}

	//------------------- ENGLISH QOD DONE -------------------//
	return nil
}

func (requestHandler *RequestHandler) insertIcelandicQOD(today string) error {
	//------------------- Start ICELANDIC QOD -------------------//
	var err error
	var quote structs.QuoteDBModel
	quote = structs.QuoteDBModel{
		IsIcelandic: true,
	}
	err = requestHandler.Db.Table("quote").Order("random()").Limit(1).First(quote).Error
	if err != nil {
		log.Fatalf("Got error when getting random icelandic quote: %s", err)
	}
	createQOD := quote.ConvertToQODDBModel(today)
	err = requestHandler.Db.Table("qods").Create(&createQOD).Error
	if err != nil {
		log.Fatalf("Got error when creating icelandic QOD: %s", err)
	}

	//------------------- ICELANDIC QOD DONE -------------------//
	return nil
}

func (requestHandler *RequestHandler) insertTopicsQOD(today string) {
	var topics []structs.TopicDBModel
	err := requestHandler.Db.Table("topics").Find(&topics).Error
	if err != nil {
		log.Fatalf("Got error when getting all topics: %s", err)
	}

	for _, topic := range topics {
		quote := structs.QuoteDBModel{
			TopicId:   topic.Id,
			TopicName: topic.Name,
		}
		err = requestHandler.Db.Table("topicsview").Order("random()").Limit(1).First(&quote).Error
		if err != nil {
			log.Fatalf("Got error when getting a random quote for topic %s: %s", topic.Name, err)
		}

		err = requestHandler.Db.Table("qods").Create(&quote).Error
		if err != nil {
			log.Fatalf("Got error when creating QOD for topic %s: %s", topic.Name, err)
		}
	}

}

func (requestHandler *RequestHandler) insertEnglishAOD(today string, isIcelandic bool) {
	//------------------- Start AOD -------------------//
	var err error
	var author structs.AuthorDBModel
	for i := 0; i < 10; i++ {
		//Note we are sampling from table "topicsview" but not quotes for more chance of a good quote
		err = requestHandler.Db.Raw("select * from authors tablesample system(0.1)").Limit(1).First(&author).Error
		if err != nil {
			log.Fatalf("Got error when getting random english author: %s", err)
		}
		if author.NrOfEnglishQuotes > 0 {
			break
		}
		author = structs.AuthorDBModel{}
	}
	createAOD := author.ConvertToAODDBModel(today)
	err = requestHandler.Db.Table("aods").Create(&createAOD).Error
	if err != nil {
		log.Fatalf("Got error when creating english AOD: %s", err)
	}

	//------------------- AOD DONE -------------------//
}

func (requestHandler *RequestHandler) insertIcelandicAOD(today string, isIcelandic bool) {
	//------------------- Start AOD -------------------//
	var err error
	var author structs.AuthorDBModel
	err = requestHandler.Db.Table("authors").Order("random()").Where("nr_of_icelandic_quotes > 0").Limit(1).First(&author).Error
	if err != nil {
		log.Fatalf("Got error when getting random english author: %s", err)
	}

	createAOD := author.ConvertToAODDBModel(today)
	err = requestHandler.Db.Table("aods").Create(&createAOD).Error
	if err != nil {
		log.Fatalf("Got error when creating english AOD: %s", err)
	}

	//------------------- AOD DONE -------------------//
}

//Creates and inserts the AOD, AODICE, QOD and QODICE for today
func (requestHandler *RequestHandler) handler(request events.APIGatewayProxyRequest) {
	//Initialize DB if requestHandler.Db = nil
	if errResponse := requestHandler.InitializeDB(); errResponse != (structs.ErrorResponse{}) {
		log.Fatalf("Could not connect to DB when creating AOD/AODICE/QOD/QODICE")
	}
	year, month, day := time.Now().Date()
	today := fmt.Sprintf("%d-%d-%d", year, month, day)

	go requestHandler.insertEnglishQOD(today)
	go requestHandler.insertIcelandicQOD(today)
	go requestHandler.insertEnglishQOD(today)
	go requestHandler.insertIcelandicQOD(today)
	go requestHandler.insertTopicsQOD(today)
}

func main() {
	lambda.Start(theReqHandler.handler)
}
