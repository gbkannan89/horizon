package handler

// ServiceRegistry holds all application service interfaces for the API gateway.
// Each service is wired at startup by the main function.
type ServiceRegistry struct {
	UserSvc           interface{} // *user_service.UserService
	InstitutionSvc    interface{}
	FinancialEventSvc interface{}
	GoalSvc           interface{}
	AccountSvc        interface{}
	AllocationSvc     interface{}
	AssetSvc          interface{}
	LiabilitySvc      interface{}
	PortfolioSvc      interface{}
}

var svc *ServiceRegistry

// InitServices sets the global service registry used by all handlers.
func InitServices(registry *ServiceRegistry) {
	svc = registry
}
