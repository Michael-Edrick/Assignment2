package DataUser

import (
	"time"
)

type User struct{
	Id       int
	Username string
	Email    string
	Password string
	Age      int
	Created_at time.Time
	Updated_at time.Time
}

type Photo struct{
	Id	uint
	Title string
	Caption string 
	Photo_url string
	User_id uint
	Created_at time.Time
	Updated_at time.Time
}

type Comment struct{
	Id uint
	User_id uint
	Photo_id uint
	Message string
	Created_at time.Time
	Updated_at time.Time
}

type SocialMedia struct{
	Id uint
	Name string
	Social_media_url string
	User_id uint
}