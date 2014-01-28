package chromecasa

import (
	"net/http"
	"os/exec"
	"fmt"
	"sync"
)

type Image struct{
    Name string `json:"name"`
    Icon string `json:"icon"`
    Content string `json:"content"`
    Height string `json:"height"`
    Width string `json:"width"`
}

type Album struct{
    Name string `json:"name"`
    Id string `json:"id"`
    Icon string `json:"icon"`
}

type Token struct{
	Provider string
	AccessToken string
	RefreshToken string
	UID string
	Type string
}

type TokenStore struct {
    lock sync.RWMutex
    data map[string]Token // Data should probably not have any reference fields
}

func (c *TokenStore)Init(){
	c.data = map[string]Token{}
}

func (c *TokenStore) get(key string) (*Token, bool) {
    c.lock.RLock()
    defer c.lock.RUnlock()
    d, ok := c.data[key]
    return &d, ok
}

func (c *TokenStore) set(key string, d *Token) {
    c.lock.Lock()
    defer c.lock.Unlock()
    c.data[key] = *d
}

func (c *TokenStore) GetToken(r *http.Request) *Token{

    cookie, _ := r.Cookie("chromecast_ref")

    if cookie == nil {
        return nil
    }

    token, in_map := c.get(cookie.Value)
    fmt.Println("Here is the map", *c)

    //refreshToken := session.Values["refresh_token"]
    if !in_map {
        //TODO remove previous cookie...
        return nil
    }else{
        return token
    }
}

func (c *TokenStore)AddToken(token Token) (string, error){

    out, err := exec.Command("uuidgen").Output()
    if err != nil {
    	fmt.Println("Oups couldn't generate uuid")
        return "", err
    }
    c.set(string(out[0:len(out) - 1]), &token)
    fmt.Println("Adding token ", c)
    return string(out), err
}