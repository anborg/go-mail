package mail

//EmailInfo - convert eft into Email Data obj
type EmailInfo struct {
	From    string
	To      string
	Cc      string
	Subject string
	Body    string
}

// MailServerConfig smtp server details
type MailServerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	CcUser   string `yaml:"ccUser"`
	OpsUser  string `yaml:"opsUser"`
	Password string `yaml:"password"`
}
