package main


import (
    "github.com/gorilla/mux"
    "net/http"
    "html/template"
    "io"
    "io/ioutil"
    "./picasa"
    "github.com/scritch007/chromecasa/chromecasa"
    "./debug"
)

var userInfoTemplate = template.Must(template.New("").Parse(`
<html><body>
This app is now authenticated to access your Google user info.  Your details are:<br />
{{.}}
</body></html>
`));

func main() {
    r := mux.NewRouter()

    tStore := new(chromecasa.TokenStore)
    tStore.Init()

    m := new(Main)
    m.TokenStore = tStore
    d := new(debug.Debug)
    d.TokenStore = tStore

    p := new(picasa.Picasa)
    p.TokenStore = tStore

    //DEBUG URL Handlers implementation defined in debug.go
    s := r.Queries("debug", "1").Subrouter()
    s.HandleFunc("/", d.HandleRoot)
    s.HandleFunc("/oauth2callback", d.HandleOAuthCallback)
    s.HandleFunc("/album", d.HandleAlbum)
    s.HandleFunc("/album/{id}", d.HandleListAlbum)
    s.HandleFunc("/authorize", d.HandleAuthorize)
    s.HandleFunc("/debug", d.HandleMain)

    //PROD URL Handlers
    r.HandleFunc("/", m.handleRoot)
    r.HandleFunc("/js/{file}", m.handleJS)
    r.HandleFunc("/img/{file}", m.handleIMG)
    r.HandleFunc("/css/{file}", m.handleCSS)

    pic := r.Queries("provider", "picasa").Subrouter()
    pic.HandleFunc("/authorize", p.HandleAuthorize)
    pic.HandleFunc("/oauth2callback", p.HandleOAuth2Callback)
    pic.HandleFunc("/album", p.HandleAlbum)
    pic.HandleFunc("/album/{id}", p.HandleListAlbum)
    http.Handle("/", r)
    //Google will redirect to this page to return your code, so handle it appropriately
    http.ListenAndServe("localhost:3000", nil)
}

type Main struct{
    TokenStore *chromecasa.TokenStore
}

func (m *Main)handleJS(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    file := vars["file"]
    file_content, err := ioutil.ReadFile("./js/" + string(file))
    if err != nil {
        io.WriteString(w, "Failed to retrieve file")
        return
    }
    io.WriteString(w, string(file_content))
}

func (m *Main)handleIMG(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    file := vars["file"]
    file_content, err := ioutil.ReadFile("./img/" + string(file))
    if err != nil {
        io.WriteString(w, "Failed to retrieve file")
        return
    }
    io.WriteString(w, string(file_content))
}

func (m *Main)handleCSS(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    file := vars["file"]
    file_content, err := ioutil.ReadFile("./css/" + string(file))
    if err != nil {
        io.WriteString(w, "Failed to retrieve file")
        return
    }
    io.WriteString(w, string(file_content))
}

func (m *Main)handleRoot(w http.ResponseWriter, r *http.Request) {
    token := m.TokenStore.GetToken(r)

    var file string
    if token == nil {
        //TODO remove previous cookie...
        file = "./html/not_authenticated.html"
    }else{
        file = "./html/index.html"
    }
    file_content, err := ioutil.ReadFile(file)
    if err != nil {
        io.WriteString(w, "Failed to retrieve file")
        return
    }
    io.WriteString(w, string(file_content))

}
