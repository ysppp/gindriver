package config

var Config = struct {
	AppName string

	DB struct {
		Name  string
		User  string
		Pass  string
		Host  string
		Port  uint
		Param string
	}

	ListenAddr    string
	RPID          string
	RPDisplayName string
}{}
