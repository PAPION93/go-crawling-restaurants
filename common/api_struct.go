package common

// ResponseLocation is Response struct from Naver Location Api
type ResponseLocation struct {
	Total   int
	Start   int
	Display int
	Items   []struct {
		Title    string
		Link     string
		Category string
		Address  string
	}
}

// ResponseGeocoding is Response struct from Naver Geocoding Api
type ResponseGeocoding struct {
	Status string
	Meta   struct {
		TotalCount int
		Page       int
		Count      int
	}
	Addresses []struct {
		RoadAddress     string
		JibunAddress    string
		EnglishAddress  string
		AddressElements []struct {
			Type      [1]string
			LongName  string
			ShortName string
			Code      string
		}
		X        string
		Y        string
		Distance int
	}
	ErrorMessage string
}
