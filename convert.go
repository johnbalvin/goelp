package goelp

import (
	"fmt"
	"log"
	"net/url"

	"github.com/johnbalvin/goelp/trace"
)

func (so searchOutput) standardize() SearchOutput {
	output := SearchOutput{
		Language: so.Locale,
	}
	pageProds := so.LegacyProps.SearchAppProps.SearchPageProps
	for _, prop := range pageProds.MainContentComponentsListProps {
		if prop.BizId == "" {
			continue
		}
		bizInfo := prop.SearchResultBusiness
		out := SearchItem{
			ID:   prop.BizId,
			Name: bizInfo.Name,
			IsAd: bizInfo.IsAd,
			//	Rating:           bizInfo.Rating,
			//ReviewCount:      bizInfo.ReviewCount,
			City:             bizInfo.City,
			ShortDescription: prop.Snippet.Text,
		}
		for _, licenseTmp := range prop.VerifiedLicenseInfo.Licenses {
			license := License(licenseTmp)
			out.Licenses = append(out.Licenses, license)
		}
		for _, cat := range prop.SearchResultBusiness.Categories {
			out.Categories = append(out.Categories, cat.Title)
		}
		if bizInfo.URLCode != "" {
			out.YelpURL = fmt.Sprintf("https://www.%s/biz/%s", so.Domain, bizInfo.URLCode)
		}
		for _, action := range prop.SearchActions {
			out.ResponseTimeMin = int(action.Content.ResponseTime.Value / 60)
			out.SubTitle = action.Content.SubtitleText.Text
		}
		out.BizURL = bizInfo.Website.Href
		if bizInfo.IsAd && bizInfo.Website.Href != "" {
			hrefUnecape, err := url.PathUnescape(bizInfo.Website.Href)
			if err == nil {
				bizInfo.Website.Href = hrefUnecape
				urlParsed, err := url.Parse(bizInfo.Website.Href)
				if err == nil {
					query := urlParsed.Query()
					bizURL := query.Get("url")
					out.BizURL = bizURL
				} else {
					errData := trace.NewOrAdd(1, "main", "Search", err, fmt.Sprintf("%+v", bizInfo))
					log.Println(errData)
				}
			} else {
				errData := trace.NewOrAdd(1, "main", "Search", err, fmt.Sprintf("%+v", bizInfo))
				log.Println(errData)
			}
		}
		output.Items = append(output.Items, out)
	}
	for _, marker := range pageProds.RightRailProps.SearchMapProps.MapState.Markers {
		for j, item := range output.Items {
			if item.ID == marker.ResourceId || item.ID == marker.BusinessId {
				output.Items[j].Coordinates.Latitude = marker.Location.Latitude
				output.Items[j].Coordinates.Longitude = marker.Location.Longitude
				break
			}
		}
	}
	for _, card := range pageProds.RightRailProps.SearchMapProps.HovercardData {
		for j, item := range output.Items {
			if item.ID == card.BizID {
				output.Items[j].AddressLines = card.AddressLines
				output.Items[j].Neighborhoods = card.Neighborhoods
				output.Items[j].Rating = card.Rating
				output.Items[j].ReviewCount = card.NumReviews
				break
			}
		}
	}
	for _, photoData := range pageProds.PhotoMetadata {
		for j, item := range output.Items {
			if item.ID == photoData.BusinessEncid {
				location := photoData.UploadedLocation
				output.Items[j].ManualLocation.Adress = location.Adress
				output.Items[j].ManualLocation.City = location.City
				output.Items[j].ManualLocation.State = location.State
				output.Items[j].ManualLocation.Zip = location.Zip
				output.Items[j].ManualLocation.Country = location.Country
				break
			}
		}
	}
	return output
}

func (wd wrapperData) standardize() Data {
	data := Data{
		Name:                wd.GeneralData.Name,
		Language:            wd.ExtraData.Locale,
		YelpURL:             wd.ExtraData.RequestUrl,
		BizURL:              wd.BizURL,
		TotalQuestionNumber: wd.QuestionsNumber,
		Claimed:             wd.Isclaimed,
		Telephone:           wd.LocalBusiness.Telephone,
		Reviews:             wd.Reviews,
		MainImage:           wd.LocalBusiness.Image,
		Images:              wd.Imgs,
		Qas:                 wd.QuestionsAnswers,
		Coordinates:         wd.Coordinates,
		Location: Location{
			Address:    wd.LocalBusiness.Address.StreetAddress,
			City:       wd.LocalBusiness.Address.AddressLocality,
			RegionCode: wd.LocalBusiness.Address.AddressRegion,
			Country:    wd.LocalBusiness.Address.AddressCountry,
			Zip:        wd.LocalBusiness.Address.PostalCode,
		},
		Rate: Rate{
			Value:       wd.GeneralData.Rating,
			ReviewCount: wd.GeneralData.ReviewCount,
		},
	}
	rt := wd.GeneralData.ReviewCountsByRating.Json
	if len(rt) == 5 {
		data.Rate.Star1 = rt[0]
		data.Rate.Start2 = rt[1]
		data.Rate.Start3 = rt[2]
		data.Rate.Start4 = rt[3]
		data.Rate.Start5 = rt[4]
	}
	for _, amenity := range wd.Amenities {
		am := Amenity{
			Title:    amenity.DisplayText,
			IsActive: amenity.IsActive,
			Icon:     amenity.IconName,
		}
		data.Amenities = append(data.Amenities, am)
	}
	return data
}
