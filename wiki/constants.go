package wiki

const (
	PREFIX = "!wiki"

	COMMAND_HELP       = "help"
	COMMAND_REGISTER   = "register"
	COMMAND_UNREGISTER = "unregister"
	COMMAND_LOGIN      = "login"
	COMMAND_ARTICLE    = "article"
	COMMAND_COMMENT    = "comment"

	INVALID_COMMAND = "SKRIIK SKRAAK SKRILLEX d:D"

	PARAM_LIST = "ls"
	PARAM_ADD  = "add"
	PARAM_DEL  = "del"
	PARAM_SHOW = "show"

	HELP_MSG          = "Usage !wiki [help | register | unregister | login | article | comment]"
	HELP_HELP         = "!wiki help [register | unregister | login | article | comment]"
	HELP_REGISTER     = "!wiki register [username] [password]"
	HELP_UNREGISTER   = "!wiki unregister [username] [password]"
	HELP_LOGIN        = "!wiki login [username] [password]"
	HELP_ARTICLE      = "!wiki article [ls | add | show | del]"
	HELP_ARTICLE_ADD  = "!wiki article add [topic] | [content]"
	HELP_ARTICLE_SHOW = "!wiki article show [id]"
	HELP_ARTICLE_DEL  = "!wiki article del [id]"
	HELP_COMMENT      = "!wiki comment [add | show] [article_id]"
)
