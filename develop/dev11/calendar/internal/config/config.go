package config

type Config struct {
	StorageFilePath string
	Host            string
	Port            string
}

func GetConfig() Config {
	return Config{
		StorageFilePath: "storage/data.json",
		Host:            "127.0.0.1",
		Port:            "8081",
	}
}
