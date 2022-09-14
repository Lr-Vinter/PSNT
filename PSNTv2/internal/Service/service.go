package service

import (
	"feed/internal/dbapi"
	"sort"
)

type Result []dbapi.Post

type Service struct {
	db        *dbapi.DataBase
	uid       int
	selection map[*dbapi.LogicFunc]Result

	complexresult []dbapi.Post
}

func NewService(db *dbapi.DataBase, uid int) *Service {
	return &Service{
		db:        db,
		uid:       uid,
		selection: make(map[*dbapi.LogicFunc]Result),
	}
}

func (s *Service) RegisterLogicFunc(funcPointer *dbapi.LogicFunc) {
	s.selection[funcPointer] = Result{}
}

func (s *Service) GetFeed(UserID int, time int64, number int) ([]dbapi.Post, error) {
	var err error

	for f, result := range s.selection {
		s.selection[f], err = (*f)(s.uid, time, number)
		if err != nil {
			return nil, err
		}

		for _, val := range result {
			s.complexresult = append(s.complexresult, val)
		}
	}

	s.makeLimitedSelection(number)
	return s.complexresult, nil
}

func (s *Service) makeLimitedSelection(number int) {
	sort.SliceStable(s.complexresult, func(i, j int) bool {
		return s.complexresult[i].CreatedAt >= s.complexresult[j].CreatedAt
	})

	s.complexresult = s.complexresult[:number+1]
}

func (s *Service) Follow(TargetID int, CreatedAt int64) error { //t!
	err := s.db.InsertDataByFields("Followers",
		dbapi.Field{Name: "UserID", Value: s.uid},
		dbapi.Field{Name: "TargetID", Value: TargetID},
		dbapi.Field{Name: "CreatedAt", Value: CreatedAt})

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UnFollow(TargetID int) error {
	err := s.db.DeleteDataByFields("Followers",
		dbapi.Field{Name: "UserID", Value: s.uid},
		dbapi.Field{Name: "TargetID", Value: TargetID})

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) MakeComment(Message string, PostID int, CreatedAt int64) error {
	err := s.db.InsertDataByFields("MessageObjects",
		dbapi.Field{Name: "Type", Value: "comment"})
	if err != nil {
		return err
	}

	err = s.db.InsertDataByFields("Comments",
		dbapi.Field{Name: "ObjectID", Value: ""},
		dbapi.Field{Name: "PostID", Value: PostID},
		dbapi.Field{Name: "OwnerID", Value: s.uid},
		dbapi.Field{Name: "Message", Value: Message},
		dbapi.Field{Name: "CreatedAt", Value: CreatedAt})

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) MakePost(Message string, CreatedAt int64) error {
	err := s.db.InsertDataByFields("MessageObjects",
		dbapi.Field{Name: "Type", Value: "post"})
	if err != nil {
		return err
	}

	err = s.db.InsertDataByFields("Posts",
		dbapi.Field{Name: "ObjectID", Value: ""},
		dbapi.Field{Name: "OwnerID", Value: s.uid},
		dbapi.Field{Name: "Message", Value: Message},
		dbapi.Field{Name: "CreatedAt", Value: CreatedAt})

	if err != nil {
		return err
	}
	return nil
}
