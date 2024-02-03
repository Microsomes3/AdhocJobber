package servers

type Server interface {
	CreateServer() (bool, error)
	GetServer()
	DeleteServer()
	ExecuteCommandOnServer()
}
