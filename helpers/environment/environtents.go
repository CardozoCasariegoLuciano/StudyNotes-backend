package environment

const (
	DEV  = "develop"
	TEST = "test"
	PROD = "production"
)

type env struct {
	environment string
}

var appEnv = func() *env {
	return &env{environment: DEV}
}()

// Getter
func GetEnvirontment() string {
	return appEnv.environment
}

// Setters
func SetTestEnvirontment() {
	appEnv.environment = TEST
}

func SetDevEnvirontment() {
	appEnv.environment = DEV
}

func SetProdEnvirontment() {
	appEnv.environment = PROD
}
