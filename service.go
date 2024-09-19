package service_init_tool

type Service struct {
	name string
}

func newService(name string) *Service {
	return &Service{
		name: name,
	}
}

func (svc *Service) GetName() {

}
