package plugin

type Module struct{}

func (m *Module) Name() string {
	return "Plugin Interface"
}

func (m *Module) Version() string {
	return "0.1.0"
}
