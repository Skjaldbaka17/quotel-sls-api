package utils

import "github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"

//Functions below are for "offline" updating Database
//Points incremented for appearing in search list
const incrementAppearInSearchList = 1

//Points incremented for direct get
const incrementIdFetch = 10

//DirectFetchAuthorsCountIncrement increments the popularity count of the Authors from a id-query
func (requestHandler *RequestHandler) DirectFetchAuthorsCountIncrement(authorIds []int) error {
	if len(authorIds) == 0 {
		return nil
	}
	return requestHandler.Db.Exec("UPDATE authors SET count = count + ? where id in (?) returning *", incrementIdFetch, authorIds).Error
}

//DirectFetchQuotesCountIncrement increments the popularity count of the Quotes from a id-query
func (requestHandler *RequestHandler) DirectFetchQuotesCountIncrement(quoteIds []int) error {
	if len(quoteIds) == 0 {
		return nil
	}
	return requestHandler.Db.Exec("UPDATE quotes SET count = count + ? where id in (?) returning *", incrementIdFetch, quoteIds).Error
}

//DirectFetchTopicCountIncrement increments the popularity count of the Topic from a id- or name-query
func (requestHandler *RequestHandler) DirectFetchTopicCountIncrement(topicId int, topicName string) error {
	return requestHandler.Db.Exec("UPDATE topics SET count = count + ? where id = ? or lower(name) = lower(?) returning *", incrementIdFetch, topicId, topicName).Error
}

//AuthorsAppearInSearchCountIncrement increments the popularity count of the Authors from a listing in a search
func (requestHandler *RequestHandler) AuthorsAppearInSearchCountIncrement(authors []structs.AuthorDBModel) error {
	if len(authors) == 0 {
		return nil
	}
	authorIds := []int{}

	for _, author := range authors {
		authorIds = append(authorIds, author.Id)
	}

	return requestHandler.Db.Exec("UPDATE authors SET count = count + ? where id in (?) returning *", incrementAppearInSearchList, authorIds).Error
}

//QuotesAppearInSearchCountIncrement increments the popularity count of the Quotes from a listing in a search
func (requestHandler *RequestHandler) QuotesAppearInSearchCountIncrement(quotes []structs.SearchViewDBModel) error {
	if len(quotes) == 0 {
		return nil
	}
	quoteIds := []int{}

	for _, quote := range quotes {
		quoteIds = append(quoteIds, quote.QuoteId)
	}

	return requestHandler.Db.Exec("UPDATE quotes SET count = count + ? where id in (?) returning *", incrementAppearInSearchList, quoteIds).Error
}

//AppearInSearchCountIncrement increments the popularity count of the Authors and quotes from a listing in a search
func (requestHandler *RequestHandler) TopicViewAppearInSearchCountIncrement(quotes []structs.TopicViewDBModel) error {
	if len(quotes) == 0 {
		return nil
	}
	authorIds := []int{}
	quoteIds := []int{}
	for _, quote := range quotes {
		authorIds = append(authorIds, quote.AuthorId)
		quoteIds = append(quoteIds, quote.QuoteId)
	}

	err := requestHandler.Db.Exec("UPDATE authors SET count = count + ? where id in (?) returning *", incrementAppearInSearchList, authorIds).Error
	if err != nil {
		return err
	}
	err = requestHandler.Db.Exec("UPDATE quotes SET count = count + ? where id in (?) returning *", incrementAppearInSearchList, quoteIds).Error
	return err
}

//AppearInSearchCountIncrement increments the popularity count of the Authors and quotes from a listing in a search
func (requestHandler *RequestHandler) SearchViewAppearInSearchCountIncrement(quotes []structs.SearchViewDBModel) error {
	if len(quotes) == 0 {
		return nil
	}
	authorIds := []int{}
	quoteIds := []int{}
	for _, quote := range quotes {
		authorIds = append(authorIds, quote.AuthorId)
		quoteIds = append(quoteIds, quote.QuoteId)
	}

	err := requestHandler.Db.Exec("UPDATE authors SET count = count + ? where id in (?) returning *", incrementAppearInSearchList, authorIds).Error
	if err != nil {
		return err
	}
	err = requestHandler.Db.Exec("UPDATE quotes SET count = count + ? where id in (?) returning *", incrementAppearInSearchList, quoteIds).Error
	return err
}
