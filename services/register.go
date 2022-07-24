package services

type Service interface{}

var (
	services []Service
)

func Register(service Service) {
	services = append(services, service)
}

func ListAll() []Service {
	return services
}
