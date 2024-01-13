package goelp

type searchOutput struct {
	Locale      string `json:"locale"`
	URL         string `json:"-"`
	Domain      string
	LegacyProps legacyProps `json:"legacyProps"`
}
type legacyProps struct {
	SearchAppProps searchAppProps `json:"searchAppProps"`
}
type searchAppProps struct {
	SearchPageProps searchPageProps `json:"searchPageProps"`
}
type searchPageProps struct {
	MainContentComponentsListProps []mainContentComponentsListProp `json:"mainContentComponentsListProps"`
	RightRailProps                 rightRailProps                  `json:"rightRailProps"`
	PhotoMetadata                  []photoMetadata                 `json:"photoMetadata"`
}
type photoMetadata struct {
	BusinessEncid    string           `json:"businessEncid"`
	UploadedLocation uploadedLocation `json:"uploadedLocation"`
}
type uploadedLocation struct {
	Adress  string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}
type rightRailProps struct {
	SearchMapProps searchMapProps `json:"searchMapProps"`
}
type searchMapProps struct {
	MapState      mapState                 `json:"mapState"`
	HovercardData map[string]hovercardData `json:"hovercardData"`
}
type mapState struct {
	Markers []marker `json:"markers"`
}
type hovercardData struct {
	BizID         string   `json:"bizId"`
	Rating        float32  `json:"rating"`
	NumReviews    int      `json:"numReviews"`
	AddressLines  []string `json:"addressLines"`
	Neighborhoods []string `json:"neighborhoods"`
}
type marker struct {
	ResourceId string   `json:"resourceId"`
	BusinessId string   `json:"businessId"`
	Location   location `json:"location"`
}
type location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type mainContentComponentsListProp struct {
	BizId                string               `json:"bizId"`
	Snippet              snippet              `json:"snippet"`
	VerifiedLicenseInfo  verifiedLicenseInfo  `json:"verifiedLicenseInfo"`
	SearchResultBusiness searchResultBusiness `json:"searchResultBusiness"`
	ScrollablePhotos     scrollablePhotos     `json:"scrollablePhotos"`
	SearchActions        []searchAction       `json:"searchActions"`
}
type verifiedLicenseInfo struct {
	Licenses []license `json:"licenses"`
}
type license struct {
	Licensee     string `json:"licensee"`
	Number       string `json:"licenseNumber"`
	IssuedBy     string `json:"issuedBy"`
	Trade        string `json:"trade"`
	VerifiedDate string `json:"verifiedDate"`
	ExpiryDate   string `json:"expiryDate"`
}
type searchAction struct {
	Content content `json:"content"`
}
type content struct {
	ResponseTime responseTime `json:"responseTime"`
	SubtitleText subtitleText `json:"subtitleText"`
}
type subtitleText struct {
	Text string `json:"text"`
}
type responseTime struct {
	Value int64 `json:"value"`
}
type snippet struct {
	Text string `json:"text"`
}
type searchResultBusiness struct {
	IsAd bool   `json:"isAd"`
	Name string `json:"name"`
	//BusinessURL      string  `json:"businessUrl"`
	//FormattedAddress string  `json:"formattedAddress"`
	City    string `json:"city"`
	URLCode string `json:"alias"`
	//Rating      float32    `json:"rating"`
	Categories []category `json:"categories"`
	//	ReviewCount int        `json:"reviewCount"`
	Website website `json:"website"`
}
type category struct {
	Title string `json:"title"`
}
type website struct {
	Href string `json:"href"`
}
type scrollablePhotos struct {
	PhotoList []photoList `json:"photoList"`
}
type photoList struct {
	Src string `json:"src"`
}
