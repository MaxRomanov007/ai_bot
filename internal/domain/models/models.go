package models

type UserRole struct {
	UserRoleID   int64  `db:"user_role_id"`
	UserRoleName string `db:"user_role_name"`
}

type UserStatus struct {
	UserStatusID   int64  `db:"user_status_id"`
	UserStatusName string `db:"user_status_name"`
}

type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	ChatID   int64  `db:"chat_id"`
	RoleID   int64  `db:"role_id"`
	StatusID int64  `db:"status_id"`

	Role   UserRole   `db:"-"`
	Status UserStatus `db:"-"`
}

type MessageRole struct {
	MessageRoleID   int64  `db:"message_role_id"`
	MessageRoleName string `db:"message_role_name"`
}

type Message struct {
	MessageID int64  `db:"message_id"`
	UserID    int64  `db:"user_id"`
	RoleID    int64  `db:"role_id"`
	Content   string `db:"content"`

	User User        `db:"-"`
	Role MessageRole `db:"-"`
}

const (
	UserRoleUser           = "user"
	UserRoleAdmin          = "admin"
	UserStatusAuthorized   = "authorized"
	UserStatusUnauthorized = "unauthorized"
	UserStatusBlocked      = "blocked"
	MessageRoleUser        = "user"
	MessageRoleAssistant   = "assistant"
)
