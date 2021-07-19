package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
)

var testingHandler = RequestHandler{}

func TestHandler(t *testing.T) {
	testingHandler.InitializeDB()

	t.Cleanup(func() {
		testingHandler.Db.Exec("delete from aods")
		testingHandler.Db.Exec("delete from qods")
	})

	year, month, day := time.Now().Date()
	today := fmt.Sprintf("%d-%d-%d", year, month, day)

	t.Run("Insert English QOD", func(t *testing.T) {
		testingHandler.insertEnglishQOD(today)

		var quotes []structs.QuoteDBModel
		testingHandler.Db.Table("qods").
			Where("date = ?", today).
			Where("topic_id = 0").
			Not("is_icelandic").
			Find(&quotes)

		if len(quotes) != 1 {
			t.Fatalf("expected single QOD but got %d", len(quotes))
		}
	})

	t.Run("Insert twice, making sure to overwrite QOD if it already exists", func(t *testing.T) {
		testingHandler.insertEnglishQOD(today)
		testingHandler.insertEnglishQOD(today)
		var quotes []structs.QuoteDBModel
		testingHandler.Db.Table("qods").
			Where("date = ?", today).
			Where("topic_id = 0").
			Not("is_icelandic").
			Find(&quotes)

		if len(quotes) != 1 {
			t.Fatalf("expected single QOD but got %d quotes", len(quotes))
		}

	})

	t.Run("Insert Icelandic QOD", func(t *testing.T) {
		testingHandler.insertIcelandicQOD(today)

		var quotes []structs.QuoteDBModel
		testingHandler.Db.Table("qods").
			Where("date = ?", today).
			Where("topic_id = 0").
			Where("is_icelandic").
			Find(&quotes)

		if len(quotes) != 1 {
			t.Fatalf("expected single QOD but got %d", len(quotes))
		}
	})

	t.Run("Insert twice, making sure to overwrite QOD if it already exists", func(t *testing.T) {
		testingHandler.insertIcelandicQOD(today)
		testingHandler.insertIcelandicQOD(today)
		var quotes []structs.QuoteDBModel
		testingHandler.Db.Table("qods").
			Where("date = ?", today).
			Where("topic_id = 0").
			Where("is_icelandic").
			Find(&quotes)

		if len(quotes) != 1 {
			t.Fatalf("expected single QOD but got %d quotes", len(quotes))
		}

	})

	t.Run("Insert English AOD", func(t *testing.T) {
		testingHandler.insertEnglishAOD(today)

		var authors []structs.AuthorDBModel
		testingHandler.Db.Table("aods").
			Where("date = ?", today).
			Not("is_icelandic").
			Find(&authors)

		if len(authors) != 1 {
			t.Fatalf("expected single QOD but got %d", len(authors))
		}
	})

	t.Run("Insert twice, making sure to overwrite AOD if it already exists", func(t *testing.T) {
		testingHandler.insertEnglishAOD(today)
		testingHandler.insertEnglishAOD(today)
		var author []structs.AuthorDBModel
		testingHandler.Db.Table("aods").
			Where("date = ?", today).
			Not("is_icelandic").
			Find(&author)

		if len(author) != 1 {
			t.Fatalf("expected single QOD but got %d quotes", len(author))
		}

	})

	t.Run("Insert Icelandic AOD", func(t *testing.T) {
		testingHandler.insertIcelandicAOD(today)

		var authors []structs.AuthorDBModel
		testingHandler.Db.Table("aods").
			Where("date = ?", today).
			Where("is_icelandic").
			Find(&authors)

		if len(authors) != 1 {
			t.Fatalf("expected single QOD but got %d", len(authors))
		}
	})

	t.Run("Insert twice, making sure to overwrite icelandic AOD if it already exists", func(t *testing.T) {
		testingHandler.insertIcelandicAOD(today)
		testingHandler.insertIcelandicAOD(today)
		var author []structs.AuthorDBModel
		testingHandler.Db.Table("aods").
			Where("date = ?", today).
			Where("is_icelandic").
			Find(&author)

		if len(author) != 1 {
			t.Fatalf("expected single QOD but got %d quotes", len(author))
		}

	})

	t.Run("Insert Topics QOD", func(t *testing.T) {
		testingHandler.insertTopicsQOD(today)

		var quotes []structs.QuoteDBModel
		testingHandler.Db.Table("qods").
			Where("date = ?", today).
			Not("topic_id = 0").
			Find(&quotes)

		if len(quotes) != 131 {
			t.Fatalf("expected 132 QOD but got %d", len(quotes))
		}
	})

}
