package models

type Agent struct {
	UUID     string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid"`
	Name     string `json:"name"`
	Username string `json:"username"`
	IP       string `json:"ip"`
	Created  string `json:"created"`
	Checkin  string `json:"checkin"`
	Token    string `json:"token"`
}

type Task struct {
	ID      string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AgentID string `gorm:"type:uuid;" json:"agentid"`
	Start   string `json:"start"`
	Command string `json:"command"`
	Result  string `json:"result"`
	End     string `json:"end"`
}

type Response struct {
	Type    string `json:"type"`
	Data    string `json:"data"`
	Message string `json:"message"`
}
