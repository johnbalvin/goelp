package goelp

type wrapperData struct {
	QuestionsNumber   int
	Isclaimed         bool
	BizURL            string
	GeneralData       generalData
	LocalBusiness     localBusiness
	Amenities         []amenity
	OperationHourData []operationHourData
	QuestionsAnswers  []Qa
	Imgs              []Img
	Reviews           []Review
	ExtraData         extraData
	Coordinates       Coordinates
}

type extraData struct {
	Locale     string `json:"locale"`
	RequestUrl string `json:"requestUrl"`
}
type generalData struct {
	Name                      string               `json:"name"`
	ReviewCount               int                  `json:"reviewCount"`
	HasStorefrontAddress      bool                 `json:"hasStorefrontAddress"`
	Rating                    float32              `json:"rating"`
	ReviewCountsByRating      reviewCountsByRating `json:"reviewCountsByRating"`
	NotRecommendedReviewCount int                  `json:"notRecommendedReviewCount"`
}
type localBusiness struct {
	Name      string        `json:"name"`
	Image     string        `json:"image"`
	Telephone string        `json:"telephone"`
	Address   postalAddress `json:"address"`
}

type postalAddress struct {
	Type            string `json:"@type"`
	StreetAddress   string `json:"streetAddress"`
	AddressLocality string `json:"addressLocality"`
	AddressRegion   string `json:"addressRegion"`
	PostalCode      string `json:"postalCode"`
	AddressCountry  string `json:"addressCountry"`
}

type operationHourData struct {
	Hours           operationHour `json:"hours"`
	RegularHours    operationHour `json:"regularHours"`
	DayOfWeekShort  string        `json:"dayOfWeekShort"`
	HasSpecialHours bool          `json:"hasSpecialHours"`
}
type operationHour struct {
	Json []string `json:"json"`
}
type amenity struct {
	DisplayText string `json:"displayText"`
	IsActive    bool   `json:"isActive"`
	IconName    string `json:"iconName"`
}
type rating struct {
	Rating int `json:"rating"`
}
type comunityQuestion struct {
	TotalCount       int    `json:"totalCount"`
	Text             string `json:"text"`
	HelpfulVoteCount int    `json:"helpfulVoteCount"`
}
type answerData struct {
	Text             string
	HelpfulVoteCount int
}
type reviewCountsByRating struct {
	Json []int `json:"json"`
}
type claimed struct {
	IsClaimed bool `json:"isClaimed"`
}
