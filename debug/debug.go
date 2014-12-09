package debug
//DEBUG implementation of the actual website...

import (
    "strconv"
    "net/http"
    "html/template"
    "io"
    "io/ioutil"
    "time"
    "encoding/json"
    "github.com/scritch007/chromecasa/chromecasa"
    "math/rand"
    "fmt"
)

var notDebugAuthenticatedTemplate = template.Must(template.New("").Parse(`
<html><body>
You have currently not given permissions to access your data. Please authenticate this app with the Google OAuth provider.
<form action="/authorize?provider=debug" method="POST"><input type="submit" value="Ok, authorize this app with my id"/></form>
</body></html>
`));


type Debug struct{
    TokenStore *chromecasa.TokenStore
}

func (d *Debug)HandleRoot(w http.ResponseWriter, r *http.Request){
    token := d.TokenStore.GetToken(r)

    if token == nil {
        //TODO remove previous cookie...
        notDebugAuthenticatedTemplate.Execute(w, nil)
    }else{
        file_content, err := ioutil.ReadFile("./html/index.html")
        if err != nil {
            io.WriteString(w, "Failed to retrieve file")
            return
        }
        io.WriteString(w, string(file_content))
    }
}

func (d *Debug) HandleOAuthCallback(w http.ResponseWriter, r *http.Request){
    token := chromecasa.Token{}
    uid, _ := d.TokenStore.AddToken(token)
    expire := time.Now().AddDate(0, 0, 1)
    cookie := http.Cookie{Name: "chromecast_ref", Value: uid, Path: "/", Expires: expire}
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/?provider=debug", http.StatusFound)
}

func (d *Debug)HandleAlbum(w http.ResponseWriter, r *http.Request){
    loop := rand.Intn(10)
    var result = make([]chromecasa.Folder, loop)

    fmt.Println("Got this loop", loop)
    for i := 0; i < loop; i++{
        i_str := strconv.Itoa(i)
        alb := chromecasa.Folder{Name:"Album" + string(i_str), Id:string(i_str), Icon: "/img/default.png", Display: rand.Intn(2)==0, Browse: rand.Intn(2)==0}
        result[i] = alb
    }
    b, _ := json.Marshal(result)
    io.WriteString(w, string(b))
}

func (d *Debug)HandleListAlbum(w http.ResponseWriter, r *http.Request){

    var result = make([]chromecasa.Image, 8)
    for i:=0;i<8;i++{
        i_str := strconv.Itoa(i)
        alb := chromecasa.Image{Name:"Image" + string(i_str), Icon: "/img/default_img.png", Content: "/img/default_content.png", Height: "128", Width: "128"}
        result[i] = alb
    }
    b, _ := json.Marshal(result)
    io.WriteString(w, string(b))
}

func (* Debug)HandleAuthorize(w http.ResponseWriter, r *http.Request){

    http.Redirect(w, r, "/oauth2callback?provider=debug", http.StatusFound)
}

func (* Debug)HandleMain(w http.ResponseWriter, r *http.Request){

    file_content, err := ioutil.ReadFile("./html/debug.html")
    if err != nil {
        io.WriteString(w, "Failed to retrieve file")
        return
    }
    io.WriteString(w, string(file_content))
}