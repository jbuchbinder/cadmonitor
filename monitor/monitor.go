package monitor

import (
	"sync"

	"github.com/pkg/errors"
)

// CadMonitor defines the interface for all monitors used to watch CAD systems.
type CadMonitor interface {
	// ConfigureFromValues populates fields specific to an implementation of
	// CadMonitor from a map[string]string.
	ConfigureFromValues(map[string]string) error
	// Login authenticates to a CAD system using the provided username and password
	Login(string, string) error
	// GetActiveCalls returns a list of active call URLs or identifiers
	GetActiveCalls() ([]string, error)
	GetStatus(string) (CallStatus, error)
	GetClearedCalls(string) (map[string]string, error)
	// SetDebug determines whether debug is enabled or not
	SetDebug(bool)
}

var (
	// ErrCadMonitorLoggedOut represents a status where the application needs to reauthenticate
	ErrCadMonitorLoggedOut = errors.New("logged out")

	cadMonitorRegistry     = map[string]func() CadMonitor{}
	cadMonitorRegistryLock = new(sync.Mutex)
)

// RegisterCadMonitor adds a new CadMonitor instance to the registry
func RegisterCadMonitor(name string, m func() CadMonitor) {
	cadMonitorRegistryLock.Lock()
	defer cadMonitorRegistryLock.Unlock()
	cadMonitorRegistry[name] = m
}

// GetCadMonitor instantiates a CadMonitor by name
func GetCadMonitor(name string) (m CadMonitor, err error) {
	var f func() CadMonitor
	var found bool
	if f, found = cadMonitorRegistry[name]; !found {
		err = errors.New("unable to locate cad monitor " + name)
		return
	}
	m = f()
	err = nil
	return
}
