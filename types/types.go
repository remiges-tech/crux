package types

type AppConfig struct {
	DBHost        string `json:"db_host"`
	DBPort        int    `json:"db_port"`
	DBUser        string `json:"db_user"`
	DBPassword    string `json:"db_password"`
	DBName        string `json:"db_name"`
	DriverName    string `json:"driver_name"`
	AppServerPort int    `json:"app_server_port"`
	ErrorTypeFile string `json:"error_type_file"`
}
