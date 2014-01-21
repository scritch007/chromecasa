package chromecasa

import (
	"encoding/json"
	//"fmt"
)

type TFeed struct {
	Value string `json:"$t"`
}

type Category struct {
	Term string `json:"term"`
	Scheme string `json:"scheme"`  	
}

type Image struct{
	Url string `json:"url"`
	Width int `json:"width"`
	Height int `json:"width"`
	Medium string `json:"medium"`
	Type string `json:"type"` 
}
type MediaGroup struct{
	Icon []Image `json:"media$thumbnail"`
	Content []Image `json:"media$content"`
}

type Entry struct {
	Category []Category `json:"category"`
	Id TFeed `json:"gphoto$id"`
	Name TFeed `json:"gphoto$name"`
	Media MediaGroup `json:"media$group"`
	Title TFeed `json:"title"`
}

type Feed struct{
	Entries []Entry `json:"entry"`
	Title TFeed `json:"title"`
	Id TFeed `json:"title"`
}

type MainResponse struct{
	Feed Feed `json:"feed"`
	Version string `json:"version"`
}



func Parse(input []byte) (*MainResponse, error){
	m := new(MainResponse)
	err := json.Unmarshal(input, m)
	//fmt.Println(temp)
	return m, err
}