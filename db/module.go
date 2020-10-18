package db

type Module struct{}

func (m *Module) Name() string {
	return "DB"
}

func (m *Module) Version() string {
	return "0.1.0"
}
