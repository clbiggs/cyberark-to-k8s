package options

// ControllerOptions holds Configuration ControllerOptions that can be set by Command Line Flag, or Config File.
type ControllerOptions struct {
	Metrics    Metrics          `cfg:",squash" yaml:"metrics" json:"metrics"`
	CyberArk   CyberArk         `cfg:",squash" yaml:"cyberark" json:"cyberark"`
	KubeConfig KubernetesConfig `cfg:",squash" yaml:"kubeConfig" json:"kubeConfig"`
	Queues     Queues           `cfg:",squash" yaml:"queues" json:"queues"`

	Version       string `cfg:"version" flag:"version" yaml:"version" json:"version"`
	AllNamespaces bool   `cfg:"all_namespaces" flag:"all-namespaces" yaml:"allNamespaces" json:"allNamespaces"` // Indicates if all namespaces in the cluster should be handled by this instance.
	Namespaces    string `cfg:"namespaces" flag:"namespaces" yaml:"namespaces" json:"namespaces"`               // List (space delimited) of namespaces to monitor. If AllNamespaces is true this value is ignored.
	Labels        string `cfg:"labels" flag:"labels" yaml:"labels" json:"labels"`                               // Label Selector to apply for watching changes.

	KubeResyncPeriod     int `cfg:"kube_resync_period" flag:"kube-resync-period" yaml:"kubeResyncPeriod" json:"kubeResyncPeriod"`                 // Resync period for kubernetes changes, in seconds.
	CyberarkResyncPeriod int `cfg:"cyberark_resync_period" flag:"cyberark-resync-period" yaml:"cyberarkResyncPeriod" json:"cyberarkResyncPeriod"` // Resync period for cyberark changes, in seconds.
}

type Metrics struct {
	Enabled bool   `flag:"metrics-enabled" cfg:"metrics_enabled" yaml:"enabled" json:"enabled"`
	Port    string `flag:"metrics-port" cfg:"metrics_port" yaml:"port" json:"port"`
}

// CyberArk options. See CyberArk documention for proper values.
type CyberArk struct {
	Subdomain string `cfg:"cyberark_subdomain" flag:"cyberark-subdomain" yaml:"subdomain" json:"subdomain"`
	LogonType string `cfg:"cyberark_logontype" flag:"cyberark-logontype" yaml:"logonType" json:"logonType"`
	Username  string `cfg:"cyberark_username" flag:"cyberark-username" yaml:"userName" json:"userName"`
	Password  string `cfg:"cyberark_password" flag:"cyberark-password" yaml:"password" json:"password"`
}

type KubernetesConfig struct {
	Config    string `cfg:"kubernetes_config" flag:"kubernetes-config" yaml:"config" json:"config"`
	MasterURL string `cfg:"kubernetes_masterurl" flag:"kubernetes-masterurl" yaml:"masterURL" json:"masterURL"`
}

type Queues struct {
	DeleteMaxRequeues    int `cfg:"queues_delete_maxrequeues" flag:"queues_delete_maxrequeues" yaml:"deleteMaxRequeues" json:"deleteMaxRequeues"`
	DeleteThreads        int `cfg:"queues_delete_threads" flag:"queues_delete_threads" yaml:"deleteThreads" json:"deleteThreads"`
	CrdChangeMaxRequeues int `cfg:"queues_crdchange_maxrequeues" flag:"queues_crdchange_maxrequeues" yaml:"crdChangeMaxRequeues" json:"crdChangeMaxRequeues"`
	CrdChangeThreads     int `cfg:"queues_crdchange_threads" flag:"queues_crdchange_threads" yaml:"crdChangeThreads" json:"crdChangeThreads"`
	CyberArkMaxRequeues  int `cfg:"queues_cyberark_maxrequeues" flag:"queues_cyberark_maxrequeues" yaml:"cyberarkMaxRequeues" json:"cyberarkMaxRequeues"`
	CyberArkThreads      int `cfg:"queues_cyberark_threads" flag:"queues_cyberark_threads" yaml:"cyberarkThreads" json:"cyberarkThreads"`
}
