package common

// special labels define
const (
	//LabelKeyNamespace is
	LabelKeyNamespace = "namespace"
	//LabelKeyService is
	LabelKeyService = "service"
	//LabelKeyApp is
	LabelKeyApp = "app"
)

// special Environment define
const (
	//EnvironmentPriority is
	EnvironmentPriority = "PRIORITY"
	//EnvironmentMinNumber is
	EnvironmentMinNumber = "MIN_NUMBER"
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
