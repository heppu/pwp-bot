package wiki

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/heppu/jun/client"
	"github.com/heppu/pwp-bot/api"
	"github.com/heppu/pwp-bot/models"
	"github.com/sorcix/irc"
)

type WikiBot struct {
	irc     *client.Client
	api     *api.ApiClient
	channel string
}

func NewWikiBot(ircClient *client.Client, apiClient *api.ApiClient, channel string) (bot *WikiBot, err error) {
	bot = &WikiBot{
		irc:     ircClient,
		api:     apiClient,
		channel: channel,
	}
	bot.addCallbacks()
	return
}

func (w *WikiBot) addCallbacks() {
	w.irc.AddCallback("PRIVMSG", func(message *irc.Message) {
		log.Println("Name:", message.Name)
		log.Println("Trailing:", message.Trailing)
		go w.handleMessage(message)
	})
}

func (w *WikiBot) handleMessage(message *irc.Message) {
	// Check that this message was ment to our bot
	if !strings.HasPrefix(message.Trailing, PREFIX) {
		return
	}

	replyTo := message.Params[0]
	name := message.Name
	if replyTo == w.irc.Nickname {
		replyTo = name
	}

	// Split the message in array by first space
	msg := strings.SplitAfterN(message.Trailing, " ", 3)
	fmt.Printf("%s\n", message.Trailing)

	// The message was just !wiki without any command
	if len(msg) == 1 {
		w.irc.Privmsg(replyTo, HELP_MSG)
		return
	}

	// Call handler based on command
	switch strings.Replace(msg[1], " ", "", -1) {
	case COMMAND_HELP:
		w.handleHelp(msg[1:], replyTo)

	case COMMAND_REGISTER:
		w.handleRegister(msg[1:], name, replyTo)

	case COMMAND_UNREGISTER:
		w.handleUnregister(msg[1:], name, replyTo)

	case COMMAND_LOGIN:
		w.handleLogin(msg[1:], name, replyTo)

	case COMMAND_ARTICLE:
		w.handleArticle(msg[1:], name, replyTo)

	case COMMAND_COMMENT:
		w.handleComment(msg[1:], name, replyTo)

	default:
		w.irc.Privmsg(replyTo, INVALID_COMMAND)
	}
}

func (w *WikiBot) handleHelp(msg []string, replyTo string) {
	w.irc.Privmsg(replyTo, HELP_MSG)
}

func (w *WikiBot) handleRegister(msg []string, name, replyTo string) {
	if len(msg) < 2 {
		w.irc.Privmsg(replyTo, HELP_REGISTER)
		return
	}

	params := strings.Split(msg[1], " ")
	if len(params) != 2 {
		w.irc.Privmsg(replyTo, HELP_REGISTER)
		return
	}

	user := &models.User{params[0], params[1]}
	err := w.api.Register(name, user)
	if err != nil {
		w.irc.Privmsg(replyTo, fmt.Sprintf("Registration failed %s", err))
		return
	}
	w.irc.Privmsg(replyTo, "You have been registered successfully!")
}

func (w *WikiBot) handleLogin(msg []string, name, replyTo string) {
	if len(msg) < 2 {
		w.irc.Privmsg(replyTo, HELP_LOGIN)
		return
	}

	params := strings.Split(msg[1], " ")
	if len(params) != 2 {
		w.irc.Privmsg(replyTo, HELP_LOGIN)
		return
	}

	user := &models.User{params[0], params[1]}
	err := w.api.Auth(name, user)
	if err != nil {
		w.irc.Privmsg(replyTo, fmt.Sprintf("Login failed %s", err))
		return
	}
	w.irc.Privmsg(replyTo, "Logged in")
}

func (w *WikiBot) handleUnregister(msg []string, name, replyTo string) {

}

func (w *WikiBot) handleArticle(msg []string, name, replyTo string) {
	if len(msg) < 2 {
		w.irc.Privmsg(replyTo, HELP_ARTICLE)
		return
	}

	params := strings.SplitAfterN(msg[1], " ", 2)
	log.Println(params)
	switch strings.Replace(params[0], " ", "", -1) {
	case PARAM_LIST:
		articles, err := w.api.ListArticles(name)
		if err != nil {
			w.irc.Privmsg(replyTo, err.Error())
		}
		for _, a := range *articles {
			w.irc.Privmsg(replyTo, fmt.Sprintf("ID: %d : Topic: %s", a.Id, a.Topic))
		}

	case PARAM_ADD:
		if len(params) != 2 {
			w.irc.Privmsg(replyTo, HELP_ARTICLE_ADD)
			return
		}
		params := strings.SplitAfterN(params[1], " | ", 2)
		if len(params) != 2 {
			w.irc.Privmsg(replyTo, HELP_ARTICLE_ADD)
			return
		}
		article := &models.Article{
			Topic: params[0][0 : len(params[0])-2],
			Text:  params[1],
		}
		err := w.api.CreateArticle(name, article)
		if err != nil {
			w.irc.Privmsg(replyTo, err.Error())
			return
		}
		w.irc.Privmsg(replyTo, "Article added succesfully")

	case PARAM_SHOW:
		if len(params) != 2 {
			w.irc.Privmsg(replyTo, HELP_ARTICLE_SHOW)
			return
		}
		log.Println(params[1])

		id, err := strconv.Atoi(params[1])
		if err != nil {
			log.Println(err)
			w.irc.Privmsg(replyTo, HELP_ARTICLE_SHOW)
		}

		article, err := w.api.GetArticle(name, id)
		if err != nil {
			w.irc.Privmsg(replyTo, err.Error())
			return
		}
		w.irc.Privmsg(replyTo, fmt.Sprintf("ID: %d Topic: %s", article.Id, article.Topic))
		w.irc.Privmsg(replyTo, article.Text)

	case PARAM_DEL:
		if len(params) != 2 {
			w.irc.Privmsg(replyTo, HELP_ARTICLE_DEL)
			return
		}
		log.Println(params[1])

		id, err := strconv.Atoi(params[1])
		if err != nil {
			log.Println(err)
			w.irc.Privmsg(replyTo, HELP_ARTICLE_DEL)
		}

		if err = w.api.RemoveArticle(name, id); err != nil {
			w.irc.Privmsg(replyTo, err.Error())
			return
		}
		w.irc.Privmsg(replyTo, fmt.Sprintf("Article with ID:%d deleted", id))

	default:
		w.irc.Privmsg(replyTo, HELP_ARTICLE)

	}
}

func (w *WikiBot) handleComment(msg []string, name, replyTo string) {
	if len(msg) < 2 {
		w.irc.Privmsg(replyTo, HELP_COMMENT)
		return
	}

	params := strings.SplitAfterN(msg[1], " ", 2)
	log.Println(params)
	switch strings.Replace(params[0], " ", "", -1) {
	case PARAM_ADD:
		if len(params) != 2 {
			w.irc.Privmsg(replyTo, HELP_COMMENT)
			return
		}

		params := strings.SplitAfterN(params[1], " ", 2)
		if len(params) != 2 {
			w.irc.Privmsg(replyTo, HELP_COMMENT)
			return
		}

		id, err := strconv.Atoi(strings.Replace(params[0], " ", "", -1))

		if err != nil {
			log.Println(err)
			w.irc.Privmsg(replyTo, HELP_COMMENT)
		}

		comment := &models.NewComment{params[1]}
		if err = w.api.CreateComment(name, id, comment); err != nil {
			w.irc.Privmsg(replyTo, err.Error())
			return
		}
		w.irc.Privmsg(replyTo, "Comment added succesfully")

	case PARAM_SHOW:
		if len(params) != 2 {
			w.irc.Privmsg(replyTo, HELP_COMMENT)
			return
		}
		log.Println(params[1])

		id, err := strconv.Atoi(params[1])
		if err != nil {
			log.Println(err)
			w.irc.Privmsg(replyTo, HELP_COMMENT)
		}

		comments, err := w.api.ListComments(name, id)
		if err != nil {
			w.irc.Privmsg(replyTo, err.Error())
		}
		for _, c := range *comments {
			w.irc.Privmsg(replyTo, fmt.Sprintf("ID: %d : Topic: %s", c.Id, c.Text))
		}

	default:
		w.irc.Privmsg(replyTo, HELP_COMMENT)

	}
}
