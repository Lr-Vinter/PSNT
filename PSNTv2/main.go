package main

import ( //"database/sql"

	"feed/internal/dbapi"
	"fmt"
)

func main() {

	DB, _ := dbapi.InitializeDataBase("sqlite3", "test1.db")

	posts, err := DB.GetLikedByFollowerPosts(4, 0, 100, DB.ModLikesCount)
	fmt.Println(posts)
	fmt.Println(err)
	fmt.Println("=====")

	posts2, err2 := DB.GetFollowerPosts(3, 0, 20, DB.ModLikesCount)
	fmt.Println(posts2)
	fmt.Println(err2)

	// TESTS:

	//s, _ := DB.InsertDataByFields("Posts", dbapi.Field{Name: "Owner id", Value: 223}, dbapi.Field{Name: "PostID", Value: "544"})
	//fmt.Println(s)
	//S := service.NewService(DB, 3)
	//S.RegisterLogicFunc((*dbapi.LogicFunc)(&DB.GetFollowerPosts))
	//S.GetFeed(2, 2, 20)
	//posts, err := DB.GetFollowerPosts(3, 2, 20)
	//posts, err := DB.GetLikedByFollowerPosts(4, 100, DB.ModLikesCount)
	//posts, err := DB.GetFollowerPosts(3, 100, DB.ModExcludeSeenPosts, DB.ModExcludeFromUser, DB.ModLikesCount)
	//posts, err := DB.GetFollowerPosts(3, 100, DB.ModExcludeFromUser, DB.ModExcludeSeenPosts)
	//posts, err := DB.GetFollowerPosts(3, 100)
}
