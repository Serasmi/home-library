package user

type LoginUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id           string `json:"id" bson:"_id"`
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"passwordHash" bson:"passwordHash"`
}
