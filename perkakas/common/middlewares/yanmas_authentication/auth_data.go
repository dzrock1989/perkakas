package yanmas_authentication

type Claims struct {
	UserUUID string `json:"user_uuid"`
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	Created  int64  `json:"created"`
}
