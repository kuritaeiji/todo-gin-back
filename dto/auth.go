package dto

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Oauth struct {
	Code              string `json:"code"`
	State             string `json:"state"`
	LocalStorageState string `json:"local_storage_state"`
}
