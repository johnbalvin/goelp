package goelp

type SearchOutput struct {
	Language    string
	SearchURL   string
	RedirectURL string
	Items       []SearchItem
}
type SearchItem struct {
	ID               string
	IsAd             bool
	Name             string
	SubTitle         string
	ShortDescription string
	YelpURL          string
	BizURL           string
	Rating           float32
	ReviewCount      int
	//yelp, rounds ResponseTimeMin to the biggest integer with a zero to the right(on the duration category), so if minutes is 12 on the screen is 20 minutes, if minutes is 26 on the screen is 30 minutes, if minutes is 68 on the screen is 2 hours,if minutes is 125 on the screen is 3 hours
	ResponseTimeMin int
	ManualLocation  ManualLocation
	AddressLines    []string
	Neighborhoods   []string
	Categories      []string
	City            string
	Licenses        []License
	Coordinates     Coordinates
}
type ManualLocation struct {
	Adress  string
	City    string
	State   string
	Zip     string
	Country string
}
type License struct {
	Licensee     string
	Number       string
	IssuedBy     string
	Trade        string
	VerifiedDate string
	ExpiryDate   string
}
type Coordinates struct {
	Latitude  float64
	Longitude float64
}
