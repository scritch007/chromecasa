package main
//DEBUG implementation of the actual website...

import (
    "strconv"
    "code.google.com/p/goauth2/oauth"
    "net/http"
    "html/template"
    "io"
    "time"
    "os/exec"
    "encoding/json"
)

var notDebugAuthenticatedTemplate = template.Must(template.New("").Parse(`
<html><body>
You have currently not given permissions to access your data. Please authenticate this app with the Google OAuth provider.
<form action="/authorize?debug=1" method="POST"><input type="submit" value="Ok, authorize this app with my id"/></form>
</body></html>
`));


func handleDebugRoot(w http.ResponseWriter, r *http.Request){
    token := getToken(r)
    
    if token == nil {
        //TODO remove previous cookie...
        notDebugAuthenticatedTemplate.Execute(w, nil)    
    }else{
        mainTemplate.Execute(w, nil)
    }
}

func handleDebugOAuthCallback(w http.ResponseWriter, r *http.Request){
    token := oauth.Token{}
    out, err := exec.Command("uuidgen").Output()
    if err != nil {
        io.WriteString(w, "Failed to generate UUID")
        return
    }
    token_map[string(out[0:len(out) - 1])] = token
    expire := time.Now().AddDate(0, 0, 1)
    cookie := http.Cookie{Name: "chromecast_ref", Value: string(out), Path: "/", Expires: expire}
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/?debug=1", http.StatusFound)
}

func handleDebugAlbum(w http.ResponseWriter, r *http.Request){
    var result = make([]Album, 3)

    for i := 0; i < 3; i++{
        i_str := strconv.Itoa(i)
        alb := Album{Name:"Album" + string(i_str), Id:string(i_str), Icon: "/img/default.png"}
        result[i] = alb
    }
    b, _ := json.Marshal(result)
    io.WriteString(w, string(b))
}

func handleDebugListAlbum(w http.ResponseWriter, r *http.Request){
    
    var result = make([]Image, 8)
    for i:=0;i<8;i++{
        i_str := strconv.Itoa(i)
        alb := Image{Name:"Image" + string(i_str), Icon: "/img/default_img.png", Content: "/img/default_content.png"}
        result[i] = alb
    }
    b, _ := json.Marshal(result)
    io.WriteString(w, string(b))   
}

func handleDebugAuthorize(w http.ResponseWriter, r *http.Request){

    http.Redirect(w, r, "/oauth2callback?debug=1", http.StatusFound)
}
