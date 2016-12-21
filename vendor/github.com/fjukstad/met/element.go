package met

func GetAllElements() ([]Data, error) {
	u := baseUrl + "/elements/v0.jsonld"
	return getData(u)
}
