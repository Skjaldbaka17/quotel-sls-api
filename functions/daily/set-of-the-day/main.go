package main

import (
	"fmt"
	"log"
	"sync"
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
	var quote structs.QuoteDBModel
	//Delete row if QOD for today has already been set
	err := requestHandler.Db.Exec("delete from qods where date = ? and topic_id = 0 and not is_icelandic", today).Error

	if err != nil {
		log.Fatalf("Got error when checking for existing QOD: %s", err)
	}

	for i := 0; i < 10; i++ {
		quote = structs.QuoteDBModel{}
		//Note we are sampling from table "topicsview" but not quotes for more chance of a good quote
		err = requestHandler.Db.Raw("select * from topicsview tablesample system(0.1)").Limit(1).First(&quote).Error
		if err != nil {
			log.Fatalf("Got error when getting random english quote: %s", err)
		}
		if !quote.IsIcelandic {
			break
		}
	}

	createQOD := quote.ConvertToQODDBModel(today)
	createQOD.TopicId = 0
	createQOD.TopicName = ""
	createQOD.IsIcelandic = false
	err = requestHandler.Db.Table("qods").Create(&createQOD).Error
	if err != nil {
		log.Fatalf("Got error when creating english QOD: %s", err)
	}

	//------------------- ENGLISH QOD DONE -------------------//
	return nil
}

func (requestHandler *RequestHandler) insertIcelandicQOD(today string) error {
	//------------------- Start ICELANDIC QOD -------------------//
	var quote structs.QuoteDBModel
	//Delete row if QOD for today has already been set
	err := requestHandler.Db.Exec("delete from qods where date = ? and topic_id = 0 and is_icelandic", today).Error

	if err != nil {
		log.Fatalf("Got error when checking for existing QOD: %s", err)
	}

	quote = structs.QuoteDBModel{}
	err = requestHandler.Db.Table("quotes").Where("is_icelandic").Order("random()").First(&quote).Error
	if err != nil {
		log.Fatalf("Got error when getting random icelandic quote: %s", err)
	}

	createQOD := quote.ConvertToQODDBModel(today)
	createQOD.TopicId = 0
	createQOD.TopicName = ""
	createQOD.IsIcelandic = true
	err = requestHandler.Db.Table("qods").Create(&createQOD).Error
	if err != nil {
		log.Fatalf("Got error when creating icelandic QOD: %s", err)
	}

	//------------------- ICELANDIC QOD DONE -------------------//
	return nil
}

func (requestHandler *RequestHandler) insertTopicsQOD(today string) {
	var wg sync.WaitGroup
	var topics []structs.TopicDBModel
	err := requestHandler.Db.Table("topics").Find(&topics).Error
	if err != nil {
		log.Fatalf("Got error when getting all topics: %s", err)
	}
	err = requestHandler.Db.Exec("delete from qods where date = ? and topic_id > 0", today).Error
	if err != nil {
		log.Fatalf("Got error second when getting all topics: %s", err)
	}
	for i, topic := range topics {
		wg.Add(1)

		go func(topic structs.TopicDBModel) {
			defer wg.Done()
			quote := structs.QuoteDBModel{
				TopicId:   topic.Id,
				TopicName: topic.Name,
			}
			err = requestHandler.Db.Table("topicsview").Order("random()").Limit(1).First(&quote).Error
			if err != nil {
				log.Fatalf("Got error when getting a random quote for topic %s: %s", topic.Name, err)
			}

			if err != nil {
				log.Fatalf("Got error when delete topic QOD for topic %s: %s", topic.Name, err)
			}
			createQuote := quote.ConvertToQODDBModel(today)
			createQuote.TopicId = uint(topic.Id)
			createQuote.TopicName = topic.Name
			err = requestHandler.Db.Table("qods").Create(&createQuote).Error
			if err != nil {
				log.Fatalf("Got error when creating QOD for topic %s: %s", topic.Name, err)
			}
		}(topic)

		if i%40 == 0 && i != 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}

func (requestHandler *RequestHandler) insertEnglishAOD(today string) {
	//------------------- Start AOD -------------------//
	var author structs.AuthorDBModel
	//Delete if QOD for today has already been set
	err := requestHandler.Db.Exec("delete from aods where date = ? and not is_icelandic", today).Error

	if err != nil {
		log.Fatalf("Got error when checking for existing AOD: %s", err)
	}

	for i := 0; i < 10; i++ {
		//Note we are sampling from table "topicsview" but not quotes for more chance of a good quote
		err = requestHandler.Db.Raw("select * from authors tablesample system(0.1)").Limit(1).Find(&author).Error
		if err != nil {
			log.Fatalf("Got error when getting random english author: %s", err)
		}
		if author != (structs.AuthorDBModel{}) && author.NrOfEnglishQuotes > 0 {
			break
		}
		author = structs.AuthorDBModel{}
	}
	createAOD := author.ConvertToAODDBModel(today, false)
	err = requestHandler.Db.Table("aods").Create(&createAOD).Error
	if err != nil {
		log.Fatalf("Got error when creating english AOD: %s", err)
	}

	//------------------- AOD DONE -------------------//
}

func (requestHandler *RequestHandler) insertIcelandicAOD(today string) {
	//------------------- Start AOD -------------------//
	var author structs.AuthorDBModel
	err := requestHandler.Db.Exec("delete from aods where date = ? and is_icelandic", today).Error

	if err != nil {
		log.Fatalf("Got error when checking for existing AOD: %s", err)
	}

	err = requestHandler.Db.Table("authors").Order("random()").Where("nr_of_icelandic_quotes > 0").Limit(1).First(&author).Error
	if err != nil {
		log.Fatalf("Got error when getting random english author: %s", err)
	}

	createAOD := author.ConvertToAODDBModel(today, true)
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

	var wg sync.WaitGroup
	wg.Add(5)
	go func() { defer wg.Done(); requestHandler.insertEnglishQOD(today) }()
	go func() { defer wg.Done(); requestHandler.insertIcelandicQOD(today) }()
	go func() { defer wg.Done(); requestHandler.insertEnglishAOD(today) }()
	go func() { defer wg.Done(); requestHandler.insertIcelandicAOD(today) }()
	go func() { defer wg.Done(); requestHandler.insertTopicsQOD(today) }()
	wg.Wait()
}

func main() {
	lambda.Start(theReqHandler.handler)
}
