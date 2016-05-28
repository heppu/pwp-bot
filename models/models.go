package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Article struct {
	Topic string `json:"topic"`
	Text  string `json:"article_text"`
}

type ListArticle struct {
	Id    int    `json:"id"`
	Topic string `json:"topic"`
	Text  string `json:"article_text"`
}

type NewComment struct {
	Text string `json:"comment_text"`
}

type Comment struct {
	Id    int    `json:"id"`
	Owner int    `json:"topic"`
	Text  string `json:"comment_text"`
}
