package module

type Module interface {
	// Name returns the name of the module.
	Name() string

	// Version returns the version of the module.
	Version() string
}
