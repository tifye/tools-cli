package security

type UserProfile struct {
	Developer   bool              `json:"developer"`
	GlobalAdmin bool              `json:"global_admin"`
	Name        string            `json:"name"`
	Roles       map[string]string `json:"roles"`
	Profile     struct {
		Firstname string      `json:"firstname"`
		Lastname  string      `json:"lastname"`
		Email     string      `json:"email"`
		Location  interface{} `json:"location"`
		Fullname  string      `json:"fullname"`
		Human     bool        `json:"human"`
	} `json:"profile"`
	Auth struct {
		SsoID        string `json:"sso_id"`
		ClientKey    string `json:"client_key"`
		ClientSecret string `json:"client_secret"`
	} `json:"auth"`
	Issued       int64  `json:"issued"`
	Expires      int64  `json:"expires"`
	APIKey       string `json:"api_key"`
	ID           string `json:"id"`
	SystemTest   bool   `json:"system_test"`
	BetaTest     bool   `json:"beta_test"`
	InternalTest bool   `json:"internal_test"`
	AccessToken  string `json:"access_token"`
}
