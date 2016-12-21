package common

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
