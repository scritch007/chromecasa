package main


import (
    "code.google.com/p/goauth2/oauth"
    "github.com/gorilla/mux"
    "net/http"
    "html/template"
    "io"
    "io/ioutil"
    "fmt"
    "time"
    "os/exec"
    "github.com/scritch007/chromecasa/parser"
    "encoding/json"
)

var token_map = map[string]oauth.Token{}

var notAuthenticatedTemplate = template.Must(template.New("").Parse(`
<html><body>
You have currently not given permissions to access your data. Please authenticate this app with the Google OAuth provider.
<form action="/authorize" method="POST"><input type="submit" value="Ok, authorize this app with my id"/></form>
</body></html>
`));


var userInfoTemplate = template.Must(template.New("").Parse(`
<html><body>
This app is now authenticated to access your Google user info.  Your details are:<br />
{{.}}
</body></html>
`));

var mainTemplate = template.Must(template.New("").Parse(`
<html>
<head>
<script src="/js/main.js"></script>
</head>
<body onload="main()">
</body>
</html>
`));

// variables used during oauth protocol flow of authentication
var (
    code = ""
    token = ""
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
        RedirectURL: "http://127.0.0.1:3000/oauth2callback",

        //This is the 'scope' of the data that you are asking the user's permission to access. For getting user's info, this is the url that Google has defined.
        Scope: "https://www.googleapis.com/auth/userinfo.profile https://picasaweb.google.com/data/",
    }

//This is the URL that Google has defined so that an authenticated application may obtain the user's info in json format
const profileInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
const albumFeedURL = "https://picasaweb.google.com/data/feed/api/user/default"

func main() {
    r := mux.NewRouter()

    //DEBUG URL Handlers implementation defined in debug.go
    s := r.Queries("debug", "1").Subrouter()
    r.HandleFunc("/", handleDebugRoot)
    s.HandleFunc("/oauth2callback", handleDebugOAuthCallback)
    s.HandleFunc("/album", handleDebugAlbum)
    s.HandleFunc("/album/{id}", handleDebugListAlbum)
    s.HandleFunc("/authorize", handleDebugAuthorize)

    //PROD URL Handlers
    r.HandleFunc("/", handleRoot)
    r.HandleFunc("/js/{file}", handleJS)
    r.HandleFunc("/img/{file}", handleIMG)
    r.HandleFunc("/authorize", handleAuthorize)
    r.HandleFunc("/oauth2callback", handleOAuth2Callback)
    r.HandleFunc("/album", handleAlbum)
    r.HandleFunc("/album/{id}", handleListAlbum)
    http.Handle("/", r)
    //Google will redirect to this page to return your code, so handle it appropriately
    http.ListenAndServe("localhost:3000", nil)
}

func getToken(r *http.Request) *oauth.Token{
    cookie, _ := r.Cookie("chromecast_ref")

    fmt.Println(cookie)
    if cookie == nil {
        return nil
    }

    token, in_map := token_map[cookie.Value]
    
    fmt.Println(token_map)
    //refreshToken := session.Values["refresh_token"]
    if !in_map {
        //TODO remove previous cookie...
        return nil
    }else{
        return &token
    }
}

type Image struct{
    Name string `json:"name"`
    Icon string `json:"icon"`
    Content string `json:"content"`
}

type Album struct{
    Name string `json:"name"`
    Id string `json:"id"`
    Icon string `json:"icon"`
}


func handleAlbum(w http.ResponseWriter, r *http.Request){

    t := &oauth.Transport{Config: oauthCfg}
    token := getToken(r)
    if token == nil{
        io.WriteString(w, "[]")
        return
    }
    t.Token = token

    resp, err := t.Client().Get(albumFeedURL + "?alt=json")

    if( err != nil){
        fmt.Println("Got an error", err);
        //TODO set headers and return an error code
    }
    
    buf, err := ioutil.ReadAll(resp.Body)
    m, _ := chromecasa.Parse(buf)

    var result = make([]Album, len(m.Feed.Entries))
    for i, album := range m.Feed.Entries{
        alb := Album{Name:album.Name.Value, Id:album.Id.Value, Icon: album.Media.Icon[0].Url}
        result[i] = alb
    }
    b, _ := json.Marshal(result)
    io.WriteString(w, string(b))
}

func handleListAlbum(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    id := vars["id"]
    t := &oauth.Transport{Config: oauthCfg}
    token := getToken(r)
    if token == nil{
        io.WriteString(w, "[]")
        return
    }
    t.Token = token

    resp, err := t.Client().Get(albumFeedURL + "/albumid/" + id + "?alt=json")

    if( err != nil){
        fmt.Println("Got an error", err);
    }
    
    buf, err := ioutil.ReadAll(resp.Body)
    
    m, _ := chromecasa.Parse(buf)

    
    var result = make([]Image, len(m.Feed.Entries))
    for i, album := range m.Feed.Entries{
        alb := Image{Name:album.Title.Value, Icon: album.Media.Icon[0].Url, Content: album.Media.Content[0].Url}
        result[i] = alb
    }
    b, _ := json.Marshal(result)
    io.WriteString(w, string(b))   
}

func handleJS(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    file := vars["file"]
    file_content, err := ioutil.ReadFile("./js/" + string(file))
    if err != nil {
        io.WriteString(w, "Failed to retrieve file")
        return
    }
    io.WriteString(w, string(file_content))
}

func handleIMG(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    file := vars["file"]
    file_content, err := ioutil.ReadFile("./img/" + string(file))
    if err != nil {
        io.WriteString(w, "Failed to retrieve file")
        return
    }
    io.WriteString(w, string(file_content))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    token := getToken(r)
    
    if token == nil {
        //TODO remove previous cookie...
        notAuthenticatedTemplate.Execute(w, nil)    
    }else{
        mainTemplate.Execute(w, nil)
    }
    
}

// Start the authorization process
func handleAuthorize(w http.ResponseWriter, r *http.Request) {
    //Get the Google URL which shows the Authentication page to the user
    url := oauthCfg.AuthCodeURL("")

    //redirect user to that page
    http.Redirect(w, r, url, http.StatusFound)
}

// Function that handles the callback from the Google server
func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
    //TODO make the actual correct checks if value is OK and deal with the errors

    //Get the code from the response
    code := r.FormValue("code")

    t := &oauth.Transport{Config: oauthCfg}

    // Exchange the received code for a token
    token, _ := t.Exchange(code)

    out, err := exec.Command("uuidgen").Output()
    if err != nil {
        io.WriteString(w, "Failed to generate UUID")
        return
    }

    token_map[string(out[0:len(out) - 1])] = *token

    expire := time.Now().AddDate(0, 0, 1)
    cookie := http.Cookie{Name: "chromecast_ref", Value: string(out), Path: "/", Expires: expire}
    http.SetCookie(w, &cookie)
    
    http.Redirect(w, r, "/", http.StatusFound)
}
