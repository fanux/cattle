package common

// special labels define
const (
	LabelKeyNamespace = "namespace"
	LabelKeyService   = "service"
	LabelKeyApp       = "app"
	LabelKeyName      = "name"
)

// special Environment define
const (
	EnvironmentPriority  = "PRIORITY"
	EnvironmentMinNumber = "MIN_NUMBER"
)

// tasks
const (
	TaskTypeCreateContainer = iota
	TaskTypeDestroyContainer
	TaskTypeStartContainer
	TaskTypeStopContainer
)

// TASK_TYPE
const (
	EnvTaskTypeKey = "TASK_TYPE"
	EnvTaskStart   = "start"
	EnvTaskStop    = "stop"
	EnvTaskCreate  = "create"
	EnvTaskDestroy = "destroy"
)

// affinity and constaint
const (
	Affinity   = "affinity"
	Constraint = "constraint"
	Applots    = "applots"
)

// default app lots
const (
	DefaultAppLots  = 1
	DefaultPriority = -1 // prevent seize all containers in cluster, when ...
)

// ScaleItem is
type ScaleItem struct {
	Filters []string
	Number  int
	ENVs    []string
	Labels  map[string]string
}

// ScaleAPI is scale http api
type ScaleAPI struct {
	Items []ScaleItem
}

// ScaleConfig is ...
type ScaleConfig ScaleAPI

// Filter is parse from filter string
type Filter struct {
	Key      string
	Operater string
	Pattern  string
}
