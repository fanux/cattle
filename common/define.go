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
	TaskTypeRemoveContainer
	TaskTypeStartContainer
	TaskTypeStopContainer
)

//ScaleItem is
type ScaleItem struct {
	Filters []string
	Number  int
	ENVs    []string
	Labels  map[string]string
}

//ScaleAPI is scale http api
type ScaleAPI struct {
	Items []ScaleItem
}

//ScaleConfig is ...
type ScaleConfig ScaleAPI

//Filter is parse from filter string
type Filter struct {
	Key      string
	Operater string
	Pattern  string
}
