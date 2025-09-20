// Zone was generated automatically by zek 0.1.28. DO NOT EDIT.
type Zone struct {
	XMLName    xml.Name `xml:"zone"`
	Text       string   `xml:",chardata"`
	Script     string   `xml:"script"`
	Background struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
	} `xml:"background"`
	Layer []struct {
		Text    string `xml:",chardata"`
		W       string `xml:"w,attr"`
		H       string `xml:"h,attr"`
		Objects struct {
			Text string `xml:",chardata"`
			Foe  struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
				X    string `xml:"x,attr"`
				Y    string `xml:"y,attr"`
			} `xml:"foe"`
			Hidden struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
				X    string `xml:"x,attr"`
				Y    string `xml:"y,attr"`
			} `xml:"hidden"`
		} `xml:"objects"`
		Row []struct {
			Text string `xml:",chardata"`
			C    []struct {
				Text string `xml:",chardata"`
				I    string `xml:"i,attr"`
			} `xml:"c"`
		} `xml:"row"`
	} `xml:"layer"`
} 
