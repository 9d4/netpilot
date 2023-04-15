package usermanager

type User struct {
	Id          string `json:".id"`
	Attributes  string `json:"attributes"`
	Disabled    bool   `json:"disabled,string"`
	Group       string `json:"group"`
	Name        string `json:"name"`
	OtpSecret   string `json:"otp-secret"`
	Password    string `json:"password"`
	SharedUsers int    `json:"shared-users,string"`
}
