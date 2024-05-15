package p2p

type Adapter interface {
	Enabled() bool

	HasTarget() bool
	Listen(stop chan error)
	PushMsgForFiles(stop chan error)
	SetTarget(target string)
	Target() string
	Shutdown()
}
