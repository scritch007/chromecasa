package dropbox

import (
	"encoding/json"
)
/*
{
	"hash": "4c64db2aeea4c3580660f3555e78a75d",
	"thumb_exists": false,
	"bytes": 0,
	"path": "/",
	"is_dir": true,
	"size": "0 bytes",
	"root": "dropbox",
	"contents": [
		{
			"revision": 465,
			"rev": "1d1015e492d",
			"thumb_exists": false,
			"bytes": 0,
			"modified": "Thu, 06 Feb 2014 21:24:28 +0000",
			"path": "/Photos",
			"is_dir": true,
			"icon": "folder",
			"root": "dropbox",
			"size": "0 bytes"
		},
		{
			"revision": 500,
			"rev": "1f4015e492d",
			"thumb_exists": true,
			"bytes": 47003,
			"modified": "Thu, 06 Feb 2014 21:28:05 +0000",
			"client_mtime": "Thu, 06 Feb 2014 21:28:05 +0000",
			"path": "/SeagateNasFinder.png",
			"is_dir": false,
			"icon": "page_white_picture",
			"root": "dropbox",
			"mime_type": "image/png",
			"size": "45.9 KB"
		},
		{
			"revision": 454,
			"rev": "1c6015e492d",
			"thumb_exists": false,
			"bytes": 0,
			"modified": "Tue, 27 Aug 2013 07:01:11 +0000",
			"path": "/test",
			"is_dir": true,
			"icon": "folder",
			"root": "dropbox",
			"size":
			"0 bytes"
		}],
	"icon": "folder"
}*/


type DropboxContent struct {
	IsDir bool `json:"is_dir"`
	Path string `json:"path"`
}
type MainDropbox struct{
	Content []DropboxContent `json:"content"`
	Path string `json:"path"`
}

func Parse(s *[]byte, m *MainDropbox) error{
	err := json.Unmarshal(*s, m)
	return err
}