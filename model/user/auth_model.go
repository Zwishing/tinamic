package user

type SignIn struct {
	LoginAccount string          `json:"loginAccount" xml:"loginAccount" form:"loginAccount" :"loginAccount"`
	Password     string          `json:"password" xml:"password" form:"password" :"password"`
	Category     AccountCategory `json:"category" xml:"category" :"category"`
}
