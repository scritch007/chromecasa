package picasa

import (
    "code.google.com/p/goauth2/oauth"
    "github.com/gorilla/mux"
    "net/http"
    "io"
    "io/ioutil"
    "fmt"
    "encoding/json"
    "github.com/scritch007/chromecasa/chromecasa"
    "time"
)

var oauthCfg = &oauth.Config {
        //TODO: put your project's Client Id here.  To be got from https://code.google.com/apis/console
        ClientId: "106373453700.apps.googleusercontent.com",

        //TODO: put your project's Client Secret value here https://code.google.com/apis/console
        ClientSecret: "x_1Ebngp5sfvKkB-vqN-Q260",

        //For Google's oauth2 authentication, use this defined URL
        AuthURL: "https://accounts.google.com/o/oauth2/auth",

        //For Google's oauth2 authentication, use this defined URL
        TokenURL: "https://accounts.google.com/o/oauth2/token",

        //To return your oauth2 code, Google will redirect the browser to this page that you have defined
        //TODO: This exact URL should also be added in your Google API console for this project within "API Access"->"Redirect URIs"
        RedirectURL: "http://127.0.0.1:3000/oauth2callback?provider=picasa",

        //This is the 'scope' of the data that you are asking the user's permission to access. For getting user's info, this is the url that Google has defined.
        Scope: "https://www.googleapis.com/auth/userinfo.profile https://picasaweb.google.com/data/",
    }

// variables used during oauth protocol flow of authentication
var (
    code = ""
    token = ""
)

//This is the URL that Google has defined so that an authenticated application may obtain the user's info in json format
const profileInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
const albumFeedURL = "https://picasaweb.google.com/data/feed/api/user/default"

type Picasa struct{
    TokenStore *chromecasa.TokenStore
}

func (p *Picasa)HandleAlbum(w http.ResponseWriter, r *http.Request){

    t := &oauth.Transport{Config: oauthCfg}
    token := p.TokenStore.GetToken(r)
    if token == nil{
        io.WriteString(w, "[]")
        return
    }
    t.Token = &oauth.Token{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken}

    resp, err := t.Client().Get(albumFeedURL + "?alt=json")

    if( err != nil){
        fmt.Println("Got an error", err);
        //TODO set headers and return an error code
    }

    buf, err := ioutil.ReadAll(resp.Body)
    m, _ := PicasaParse(buf)

    var result = make([]chromecasa.Folder, len(m.Feed.Entries))
    for i, album := range m.Feed.Entries{
        alb := chromecasa.Folder{Name:album.Name.Value, Id:album.Id.Value, Icon: album.Media.Icon[0].Url, Display: true, Browse: false}
        result[i] = alb
    }
    b, _ := json.Marshal(result)
    io.WriteString(w, string(b))
}

func (p *Picasa)HandleListAlbum(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    id := vars["id"]
    t := &oauth.Transport{Config: oauthCfg}
    token := p.TokenStore.GetToken(r)
    if token == nil{
        io.WriteString(w, "[]")
        return
    }
    t.Token = &oauth.Token{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken}

    //TODO we should be able to specify the imgmax, maybe depending on the network speed.
    //Need to find a way on how to find out the speed
    //For now the 1600 should be enough
    resp, err := t.Client().Get(albumFeedURL + "/albumid/" + id + "?alt=json&imgmax=1600")

    if( err != nil){
        fmt.Println("Got an error", err);
        io.WriteString(w, string("{\"error\": \"Failed to retrieve album information\"}"))
        return
    }

    buf, err := ioutil.ReadAll(resp.Body)

    fmt.Println(string(buf))

    m, _ := PicasaParse(buf)

    var result = make([]chromecasa.Image, len(m.Feed.Entries))
    for i, album := range m.Feed.Entries{
        alb := chromecasa.Image{Name:album.Title.Value, Icon: album.Media.Icon[0].Url, Content: album.Media.Content[0].Url, Height: album.Height.Value, Width:album.Width.Value}
        result[i] = alb
    }
    b, _ := json.Marshal(result)
    io.WriteString(w, string(b))
}

// Start the authorization process
func (p *Picasa)HandleAuthorize(w http.ResponseWriter, r *http.Request) {
    //Get the Google URL which shows the Authentication page to the user
    url := oauthCfg.AuthCodeURL("")

    //redirect user to that page
    http.Redirect(w, r, url, http.StatusFound)
}

// Function that handles the callback from the Google server
func (p *Picasa)HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
    //TODO make the actual correct checks if value is OK and deal with the errors

    //Get the code from the response
    code := r.FormValue("code")

    t := &oauth.Transport{Config: oauthCfg}

    // Exchange the received code for a token
    token, err := t.Exchange(code)

    if (err != nil){
        io.WriteString(w, "Failed to get Credentials")
        return
    }
    new_t := chromecasa.Token{Provider:"picasa", AccessToken:token.AccessToken, RefreshToken:token.RefreshToken}
    uid, _ := p.TokenStore.AddToken(new_t)
    expire := time.Now().AddDate(0, 0, 1)
    cookie := http.Cookie{Name: "chromecast_ref", Value: uid, Path: "/", Expires: expire}
    http.SetCookie(w, &cookie)

    http.Redirect(w, r, "/?provider=picasa", http.StatusFound)
}
