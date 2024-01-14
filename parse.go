package goelp

import (
	"bytes"
	"encoding/json"
	"html"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/johnbalvin/goelp/trace"

	"github.com/PuerkitoBio/goquery"
)

func ParseBodySearch(body []byte) (SearchOutput, error) {
	out, err := parseBodySearch(body)
	if err != nil {
		return SearchOutput{}, trace.NewOrAdd(1, "main", "ParseBodySearch", err, "")
	}
	return out.standardize(), nil
}
func parseBodySearch(body []byte) (searchOutput, error) {
	reader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return searchOutput{}, trace.NewOrAdd(1, "main", "parseBodySearch", err, "")
	}
	var si searchOutput
	doc.Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		typeData := RemoveSpace(s.AttrOr("data-hypernova-key", ""))
		if typeData == "" {
			return true
		}
		htmlData, err := s.Html()
		if err != nil {
			errData := trace.NewOrAdd(2, "main", "parseBodySearch", err, "")
			log.Println(errData)
			return true
		}
		htmlData = RemoveSpace(html.UnescapeString(htmlData))
		htmlData = strings.TrimPrefix(htmlData, "<!--")
		htmlData = strings.TrimSuffix(htmlData, "-->")
		domain := regexDomain.FindString(htmlData)
		domain = strings.ReplaceAll(domain, `"domain":`, "")
		domain = strings.ReplaceAll(domain, `"`, "")
		si.Domain = domain
		if err := json.Unmarshal([]byte(htmlData), &si); err != nil {
			errData := trace.NewOrAdd(3, "main", "parseBodySearch", err, "")
			log.Println(errData)
			return true
		}
		return false
	})
	return si, nil
}
func ParseBodyDetails(body []byte) (Data, error) {
	out, err := parseBodyDetails(body)
	if err != nil {
		return Data{}, trace.NewOrAdd(1, "main", "ParseBodyDetails", err, "")
	}
	return out.standardize(), nil
}
func parseBodyDetails(body []byte) (wrapperData, error) {
	reader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return wrapperData{}, trace.NewOrAdd(1, "main", "parseBodyDetails", err, "")
	}
	var htmlDataLocalBussiness, htmlExtraData, htmlMainData string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		kind := s.AttrOr("type", "")
		htmlData, err := s.Html()
		if err != nil {
			errData := trace.NewOrAdd(2, "main", "parseBodyDetails", err, "")
			log.Println(errData)
			return
		}
		htmlData = RemoveSpace(html.UnescapeString(htmlData))
		htmlData = strings.TrimPrefix(htmlData, "<!--")
		htmlData = strings.TrimSuffix(htmlData, "-->")
		switch kind {
		case "application/ld+json":
			if strings.Contains(htmlData, `"@type":"LocalBusiness"`) {
				htmlDataLocalBussiness = htmlData
			}
		case "application/json":
			typeData1 := RemoveSpace(s.AttrOr("data-hypernova-id", ""))
			if typeData1 != "" {
				htmlExtraData = htmlData
				return
			}
			typeData2 := RemoveSpace(s.AttrOr("data-apollo-state", ""))
			if typeData2 == "" {
				return
			}
			htmlData = RemoveSpace(html.UnescapeString(htmlData))
			htmlMainData = htmlData
		}
	})
	data, err := parseMainDataDetails(htmlMainData)
	if err != nil {
		return wrapperData{}, trace.NewOrAdd(3, "main", "parseBodyDetails", err, "")
	}
	if err := json.Unmarshal([]byte(htmlDataLocalBussiness), &data.LocalBusiness); err != nil {
		return wrapperData{}, trace.NewOrAdd(4, "main", "parseBodyDetails", err, "")
	}
	var ed extraData
	if err := json.Unmarshal([]byte(htmlExtraData), &ed); err != nil {
		return wrapperData{}, trace.NewOrAdd(5, "main", "parseBodyDetails", err, "")
	}
	data.ExtraData.Locale = ed.Locale
	data.ExtraData.RequestUrl = ed.RequestUrl

	bizURL := regexBizURL.FindString(htmlMainData)
	bizURL = strings.ReplaceAll(bizURL, `"websiteUrl":"`, "")
	bizURL = strings.ReplaceAll(bizURL, `"`, "")
	data.BizURL = bizURL
	return data, nil
}
func parseMainDataDetails(htmlMainData string) (wrapperData, error) {
	var data wrapperData
	mapReviewsAutors := make(map[string]string)
	mapReviewsCreated := make(map[string]string)
	mapReviewsText := make(map[string]string)
	mapReviewsRating := make(map[string]int)
	mapQaQuestions := make(map[string]string)
	mapQaAnwers := make(map[string]answerData)
	mapQaCreated := make(map[string]string)
	mapQaAuthors := make(map[string]string)
	mainData := make(map[string]json.RawMessage)
	if err := json.Unmarshal([]byte(htmlMainData), &mainData); err != nil {
		return wrapperData{}, trace.NewOrAdd(1, "main", "parseMainDataDetails", err, "")
	}
	for label, value := range mainData {
		if strings.Contains(label, "$BusinessPhoto") && strings.HasSuffix(label, ".photoUrl") {
			mapData := make(map[string]string)
			if err := json.Unmarshal(value, &mapData); err != nil {
				return wrapperData{}, trace.NewOrAdd(2, "main", "parseMainDataDetails", err, "")
			}
			for label, value := range mapData {
				if strings.Contains(label, "LARGE") {
					img := Img{
						URL:       value,
						Extension: filepath.Ext(value),
					}
					data.Imgs = append(data.Imgs, img)
					break
				}
			}
			continue
		}
		if !strings.Contains(label, ".business") {
			continue
		}
		if strings.Contains(string(value), "reviewCount") && strings.HasSuffix(label, "})") {
			if err := json.Unmarshal(value, &data.GeneralData); err != nil {
				return wrapperData{}, trace.NewOrAdd(3, "main", "parseMainDataDetails", err, "")
			}
			continue
		}
		if strings.Contains(label, ".claimability") && strings.HasSuffix(label, "})") {
			var claimed claimed
			if err := json.Unmarshal(value, &claimed); err != nil {
				return wrapperData{}, trace.NewOrAdd(4, "main", "parseMainDataDetails", err, "")
			}
			data.Isclaimed = claimed.IsClaimed
			continue
		}
		if strings.Contains(label, ".operationHours") && strings.Contains(label, ".regularHoursMergedWithSpecialHoursForCurrentWeek") {
			var operationHourData operationHourData
			if err := json.Unmarshal(value, &operationHourData); err != nil {
				return wrapperData{}, trace.NewOrAdd(5, "main", "parseMainDataDetails", err, "")
			}
			data.OperationHourData = append(data.OperationHourData, operationHourData)
			continue
		}
		if strings.Contains(label, ".organizedProperties") && strings.Contains(label, ".properties") && strings.Contains(label, "clientPlatform") {
			var amenity amenity
			if err := json.Unmarshal(value, &amenity); err != nil {
				return wrapperData{}, trace.NewOrAdd(6, "main", "parseMainDataDetails", err, "")
			}
			data.Amenities = append(data.Amenities, amenity)
			continue
		}
		if strings.Contains(label, ".map") {
			mapData := make(map[string]string)
			if err := json.Unmarshal(value, &mapData); err != nil {
				return wrapperData{}, trace.NewOrAdd(7, "main", "parseMainDataDetails", err, "")
			}
			src := mapData["src"]
			if src == "" {
				continue
			}
			urlParsed, err := url.Parse(src)
			if err != nil {
				return wrapperData{}, trace.NewOrAdd(8, "main", "parseMainDataDetails", err, "")
			}
			query := urlParsed.Query()
			center := query.Get("center")
			splited := strings.Split(center, ",")
			if len(splited) == 2 {
				latitude, err := strconv.ParseFloat(splited[0], 64)
				if err != nil {
					return wrapperData{}, trace.NewOrAdd(9, "main", "parseMainDataDetails", err, "")
				}
				longitude, err := strconv.ParseFloat(splited[1], 64)
				if err != nil {
					return wrapperData{}, trace.NewOrAdd(10, "main", "parseMainDataDetails", err, "")
				}
				data.Coordinates.Latitude = latitude
				data.Coordinates.Longitude = longitude
			}
			continue
		}
		if strings.Contains(label, ".reviews") {
			if !(strings.Contains(label, "author") || strings.Contains(label, "createdAt") || strings.Contains(label, "text") || strings.Contains(label, "rating") || strings.HasSuffix(label, ".node")) {
				continue
			}
			if strings.HasSuffix(label, ".node") {
				var rating rating
				if err := json.Unmarshal(value, &rating); err != nil {
					return wrapperData{}, trace.NewOrAdd(11, "main", "parseMainDataDetails", err, "")
				}
				mapReviewsRating[label] = rating.Rating
				continue
			}
			mapData := make(map[string]string)
			if err := json.Unmarshal(value, &mapData); err != nil {
				return wrapperData{}, trace.NewOrAdd(12, "main", "parseMainDataDetails", err, "")
			}
			if strings.Contains(label, ".author") {
				name := mapData["displayName"]
				mapReviewsAutors[label] = name
			} else if strings.Contains(label, "createdAt") {
				for label2, value2 := range mapData {
					if strings.Contains(label2, "localDate") {
						mapReviewsCreated[label] = value2
						break
					}
				}
			} else if strings.Contains(label, "text") {
				review := mapData["full"]
				mapReviewsText[label] = review
			}
			continue
		}
		if strings.Contains(label, ".communityQuestions") {
			if strings.HasSuffix(label, "})") && strings.Contains(string(value), "edges") {
				var comunityQuestion comunityQuestion
				if err := json.Unmarshal(value, &comunityQuestion); err != nil {
					return wrapperData{}, trace.NewOrAdd(13, "main", "parseMainDataDetails", err, "")
				}
				data.QuestionsNumber = comunityQuestion.TotalCount
			} else if strings.HasSuffix(label, ".node") {
				var comunityQuestion comunityQuestion
				if err := json.Unmarshal(value, &comunityQuestion); err != nil {
					return wrapperData{}, trace.NewOrAdd(14, "main", "parseMainDataDetails", err, "")
				}
				mapQaQuestions[label] = comunityQuestion.Text
			} else if strings.HasSuffix(label, ".topAnswer") {
				var comunityQuestion comunityQuestion
				if err := json.Unmarshal(value, &comunityQuestion); err != nil {
					return wrapperData{}, trace.NewOrAdd(15, "main", "parseMainDataDetails", err, "")
				}
				mapQaAnwers[label] = answerData{
					Text:             comunityQuestion.Text,
					HelpfulVoteCount: comunityQuestion.HelpfulVoteCount,
				}
			} else if strings.HasSuffix(label, ".createdAt") {
				mapData := make(map[string]string)
				if err := json.Unmarshal(value, &mapData); err != nil {
					return wrapperData{}, trace.NewOrAdd(16, "main", "parseMainDataDetails", err, "")
				}
				mapQaCreated[label] = mapData["humanRelativeDateTime"]
			} else if strings.HasSuffix(label, ".author") {
				mapData := make(map[string]string)
				if err := json.Unmarshal(value, &mapData); err != nil {
					return wrapperData{}, trace.NewOrAdd(17, "main", "parseMainDataDetails", err, "")
				}
				mapQaAuthors[label] = mapData["displayName"]
			}
		}
	}
	for id, rating := range mapReviewsRating {
		rv := Review{
			Rate: rating,
		}
		for id2, date := range mapReviewsCreated {
			if strings.Contains(id2, id) {
				rv.Date = date
				break
			}
		}
		for id2, value := range mapReviewsText {
			if strings.Contains(id2, id) {
				rv.Value = value
				break
			}
		}
		for id2, authorName := range mapReviewsAutors {
			if strings.Contains(id2, id) {
				rv.Name = authorName
				break
			}
		}
		if rv.Name == "" {
			continue
		}
		data.Reviews = append(data.Reviews, rv)
	}
	for id, question := range mapQaQuestions {
		qa := Qa{
			Question: question,
		}
		for id2, awns := range mapQaAnwers {
			if strings.Contains(id2, id) {
				qa.Answer.Answer = awns.Text
				qa.Answer.HelpfullCount = awns.HelpfulVoteCount
				break
			}
		}
		for id2, date := range mapQaCreated {
			if strings.Contains(id2, id) {
				qa.Answer.Date = date
				break
			}
		}
		for id2, authorName := range mapQaAuthors {
			if strings.Contains(id2, id) {
				qa.Answer.AuthorName = authorName
				break
			}
		}
		data.QuestionsAnswers = append(data.QuestionsAnswers, qa)
	}
	return data, nil
}
