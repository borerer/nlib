package signals

type SignalManager struct {
}

func NewSignalManager() *SignalManager {
	return &SignalManager{}
}

func (m *SignalManager) Start() error {
	return nil
}

func (m *SignalManager) Stop() error {
	return nil
}
