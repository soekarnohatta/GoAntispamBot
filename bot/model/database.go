package model

type (
	// Struct Lang holds the language preference of a chat or an user.
	Lang struct {
		ChatID   int    `bson:"chat_id" json:"chat_id"`
		Language string `bson:"lang" json:"lang"`
	}

	// Struct Private holds the notification preference of an user.
	Private struct {
		UserID int  `bson:"user_id" json:"user_id"`
		Notif  bool `bson:"notification" json:"notification"`
	}

	// Struct GlobalBan holds the global ban data of an user.
	GlobalBan struct {
		UserID     int    `bson:"user_id" json:"user_id"`
		Reason     string `bson:"reason" json:"reason_id"`
		BannedBy   int    `bson:"banner" json:"banner"`
		BannedFrom int    `bson:"appeal" json:"appeal"`
		TimeAdded  int    `bson:"time_added" json:"time_added"`
	}

	// Struct GroupSetting holds the settings data of a chat.
	GroupSetting struct {
		ChatID         int  `bson:"chat_id" json:"chat_id"`
		Gban           bool `bson:"enforce_gban" json:"enforce_gban"`
		Username       bool `bson:"enforce_username" json:"enforce_username"`
		ProfilePicture bool `bson:"enforce_profile_picture" json:"enforce_profile_picture"`
		Time           int  `bson:"time_settings" json:"time_settings"`
	}

	LockSetting struct {
	}

	// Struct ChatLog holds the basic chat data of a chat.
	ChatLog struct {
		ChatID    int    `bson:"chat_id" json:"chat_id"`
		ChatType  string `bson:"chat_type" json:"chat_type"`
		ChatLink  string `bson:"chat_link" json:"chat_link"`
		ChatTitle string `bson:"chat_title" json:"chat_title"`
	}

	// Struct UserLog holds the basic user data of an user.
	UserLog struct {
		UserID    int    `bson:"user_id" json:"user_id"`
		FirstName string `bson:"user_first_name" json:"user_first_name"`
		LastName  string `bson:"user_last_name" json:"user_last_name"`
		UserName  string `bson:"user_username" json:"user_username"`
	}
)
