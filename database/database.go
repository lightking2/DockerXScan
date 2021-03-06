package database

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrBackendException is an error that occurs when the database backend does
	// not work properly (ie. unreachable).
	ErrBackendException = errors.New("database: an error occured when querying the backend")

	// ErrInconsistent is an error that occurs when a database consistency check
	// fails (i.e. when an entity which is supposed to be unique is detected
	// twice)
	ErrInconsistent = errors.New("database: inconsistent database")
)

// RegistrableComponentConfig is a configuration block that can be used to
// determine which registrable component should be initialized and pass custom
// configuration to it.
type RegistrableComponentConfig struct {
	Type    string
	Options map[string]interface{}
}

var drivers = make(map[string]Driver)

// Driver is a function that opens a Datastore specified by its database driver type and specific
// configuration.
type Driver func(RegistrableComponentConfig) (Datastore, error)

// Register makes a Constructor available by the provided name.
//
// If this function is called twice with the same name or if the Constructor is
// nil, it panics.
func Register(name string, driver Driver) {
	if driver == nil {
		panic("database: could not register nil Driver")
	}
	if _, dup := drivers[name]; dup {
		panic("database: could not register duplicate Driver: " + name)
	}
	drivers[name] = driver
}

// Open opens a Datastore specified by a configuration.
func Open(cfg RegistrableComponentConfig) (Datastore, error) {
	driver, ok := drivers[cfg.Type]
	if !ok {
		return nil, fmt.Errorf("database: unknown Driver %q (forgotten configuration or import?)", cfg.Type)
	}
	return driver(cfg)
}

// Datastore represents the required operations on a persistent data store for
// a Clair deployment.
type Datastore interface {
	//关闭数据
	Close()
	//插入namespace
	InsertNamespace(namespace Namespace) (int ,error)

	//查询namespace
	ListNamespaces() ([]Namespace, error)

	//插入layer
	InsertLayer(Layer) error

	//查询layer
	FindLayer(name string, withFeatures, withVulnerabilities bool) (Layer, error)

	//删除layer
	DeleteLayer(name string) error

	//插入特征
	InsertFeature(feature Feature) (int, error)

	//插入特征版本
	InsertFeatureVersion(fv FeatureVersion) (id int, err error)

	//列出漏洞
	ListVulnerabilities(namespaceName string, limit int, page int) ([]Vulnerability, int, error)

	//插入漏洞
	InsertVulnerabilities(vulnerabilities []Vulnerability, createNotification bool) error

	//查找漏洞
	FindVulnerability(namespaceName, name string) (Vulnerability, error)

	//插入漏洞修复
	InsertVulnerabilityFixes(vulnerabilityNamespace, vulnerabilityName string, fixes []FeatureVersion) error

	//删除漏洞
	DeleteVulnerability(namespaceName, name string) error

	//删除漏洞修复
	DeleteVulnerabilityFix(vulnerabilityNamespace, vulnerabilityName, featureName string) error

	InsertKeyValue(key, value string) error

	GetKeyValue(key string) (string, error)

	Lock(name string, owner string, duration time.Duration, renew bool) (bool, time.Time)

	Unlock(name, owner string)

	FindLock(name string) (string, time.Time, error)

	GetAvailableNotification(renotifyInterval time.Duration) (VulnerabilityNotification, error)

	GetNotification(name string, limit int, page VulnerabilityNotificationPageNumber) (VulnerabilityNotification, VulnerabilityNotificationPageNumber, error)

	SetNotificationNotified(name string) error

	DeleteNotification(name string) error

	Ping() bool
}