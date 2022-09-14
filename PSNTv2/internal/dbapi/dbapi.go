package dbapi

import ( //"database/sql"

	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	db *sql.DB

	queryParam string
	answer     []Post
}

type LogicFunc func(int, int64, int, ...modifier) ([]Post, error)

func InitializeDataBase(driverName string, sourceName string) (*DataBase, error) {
	db, err := sql.Open(driverName, sourceName)
	if err != nil {
		return nil, err
	}

	return &DataBase{db, "", nil}, nil
}

// ------------------------------------------------------------main select args------>
// -----------------------------mod func arg(s)----->|                               |
func (d *DataBase) runQueryWithParams(tx *sql.Tx, UserID int, modfuncs []modifier, qargs ...interface{}) error {
	var res *sql.Rows
	len := len(modfuncs)
	var err error

	switch len {
	case 0:
		res, err = tx.Query(d.queryParam, qargs...)
		if err != nil {
			tx.Rollback()
			return err
		}
		d.loadAnswer(res, false)

	default:
		_, err = tx.Exec(`create temp table Mods AS`+d.queryParam, qargs...)
		if err != nil {
			tx.Rollback()
			return err
		}

		for i := range modfuncs {
			err = modfuncs[i](tx, UserID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		res = d.getDataFromTemp(tx)
		d.loadAnswer(res, true)
	}
	return nil
}

// *
// Followers Posts
func (d *DataBase) GetFollowerPosts(UserID int, time int64, number int, modfuncs ...modifier) ([]Post, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	d.queryParam = ` select Posts.PostID, Posts.ObjectID, Posts.OwnerID, Posts.Message, Posts.CreatedAt from Posts join Followers 
		ON Posts.OwnerID = Followers.TargetID AND Followers.SubID = ? 
		where Posts.CreatedAt > ?
		limit ?`

	err = d.runQueryWithParams(tx, UserID, modfuncs, UserID, time, number)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(fmt.Sprint("drop table ", "Mods"))
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return d.answer, nil // *sql.Rows link and ?
}

// *
// Posts with likes from target
func (d *DataBase) GetLikedByFollowerPosts(UserID int, time int64, number int, modfuncs ...modifier) ([]Post, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	d.queryParam = ` select distinct Posts.PostID, Posts.ObjectID, Posts.OwnerID, Posts.Message, Posts.CreatedAt from Posts left join Likes
		ON Posts.ObjectID = Likes.ObjectID
			left join Followers as LikeFromTarget
				ON Likes.UserID = LikeFromTarget.TargetID AND LikeFromTarget.SubID = ?
					where LikeFromTarget.SubID = ? AND Posts.CreatedAt > ?
					limit ?`

	err = d.runQueryWithParams(tx, UserID, modfuncs, UserID, UserID, time, number)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(fmt.Sprint("drop table ", "Mods"))
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return d.answer, nil
}

// ***
// Not good code
// ***
func (d *DataBase) ViewComments(UserID int, PostID int) ([]Post, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`
	WITH Answer AS 
	(
		select Posts.PostID from Posts
			left join Followers 
				ON Posts.OwnerID = Followers.TargetID AND followers.SubID = ? AND Posts.PostID = ?
	)
	select Comments.CommentID, Comments.OwnerID, Comments.Message, Comments.CreatedAt, count(Likes.ObjectID) as LikeNumber from Answer 
		join Comments 
			ON Comments.PostID = Answer.PostID
		left join Likes 
			ON Comments.ObjectID = Likes.ObjectID
	group by Comments.ObjectID
	`)
	defer stmt.Close()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	t := Post{}
	checknorows := stmt.QueryRow(UserID, PostID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = checknorows.Scan(&t.PostID, &t.OwnerID, &t.Message, &t.CreatedAt, &t.LikeNumber)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return nil, err
	}

	getcomments, err := stmt.Query(UserID, PostID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var comments []Post
	currpost := Post{}
	for getcomments.Next() {
		err = getcomments.Scan(&currpost.PostID, &currpost.OwnerID, &currpost.Message, &currpost.CreatedAt, &currpost.LikeNumber)
		comments = append(comments, currpost)
	}

	tx.Commit()
	return comments, nil
}

func (d *DataBase) loadAnswer(result *sql.Rows, likes bool) error {
	d.answer = nil
	currpost := Post{}
	objectid := 0

	var err error
	for result.Next() {
		if likes {
			err = result.Scan(&currpost.PostID, &objectid, &currpost.OwnerID,
				&currpost.Message, &currpost.CreatedAt, &currpost.LikeNumber)
		} else {
			err = result.Scan(&currpost.PostID, &objectid, &currpost.OwnerID,
				&currpost.Message, &currpost.CreatedAt)
		}
		if err != nil {
			return err
		}
		d.answer = append(d.answer, currpost)
	}

	return nil
}

func (d *DataBase) getDataFromTemp(tx *sql.Tx) *sql.Rows {
	res, err := tx.Query(`select Mods.PostID, Mods.ObjectID, Mods.OwnerID, Mods.Message, Mods.CreatedAt, LikeNumber from Mods`)
	if err != nil {
		panic(err)
	}
	return res
}
