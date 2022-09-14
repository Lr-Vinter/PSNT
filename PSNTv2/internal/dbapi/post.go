package dbapi

type Post struct {
	PostID     int
	OwnerID    int
	CreatedAt  int64
	LikeNumber int
	Message    string
}

type Field struct {
	Name  string
	Value interface{}
}

type Selection []Post
