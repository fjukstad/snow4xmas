package met

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Context          string    `json:"@context"`
	Type             string    `json:"@type"`
	APIVersion       string    `json:"apiVersion"`
	License          string    `json:"license"`
	CreatedAt        time.Time `json:"createdAt"`
	QueryTime        float64   `json:"queryTime"`
	CurrentItemCount int       `json:"currentItemCount"`
	ItemsPerPage     int       `json:"itemsPerPage"`
	Offset           int       `json:"offset"`
	TotalItemCount   int       `json:"totalItemCount"`
	NextLink         string    `json:"nextLink"`
	PreviousLink     string    `json:"previousLink"`
	CurrentLink      string    `json:"currentLink"`
	Data             []Data    `json:"data"`
}

type Data struct {
	Id                    string `json:"id"`
	Name                  string `json:"name"`
	Country               string `json:"country"`
	SourceID              string `json:"sourceId"`
	Geometry              `json:"geometry"`
	Levels                `json:"levels"`
	ReferenceTime         time.Time `json:"referenceTime"`
	Observations          `json:"observations"`
	ValidFrom             string `json:"validFrom"`
	LegacyMetNoConvention `json:"legacyMetNoConvention"`
	CfConvention          `json:"cfConvention"`
}

type CfConvention struct {
	StandardName string `json:"standardName"`
	Unit         string `json:"unit"`
	Status       string `json:"status"`
}

type LegacyMetNoConvention struct {
	ElemCodes []string `json:"elemCodes"`
	Category  string   `json:"category"`
	Unit      string   `json:"unit"`
}

type Level struct {
	LevelType string `json:"levelType"`
	Value     int    `json:"value"`
	Unit      string `json:"unit"`
}

type Observation struct {
	ElementId           string  `json:"elementId"`
	Value               float64 `json:"value"`
	Unit                string  `json:"unit"`
	CodeTable           string  `json:"codeTable"`
	PerformanceCategory string  `json:"performanceCategory"`
	ExposureCategory    string  `json:"exposureCategory"`
	QualityCode         int     `json:"qualityCode"`
	DataVersion         string  `json:"dataVersion"`
}

type Observations []Observation

type Geometry struct {
	Type         string    `json:"@type"`
	Coordinates  []float64 `json:"coordinates"`
	Interpolated bool      `json:"interpolated"`
}

// Use Filter to set different query parameters. See https://data.met.no/docs
// for more information. Some queries require some of the parameters set.
type Filter struct {
	Sources               []string
	ReferenceTime         string
	Elements              []string
	PerformanceCategories []string
	ExposureCategories    []string
	Fields                []string
	Ids                   []string
	Types                 []string
	Geometry              string
	ValidTime             string
}

func (l Level) String() string {
	return "LevelType: " + l.LevelType + " Value: " + strconv.Itoa(l.Value) + " Unit: " + l.Unit
}

type Levels []Level

func (ls Levels) String() string {
	return sliceToString("\n", ls)
}

func (g *Geometry) String() string {
	lat := strconv.FormatFloat(g.Coordinates[0], 'f', -1, 64)
	long := strconv.FormatFloat(g.Coordinates[1], 'f', -1, 64)
	return "Type: " + g.Type + "Coordinates: " + lat + "," + long
}

func (o Observation) String() string {

	v := strconv.FormatFloat(o.Value, 'f', -1, 64)
	qc := strconv.Itoa(o.QualityCode)

	return "ElementId: " + o.ElementId +
		" Value: " + v +
		" Unit: " + o.Unit +
		" CodeTable: " + o.CodeTable +
		" Performance Category: " + o.PerformanceCategory +
		" Quality Code: " + qc +
		" DataVersion:" + o.DataVersion

}

func (obs Observations) String() string {
	return sliceToString("\n", obs)
}

func sliceToString(sep string, t interface{}) string {
	buf := new(bytes.Buffer)
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)
		for i := 0; i < s.Len(); i++ {
			if i > 0 {
				buf.WriteString(sep)
			}
			item := s.Index(i).Interface()
			fmt.Fprint(buf, item)
		}
	}
	return buf.String()
}

func (lmnc *LegacyMetNoConvention) String() string {
	return "ElemCodes" + sliceToString(",", lmnc.ElemCodes) + "Category: " + lmnc.Category + "Unit: " + lmnc.Unit
}

func (cc *CfConvention) String() string {
	return "StandardName: " + cc.StandardName + " Unit: " + cc.Unit +
		" Status:" + cc.Status
}

var baseUrl = "https://data.met.no"

func getClientId() (string, error) {
	id := os.Getenv("CLIENT_ID")
	if id == "" {
		return "", errors.New("CLIENT_ID not set")
	}
	return id, nil
}

func get(endpoint string) ([]byte, error) {
	id, err := getClientId()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(id, "")

	c := http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 400 {
		return nil, errors.New("Invalid parameter value or malformed request.")
	}

	if resp.StatusCode == 401 {
		return nil, errors.New("Unauthorized client ID.")
	}

	if resp.StatusCode == 404 {
		return nil, errors.New("No data was found for the list of query Ids.")
	}

	if resp.StatusCode == 400 {
		return nil, errors.New("Internal server error.")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func createUrl(endpoint string, f Filter) string {

	sources := strings.Join(f.Sources, ",")
	elements := strings.Join(f.Elements, ",")
	performanceCategories := strings.Join(f.PerformanceCategories, ",")
	exposureCategories := strings.Join(f.ExposureCategories, ",")
	fields := strings.Join(f.Fields, ",")
	ids := strings.Join(f.Fields, ",")
	types := strings.Join(f.Types, ",")

	v := url.Values{}

	if len(sources) > 0 {
		v.Add("sources", sources)
	}

	if len(elements) > 0 {
		v.Add("elements", elements)
	}

	if f.ReferenceTime != "" {
		v.Add("referencetime", f.ReferenceTime)
	}

	if len(performanceCategories) > 0 {
		v.Add("performancecategory", performanceCategories)
	}

	if len(exposureCategories) > 0 {
		v.Add("exposurecategory", exposureCategories)
	}

	if len(fields) > 0 {
		v.Add("fields", fields)
	}

	if len(ids) > 0 {
		v.Add("ids", ids)
	}

	if len(types) > 0 {
		v.Add("types", types)
	}

	if f.Geometry != "" {
		v.Add("geometry", f.Geometry)
	}

	if f.ValidTime != "" {
		v.Add("validtime", f.ValidTime)
	}

	return endpoint + v.Encode()
}

func getData(u string) ([]Data, error) {
	body, err := get(u)
	if err != nil {
		return []Data{}, err
	}
	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return []Data{}, err
	}

	return response.Data, nil

}

func (d *Data) String() string {
	str := "Id: " + d.Id +
		"\nName: " + d.Name +
		"\nCountry: " + d.Country +
		"\nSourceId: " + d.SourceID +
		"\nGeometry: " + d.Geometry.String() +
		"\nLevels: " + d.Levels.String() +
		"\nReferenceTime: " + d.ReferenceTime.String() +
		"\nObservations: " + d.Observations.String() +
		"\nValidFrom: " + d.ValidFrom +
		"\nLegacyMetNoConvention: " + d.LegacyMetNoConvention.String() +
		"\nCfConvention: " + d.CfConvention.String()

	return str
}
