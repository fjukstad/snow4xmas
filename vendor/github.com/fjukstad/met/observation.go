package met

func GetObservations(f Filter) ([]Data, error) {
	endpoint := baseUrl + "/observations/v0.jsonld?"
	u := createUrl(endpoint, f)
	return getData(u)
}
