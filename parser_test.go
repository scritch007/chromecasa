package main

import (
	"testing"
	"fmt"
	"io/ioutil"
	"./parser"
)

func TestParser(t *testing.T){
	js_byte, err := ioutil.ReadFile("test_files/album_feed_example.json")
	if(err != nil){
		t.Errorf("Failed to open the file")
	}
	
	output, err := chromecasa.Parse(js_byte)
	if err != nil{
		t.Errorf("Failed to load the informations", err)
	}
	fmt.Println("Title ", output.Feed.Title)
	fmt.Println("Version ", output.Version)
	fmt.Println("Num entries = ", len(output.Feed.Entries))
	fmt.Println(output.Feed.Entries[0].Name.Value, ": ", output.Feed.Entries[0].Id.Value, "=>", output.Feed.Entries[0].Media.Icon[0].Url)
	js_byte, err = ioutil.ReadFile("test_files/list_photo.json")
	if(err != nil){
		t.Errorf("Failed to open the file")
	}
	
	output, err = chromecasa.Parse(js_byte)
	if err != nil{
		t.Errorf("Failed to load the informations", err)
	}
	fmt.Println(output.Feed.Entries[0].Title.Value, ": ", output.Feed.Entries[0].Media.Content[0].Url, "=>", output.Feed.Entries[0].Media.Icon[0].Url)
}
