package types

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	BASE10          = 10
	defaultHashCost = 14
)

type User struct {
	Username string `xorm:"varchar(32) not null unique index"`
	Hash     string `xorm:"varchar(128) not null"`
	ID       uint16 `xorm:"id pk autoincr"`
	Active   bool   `xorm:"active DEFAULT 1"`
	IsAdmin  bool   `xorm:"bool NOT NULL is_admin DEFAULT 0"`
}

func (u *User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultHashCost)
	return string(bytes), err
}

func (u *User) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Note struct {
	Title   string `xorm:"varchar(32) NOT NULL index('tag_note')"`
	Content string `xorm:"TEXT NOT NULL"`
	Tag     string `xorm:"varchar(32) NOT NULL index('tag_note')"`
	ID      uint16 `xorm:"id pk autoincr"`
	UserID  uint16 `xorm:"user_id INTEGER NOT NULL"`
}

type Notes []Note

func (n *Note) CalcHash() string {
	return ""
}

type SharedNote struct {
	ValidUntil time.Time `xorm:"valid_until DATE NOT NULL"`
	Hash       string    `xorm:"varchar(32) NOT NULL unique index('hash')"`
	NoteID     uint16    `xorm:"note_id INTEGER NOT NULL index"`
	UserID     uint16    `xorm:"user_id INTEGER NOT NULL index"`
}

type NoteTag struct {
	Alias    string
	Snippets uint64
}
