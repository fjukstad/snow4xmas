package met

func GetAllLocations() ([]Data, error) {
	u := baseUrl + "/locations/v0.jsonld"
	return getData(u)
}
