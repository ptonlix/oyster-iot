package thirdpartaccess

type WxSession struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int64  `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

type WxToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type WxPhoneNum struct {
	Errcode   int64     `json:"errcode"`
	Errmsg    string    `json:"errmsg"`
	PhoneInfo PhoneInfo `json:"phone_info"`
}

type PhoneInfo struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string    `json:"countryCode"`
	Watermark       Watermark `json:"watermark"`
}

type Watermark struct {
	Timestamp int64  `json:"timestamp"`
	Appid     string `json:"appid"`
}
