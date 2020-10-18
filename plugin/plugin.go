package plugin

type Plugin interface {
	// ID returns the ID of the plugin.
	ID() string

	// IsLoaded returns whether the plugin is loaded.
	IsLoaded() bool

	// Load loads the plugin.
	Load() error

	// Unload unloads the plugin
	Unload() error
}
