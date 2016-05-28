package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/heppu/pwp-bot/models"
)

type ApiClient struct {
	tokens    map[string]string
	tokenLock *sync.RWMutex
	client    *http.Client
	apiUrl    string
}

func NewApiClient(apiUrl string) *ApiClient {
	tokens := make(map[string]string)
	return &ApiClient{
		tokens:    tokens,
		tokenLock: &sync.RWMutex{},
		client:    &http.Client{},
		apiUrl:    apiUrl,
	}
}

func (a *ApiClient) setTokenForNick(nick, token string) {
	a.tokenLock.Lock()
	defer a.tokenLock.Unlock()
	a.tokens[nick] = token
}

func (a *ApiClient) getTokenForNick(nick string) (token string, ok bool) {
	a.tokenLock.RLock()
	defer a.tokenLock.RUnlock()
	token, ok = a.tokens[nick]
	return
}

func (a *ApiClient) Register(nick string, user *models.User) (err error) {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(user)
	res, err := a.send("POST", a.apiUrl+"/users", nick, body)
	if err != nil {
		log.Println(err)
		return err
	}

	if res.StatusCode > 299 {
		log.Printf("Status not ok: %s\n", res.StatusCode)
		return errors.New(fmt.Sprintf("Request failed with code: %d", res.StatusCode))
	}

	token := res.Header.Get("Authorization")
	a.setTokenForNick(nick, token)

	return nil
}

func (a *ApiClient) Auth(nick string, user *models.User) (err error) {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(user)
	res, err := a.send("POST", a.apiUrl+"/auth", nick, body)
	if err != nil {
		log.Println(err)
		return err
	}

	if res.StatusCode > 299 {
		log.Printf("Status not ok: %s\n", res.StatusCode)
		return errors.New(fmt.Sprintf("Request failed with code: %d", res.StatusCode))
	}

	token := res.Header.Get("Authorization")
	a.setTokenForNick(nick, token)

	return nil
}

func (a *ApiClient) CreateArticle(nick string, article *models.Article) (err error) {

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(article)
	res, err := a.send("POST", a.apiUrl+"/articles", nick, body)
	if err != nil {
		log.Println(err)
		return err
	}

	if res.StatusCode > 299 {
		log.Printf("Status not ok: %s\n", res.StatusCode)
		return errors.New(fmt.Sprintf("Request failed with code: %d", res.StatusCode))
	}
	return nil
}

func (a *ApiClient) ListArticles(nick string) (*[]models.ListArticle, error) {
	res, err := a.send("GET", a.apiUrl+"/articles", nick, new(bytes.Buffer))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if res.StatusCode > 299 {
		log.Printf("Status not ok: %s\n", res.StatusCode)
		return nil, errors.New(fmt.Sprintf("Request failed with code: %d", res.StatusCode))
	}
	articles := make([]models.ListArticle, 0)
	if err = parseResponse(res, &articles); err != nil {
		return nil, err
	}
	return &articles, nil
}

func (a *ApiClient) GetArticle(nick string, id int) (*models.ListArticle, error) {
	addr := fmt.Sprintf("%s/articles/%d", a.apiUrl, id)
	res, err := a.send("GET", addr, nick, new(bytes.Buffer))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if res.StatusCode > 299 {
		log.Printf("Status not ok: %s\n", res.StatusCode)
		return nil, errors.New(fmt.Sprintf("Request failed with code: %d", res.StatusCode))
	}
	article := &models.ListArticle{}
	if err = parseResponse(res, &article); err != nil {
		return nil, err
	}
	return article, nil

}

func (a *ApiClient) RemoveArticle(nick string, id int) error {
	addr := fmt.Sprintf("%s/articles/%d", a.apiUrl, id)
	res, err := a.send("DELETE", addr, nick, new(bytes.Buffer))
	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode > 299 {
		log.Printf("Status not ok: %s\n", res.StatusCode)
		return errors.New(fmt.Sprintf("Request failed with code: %d", res.StatusCode))
	}
	return nil

}

func (a *ApiClient) CreateComment(nick string, id int, comment *models.NewComment) (err error) {
	addr := fmt.Sprintf("%s/articles/%d/comments", a.apiUrl, id)
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(comment)
	res, err := a.send("POST", addr, nick, body)
	if err != nil {
		log.Println(err)
		return err
	}

	if res.StatusCode > 299 {
		log.Printf("Status not ok: %s\n", res.StatusCode)
		return errors.New(fmt.Sprintf("Request failed with code: %d", res.StatusCode))
	}
	return nil
}

func (a *ApiClient) ListComments(nick string, id int) (*[]models.Comment, error) {
	addr := fmt.Sprintf("%s/articles/%d/comments", a.apiUrl, id)
	res, err := a.send("GET", addr, nick, new(bytes.Buffer))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if res.StatusCode > 299 {
		log.Printf("Status not ok: %s\n", res.StatusCode)
		return nil, errors.New(fmt.Sprintf("Request failed with code: %d", res.StatusCode))
	}
	comments := make([]models.Comment, 0)
	if err = parseResponse(res, &comments); err != nil {
		return nil, err
	}
	return &comments, nil
}

func (a *ApiClient) DeleteUser(nick string) error {
	addr := a.apiUrl + "/users/me"
	res, err := a.send("DELETE", addr, nick, new(bytes.Buffer))
	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode > 299 {
		log.Printf("Status not ok: %s\n", res.StatusCode)
		return errors.New(fmt.Sprintf("Request failed with code: %d", res.StatusCode))
	}
	return nil
}

func (a *ApiClient) send(method, path, nick string, body *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	token, ok := a.getTokenForNick(nick)

	if ok {
		req.Header.Set("Authorization", token)
	}

	return a.client.Do(req)
}

func parseResponse(r *http.Response, i interface{}) (err error) {
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(i)
	if err != nil {
		log.Printf("[error][api] : %s\n", err)
	}
	return
}
