package goelp

import "net/url"

type Client struct {
	Language  string    //for example es_CL,es_MX
	SortValue SortValue //for example SortHighestRate,SortMostViewed,SortRecomended
	Location  string
	ProxyURL  *url.URL
}

func DefaulClient() Client {
	client := Client{
		Language:  "en_US",
		SortValue: SortRecomended,
		Location:  "Austin, TX 78701", //default, yelp requeries to input a location, this is random, change it if you need it from other location
		ProxyURL:  nil,
	}
	return client
}
func NewClient(language, location string, sortValue SortValue, proxyURL *url.URL) Client {
	client := Client{
		Language:  language,
		SortValue: sortValue,
		Location:  location,
		ProxyURL:  proxyURL,
	}
	return client
}

func (cl Client) Search(index int, searchValue string) (SearchOutput, error) {
	return Search(index, searchValue, cl.Location, cl.Language, cl.SortValue, cl.ProxyURL)
}

func (cl Client) GetFromYelpBizURL(yelpURL string) (Data, error) {
	return GetFromYelpBizURL(yelpURL, cl.ProxyURL)
}
