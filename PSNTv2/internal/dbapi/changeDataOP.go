package dbapi

import (
	"errors"
	"strconv"
)

// Normal delete
// add int64
func (d *DataBase) DeleteDataByFields(tablename string, fields ...Field) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	query := "delete from " + tablename + " where "
	for i, field := range fields {
		if i != 0 {
			query = query + " AND "
		}

		if strval, ok := field.Value.(string); ok {
			query = query + field.Name + " = " + strval
		} else {
			if numval, ok := field.Value.(int); ok {
				query = query + field.Name + " = " + strconv.Itoa(numval)
			} else {
				return errors.New("wrong arg passed")
			}
		}
	}
	_, err = tx.Exec(query + ";")
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Normal Insert
// add int64
func (d *DataBase) InsertDataByFields(tablename string, fields ...Field) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	queryNames := "insert into " + tablename + "("
	queryValues := " values("
	for i, field := range fields {

		queryNames = queryNames + field.Name

		if strval, ok := field.Value.(string); ok {
			queryValues = queryValues + strval
		} else {
			if numval, ok := field.Value.(int); ok {
				queryValues = queryValues + strconv.Itoa(numval)
			} else {
				return errors.New("wrong arg passed")
			}
		}
		if i == len(fields)-1 {
			queryNames += ")"
			queryValues += ");"
			break
		}
		queryNames += ", "
		queryValues += ", "
	}
	_, err = tx.Exec(queryNames + queryValues)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

/*
// delete with 2 tables
func (d *DataBase) PushPost(UserID int, Message string, CreatedAt int64) error {
res, err := tx.Exec(`INSERT INTO MessageObjects(Type) VALUES ('post');`)
ObjectId, err := res.LastInsertId()
insertpost, err := tx.Prepare(`

	INSERT INTO Posts(ObjectID, OwnerID, Message, CreatedAt)
	VALUES (?, ?, ?, ?);

`)
*/
