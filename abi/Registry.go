package abi

type Registry interface {
	Auth() error
	SetToken(token string)
	Logout()
	Send(path string, inputData interface{}) (interface{}, error)
}

var defaultRegistry Registry = nil

func SetRegistry(registry Registry) {
	defaultRegistry = registry
}

func GetRegistry() Registry {
	return defaultRegistry
}
