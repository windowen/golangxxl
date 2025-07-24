package rpcclient

type ServiceClients struct {
	// 直播服务
	LiveClient    *LiveClient
	FinanceClient *FinanceClient
}

var (
	ServiceClientsInstance *ServiceClients
)

// NewServiceClients 创建 ServiceClients
func NewServiceClients() *ServiceClients {
	ServiceClientsInstance = &ServiceClients{}

	return ServiceClientsInstance
}

func (s *ServiceClients) Start() {
	s.LiveClient = newLiveClient()
	s.FinanceClient = newFinanceClient()
}

func (s *ServiceClients) Stop() {

}
