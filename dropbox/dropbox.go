package dropbox
//Dropbox implementation of the actual website...

import (
    "strconv"
    "net/http"
    "html/template"
    "net/url"
    "io"
    "io/ioutil"
    "time"
    "encoding/json"
    "github.com/scritch007/chromecasa/chromecasa"
    "fmt"
)

var notDropboxAuthenticatedTemplate = template.Must(template.New("").Parse(`
<html><body>
You have currently not given permissions to access your data. Please authenticate this app with the Google OAuth provider.
<form action="/authorize?Dropbox=1" method="POST"><input type="submit" value="Ok, authorize this app with my id"/></form>
</body></html>
`));

var (
	clientId = "d81phd386hi5xs5"
	clientSecret = "o5yjaq9k5kah8b5"
	redirectURI = "http://127.0.0.1:3000/oauth2callback?provider=dropbox"
)


type Dropbox struct{
    TokenStore *chromecasa.TokenStore
}

func (d *Dropbox)HandleRoot(w http.ResponseWriter, r *http.Request){
    token := d.TokenStore.GetToken(r)

    if token == nil {
        //TODO remove previous cookie...
        notDropboxAuthenticatedTemplate.Execute(w, nil)
    }else{
        file_content, err := ioutil.ReadFile("./html/index.html")
        if err != nil {
            io.WriteString(w, "Failed to retrieve file")
            return
        }
        io.WriteString(w, string(file_content))
    }
}

type DropboxAccess struct{
	AccessToken string `json:"access_token"`
	Uid string `json:"uid"`

}

func (d *Dropbox) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request){
	code := r.FormValue("code")

	//Now get an access token
	remote_url := "https://api.dropbox.com/1/oauth2/token"
	values := url.Values{
		"client_id": {clientId},
		"client_secret": {clientSecret},
		"redirect_uri": {redirectURI},
		"grant_type": {"authorization_code"},
		"code": {code},
	}

	resp, err := http.PostForm(remote_url, values)
	if (nil != err){
		fmt.Println("Failed to get token")
		return
	}
	fmt.Println("Got something")
	raw_resp, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(raw_resp))

	var exchangedToken DropboxAccess
	//TODO Handle error properly
	err = json.Unmarshal(raw_resp, &exchangedToken)
	if (nil != err){
		fmt.Println(err)
		return
	}

	fmt.Println(exchangedToken)

    token := chromecasa.Token{AccessToken: exchangedToken.AccessToken, UID: exchangedToken.Uid}
    uid, _ := d.TokenStore.AddToken(token)
    expire := time.Now().AddDate(0, 0, 1)
    cookie := http.Cookie{Name: "chromecast_ref", Value: uid, Path: "/", Expires: expire}
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/?provider=dropbox", http.StatusFound)
}

func (d *Dropbox)HandleAlbum(w http.ResponseWriter, r *http.Request){
	album := r.FormValue("path")
    token := d.TokenStore.GetToken(r)

	client := new(http.Client)
	url := "https://api.dropbox.com/1/metadata/dropbox" + album + "?list=true"
	req, err := http.NewRequest("GET", url , nil)
	if(nil != err){
		io.WriteString(w, "Failed to create new Request")
		return;
	}
	req.Header.Add("Authorization", "Bearer " + token.AccessToken)
	resp, err :=client.Do(req)
	if(nil != err){
		io.WriteString(w, "Failed to send request")
		return;
	}
	raw_body, err := ioutil.ReadAll(resp.Body)
	if(nil != err){
		io.WriteString(w, "Failed to read response body")
		return;
	}

	io.WriteString(w, string(raw_body))

	m := new(MainDropbox)
	err = Parse(&raw_body, m)

    var result = make([]chromecasa.Folder, len(m.Content))

    for i := 0; i < len(m.Content); i++{
        i_str := strconv.Itoa(i)
        alb := chromecasa.Folder{Name:"Album" + string(i_str), Id:string(i_str), Icon: "/img/default.png", Display: true}
        result[i] = alb
    }
    b, _ := json.Marshal(result)
    io.WriteString(w, string(b))
}

func (d *Dropbox)HandleListAlbum(w http.ResponseWriter, r *http.Request){

    var result = make([]chromecasa.Image, 8)
    for i:=0;i<8;i++{
        i_str := strconv.Itoa(i)
        alb := chromecasa.Image{Name:"Image" + string(i_str), Icon: "/img/default_img.png", Content: "/img/default_content.png", Height: "128", Width: "128"}
        result[i] = alb
    }
    b, _ := json.Marshal(result)
    io.WriteString(w, string(b))
}

func (* Dropbox)HandleAuthorize(w http.ResponseWriter, r *http.Request){
	url_var := "https://www.dropbox.com/1/oauth2/authorize?client_id=" + clientId
	url_var += "&response_type=code"
	url_var += "&redirect_uri=" + url.QueryEscape(redirectURI)
 	url_var += "&state=nothing"
 	http.Redirect(w, r, url_var, http.StatusFound)
}

func (* Dropbox)HandleMain(w http.ResponseWriter, r *http.Request){

    file_content, err := ioutil.ReadFile("./html/Dropbox.html")
    if err != nil {
        io.WriteString(w, "Failed to retrieve file")
        return
    }
    io.WriteString(w, string(file_content))
}