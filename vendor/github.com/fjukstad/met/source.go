package met

func GetAllSources() ([]Data, error) {
	return GetSources(Filter{})
}

func GetSources(f Filter) ([]Data, error) {
	endpoint := baseUrl + "/sources/v0.jsonld?"
	u := createUrl(endpoint, f)
	return getData(u)
}
