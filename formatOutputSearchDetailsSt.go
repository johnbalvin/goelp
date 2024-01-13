package goelp

type Data struct {
	Language            string
	TotalQuestionNumber int
	Name                string
	Claimed             bool
	Telephone           string
	YelpURL             string
	BizURL              string
	MainImage           string
	Rate                Rate
	Location            Location
	Coordinates         Coordinates
	Images              []Img
	Amenities           []Amenity
	Qas                 []Qa
	Reviews             []Review
}
type Rate struct {
	Value       float32
	ReviewCount int
	Star1       int
	Start2      int
	Start3      int
	Start4      int
	Start5      int
}
type Qa struct {
	Question string
	Answer   Answer
}
type Answer struct {
	Answer        string
	AuthorName    string
	Date          string
	HelpfullCount int
}
type Amenity struct {
	Title    string
	IsActive bool
	Icon     string
}
type Img struct {
	URL         string
	ContentType string
	Extension   string
	Content     []byte `json:"-"`
}
type Review struct {
	Name  string
	Rate  int
	Date  string
	Value string
}
type Location struct {
	Address    string
	City       string
	RegionCode string
	Zip        string
	Country    string
}
