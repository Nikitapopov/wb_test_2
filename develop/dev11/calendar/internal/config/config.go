package config

// Конфигурация приложения
type Config struct {
	StorageFilePath string
	Host            string
	Port            string
}

// Геттер конфигурации
func GetConfig() Config {
	return Config{
		StorageFilePath: "storage/data.json",
		Host:            "127.0.0.1",
		Port:            "8081",
	}
}
