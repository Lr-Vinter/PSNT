package dbapi

import (
	"database/sql"
	"strconv"
)

type modifier func(*sql.Tx, int) error

func (d *DataBase) runModByCondition(tx *sql.Tx, modcond string) error {
	_, err := tx.Exec(`create temp table Mods` + modcond)

	if err != nil {
		_, err = tx.Exec(`create temp table TempData` + modcond)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`drop table Mods;`)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`alter table TempData RENAME TO Mods;`)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func (d *DataBase) ModExcludeSeenPosts(tx *sql.Tx, UserID int) error {

	modcond := ` AS
	select * from Mods exi

	left join SeenPosts 
		ON Mods.PostID = SeenPosts.PostID AND SeenPosts.FollowerID = ` + strconv.Itoa(UserID) +
		` where SeenPosts.FollowerID is NULL;`

	err := d.runModByCondition(tx, modcond)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) ModExcludeFromUser(tx *sql.Tx, UserID int) error {

	modcond := ` AS 
	select * from Mods 
	where Mods.OwnerID != ?;`

	err := d.runModByCondition(tx, modcond)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) ModLikesCount(tx *sql.Tx, UserID int) error {

	modcond := ` AS select Mods.PostID, Mods.ObjectID, Mods.OwnerID, Mods.Message, Mods.CreatedAt, count(Likes.ObjectID) as LikeNumber from Mods 
	left join Likes ON Mods.ObjectID = Likes.ObjectID group by Mods.ObjectID`

	err := d.runModByCondition(tx, modcond)
	if err != nil {
		return err
	}
	return nil
}
