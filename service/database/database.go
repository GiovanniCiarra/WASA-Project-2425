/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Creo la struct di un utente
type User struct {
	UserId int    `json:"userId"`
	Name   string `json:"username"`
	Photo  string `json:"photo"`
}

// Creo la struct di una conversation
type Conversation struct {
	ConversationId int     `json:"conversationId"`
	Type           string  `json:"type"`
	Name           string  `json:"name"`
	Photo          string  `json:"photo"`
	LastMessage    Message `json:"lastMessage"`
}

type ReplyInfo struct {
	MessageId      int    `json:"messageId"`
	Content        string `json:"content"`
	Type           string `json:"type"`
	SenderUsername string `json:"senderUsername"`
	Deleted        bool   `json:"deleted,omitempty"`
}

type Message struct {
	MessageId      int       `json:"messageId"`
	ConversationId int       `json:"conversationId"`
	Type           string    `json:"type"`
	Content        string    `json:"content"`
	SenderId       int       `json:"senderId"`
	SenderUsername string    `json:"senderUsername"`
	Timestamp      string    `json:"timestamp"`
	Status         string    `json:"status,omitempty"`
	Comments       []Comment `json:"comments,omitempty"`
	IsForwarded    int       `json:"isForwarded,omitempty"`
	IsReply        int       `json:"isReply,omitempty"`
	ReplyTo        ReplyInfo `json:"replyTo,omitempty" `
}

type Comment struct {
	CommentId int    `json:"commentId"`
	MessageId int    `json:"messageId"`
	SenderId  int    `json:"senderId"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetUserIdByName(name string) (int, error)
	CreateUser(name string) (int, error)
	GetUsers(userId int) ([]User, error)
	UpdateUserName(userId int, username string) error
	CheckUsernameExists(username string) (bool, error)
	UpdateUserPhoto(userId int, photo string) error
	GetMyConversations(userId int) ([]Conversation, error)
	GetConversation(userId int, conversationId int) ([]Message, error)
	GetOtherUserInformation(userId int, conversationId int) (string, string, error)
	CreateConversation(userId int, name string, photo string, members []int) (int, error)
	GetUserInformationById(userId int) (string, string, error)
	IsUserInConversation(userId int, conversationId int) (bool, error)
	SaveMessage(conversationId, senderId, isForwarded, replyTo int, content, msgType string) (int, error)
	IsMessageOwner(senderId, messageId, conversationId int) (bool, error)
	DeleteMessage(messageId int) error
	DoesConversationHaveMessages(conversationId int) (bool, error)
	DeleteConversation(conversationId int) error
	GetMessageById(messageId int) (Message, error)
	IsUserInMessageConversation(userId, messageId int) (bool, error)
	HasUserCommentedMessage(senderId, messageId int) (bool, error)
	SaveComment(messageId, senderId int, content string) (int, error)
	DeleteComment(messageId, senderId int) error
	AddToGroup(userId, conversationId int) error
	IsGroup(conversationId int) (bool, error)
	RemoveFromGroup(userId, conversationId int) error
	DoesGroupHaveMembers(conversationId int) (bool, error)
	UpdateGroupName(conversationId int, newName string) error
	UpdateGroupPhoto(conversationId int, image string) error
	UpdateMessagesToReceived(userId int) error
	UpdateMessagesToRead(userId, conversationId int) error
	GetConversationMembers(conversationId int) ([]User, error)
	UpdateMessageStatus(messageId int) error
	DoesUserExist(userId int) (bool, error)
	DoesConversationExists(userId, targetUserId int) (int, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='user_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE user_table(
										userId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
										name TEXT NOT NULL,
										photo TEXT NOT NULL,
										UNIQUE (name)
										);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='conversation_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE conversation_table(
										conversationId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
										type TEXT CHECK (type IN ('chat', 'group')) NOT NULL,
										name TEXT NOT NULL,
										photo TEXT NOT NULL
										);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='conversation_members_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE conversation_members_table(
										conversationId INTEGER NOT NULL,
										userId NOT NULL,
										PRIMARY KEY (conversationId, userId),
										FOREIGN KEY (conversationId) REFERENCES conversation_table(conversationId) ON DELETE CASCADE,
										FOREIGN KEY (userId) REFERENCES user_table(userId) ON DELETE CASCADE
										);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='message_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE message_table(
										messageId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
										conversationId INTEGER NOT NULL,
										senderId INTEGER NOT NULL,
										content TEXT,
										timestamp DATETIME,
										type TEXT NOT NULL CHECK(type IN ('text', 'photo')),
										status TEXT NOT NULL CHECK(status IN ('sent', 'received', 'read')),
										isForwarded INTEGER NOT NULL CHECK(isForwarded IN(0,1)),
										replyTo INTEGER NOT NULL,
										FOREIGN KEY (senderId) REFERENCES user_table(userId) ON DELETE CASCADE,
										FOREIGN KEY (conversationId) REFERENCES conversation_table(conversationId) ON DELETE CASCADE,
										FOREIGN KEY (replyTo) REFERENCES message_table(messageId) ON DELETE CASCADE
										);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comment_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE comment_table(
										commentId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
										messageId INTEGER NOT NULL,
										senderId INTEGER NOT NULL,
										content TEXT,
										timestamp DATETIME,
										UNIQUE (messageId, senderId),
										FOREIGN KEY (senderId) REFERENCES user_table(userId) ON DELETE CASCADE,
										FOREIGN KEY (messageId) REFERENCES message_table(messageId) ON DELETE CASCADE
										);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='message_status_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE message_status_table(
										messageId INTEGER NOT NULL,
										userId INTEGER NOT NULL,
										status TEXT NOT NULL CHECK(status IN ('sent', 'received', 'read')),
										PRIMARY KEY(messageId, userId),
										FOREIGN KEY (messageId) REFERENCES message_table(messageId) ON DELETE CASCADE,
										FOREIGN KEY (userId) REFERENCES user_table(userId) ON DELETE CASCADE
										);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
