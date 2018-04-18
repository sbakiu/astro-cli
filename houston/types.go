package houston

type CreateTokenResponse struct {
	Data struct {
		CreateToken Token `json:"createToken"`
	} `json:"data"`
}

type Token struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Token   string  `json:"token"`
	Decoded Decoded `json:"decoded"`
}

type Deployment struct {
	Id          string `json:"uuid"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	ReleaseName string `json:"release_name"`
	Version     string `json:"version"`
}

type Decoded struct {
	ID  string `json:"id"`
	SU  bool   `json:"sU"`
	Iat int    `json:"iat"`
	Exp int    `json:"exp"`
}
