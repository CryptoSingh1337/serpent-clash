package system

type System interface {
	Name() string
	Update()
	Stop()
}
