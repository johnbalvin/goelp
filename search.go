package goelp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/johnbalvin/goelp/trace"
)

func Search(index int, searchValue, location, language string, sv SortValue, proxyURL *url.URL) (SearchOutput, error) {
	urlParsed, err := url.Parse(epSearch)
	if err != nil {
		return SearchOutput{}, trace.NewOrAdd(1, "main", "search", err, "")
	}
	datataSend := url.Values{}
	datataSend.Add("find_desc", searchValue)
	datataSend.Add("find_loc", location)
	datataSend.Add("start", strconv.Itoa(index))
	datataSend.Add("sortby", string(sv))
	datataSend.Add("hl", language)
	urlParsed.RawQuery = datataSend.Encode()
	urlToUse := urlParsed.String()
	resp, body, err := makeGetRequest(urlToUse, proxyURL)
	if err != nil {
		return SearchOutput{}, trace.NewOrAdd(2, "main", "search", err, "")
	}
	var redirectURL string
	if resp.StatusCode == 301 {
		location := resp.Header.Get("location")
		resp, body, err = makeGetRequest(location, proxyURL)
		if err != nil {
			return SearchOutput{}, trace.NewOrAdd(3, "main", "search", err, "")
		}
		redirectURL = location
	}
	if resp.StatusCode == 503 {
		return SearchOutput{}, trace.NewOrAdd(4, "main", "search", trace.ErrPermisions, "")
	}
	if resp.StatusCode != 200 {
		errData := fmt.Sprintf("status: %d headers: %+v", resp.StatusCode, resp.Header)
		return SearchOutput{}, trace.NewOrAdd(5, "main", "search", trace.ErrStatusCode, errData)
	}
	si, err := ParseBodySearch(body)
	if err != nil {
		return SearchOutput{}, trace.NewOrAdd(6, "main", "search", err, "")
	}
	si.SearchURL = urlToUse
	si.RedirectURL = redirectURL
	return si, nil
}
func makeGetRequest(urlToUse string, proxyURL *url.URL) (*http.Response, []byte, error) {
	req, err := http.NewRequest("GET", urlToUse, nil)
	if err != nil {
		return nil, nil, trace.NewOrAdd(1, "main", "makeGetRequest", err, "")
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "en")
	req.Header.Add("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	transport := &http.Transport{
		MaxIdleConnsPerHost: 30,
		DisableKeepAlives:   true,
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true,
		},
	}
	if proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}
	client := &http.Client{
		Timeout: time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, trace.NewOrAdd(2, "main", "makeGetRequest", err, "")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, trace.NewOrAdd(3, "main", "makeGetRequest", err, "")
	}
	return resp, body, nil
}
