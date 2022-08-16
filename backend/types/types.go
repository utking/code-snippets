package types

type NoteTag struct {
	Alias string `xorm:"varchar(32) not null"`
	ID    uint16 `xorm:"id pk autoincr"`
	Notes uint   `xorm:"-"`
}

type Note struct {
	Title   string `xorm:"varchar(32) NOT NULL"`
	Content string `xorm:"TEXT NOT NULL"`
	TagID   uint16 `xorm:"tag_id NOT NULL"`
	ID      uint16 `xorm:"id pk autoincr"`
	Indent  uint8  `xorm:"NOT NULL"`
}

type Notes []Note

const (
	BASE10 = 10
)
