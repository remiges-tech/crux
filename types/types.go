package types

type AppConfig struct {
	DBHost        string   `json:"db_host"`
	DBPort        string   `json:"db_port"`
	DBUser        string   `json:"user"`
	DBPassword    string   `json:"db_password"`
	DBName        string   `json:"db_name"`
	AppServerPort string   `json:"app_server_port"`
	EtcdEndPoints []string `json:"etcd_endpoints"`
	ErrorTypeFile string   `json:"error_type_file"`
}
