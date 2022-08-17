package types

import "golang.org/x/crypto/bcrypt"

const (
	BASE10          = 10
	defaultHashCost = 14
)

type User struct {
	Username string `xorm:"varchar(32) not null"`
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
	Title   string `xorm:"varchar(32) NOT NULL index(tag)"`
	Content string `xorm:"TEXT NOT NULL"`
	Tag     string `xorm:"varchar(32) NOT NULL index(tag)"`
	ID      uint16 `xorm:"id pk autoincr"`
	UserID  uint16 `xorm:"user_id INTEGER NOT NULL"`
}

type Notes []Note
type NoteTag struct {
	Alias    string
	Snippets uint64
}
type NoteTags []NoteTag
