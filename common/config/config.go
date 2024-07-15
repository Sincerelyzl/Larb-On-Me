package config

type MicroserviceConfig struct {
	AppPort       string
	AppDatabase   string
	AppCollection map[string]string
}

type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type LOMAuthConfig struct {
	Secret       string
	Argon2Config Argon2Config
}

type LOMServices struct {
	UserService     MicroserviceConfig
	ChatRoomService MicroserviceConfig
	GatewayService  MicroserviceConfig
}

type LOMConfig struct {
	LOMAuth     LOMAuthConfig
	LOMServices LOMServices
}

func Initialize() *LOMConfig {
	// Load configuration from file.

	var lomConfig LOMConfig
	return &lomConfig
}
