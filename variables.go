package goelp

import "regexp"

const epSearch = "https://www.yelp.com/search"

type SortValue string

var (
	SortHighestRate = SortValue("rating")
	SortMostViewed  = SortValue("review_count")
	SortRecomended  = SortValue("recommended")
	regexBizURL     = regexp.MustCompile(`"websiteUrl":".+?"`)
	regexDomain     = regexp.MustCompile(`"domain":".+?"`)
)
