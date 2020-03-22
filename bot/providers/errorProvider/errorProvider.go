/*
Package "providers" is a package that provides required things reqired by the bot
to be used by other funcs.
This package should has all providers for the bot.
*/
package errorProvider

const (
	// BadRequest Error
	NotEnoffRight  = "Bad Request: not enough rights to restrict/unrestrict chat member"
	RepMsgNotFound = "Bad Request: reply message not found"
	MsgCannotDel   = "Bad Request: message can't be deleted"
	UserNotPart    = "Bad Request: USER_NOT_PARTICIPANT"
	UserInvalid    = "Bad Request: USER_ID_INVALID"
	ChatMemPriv    = "Bad Request: chat member status can't be changed in private chats"
	ChatNotFound   = "Bad Request: chat not found"

	// Forbidden Error
	BotKicked  = "Forbidden: bot was kicked"
	BotBlocked = "Forbidden: bot blocked by user"
)
