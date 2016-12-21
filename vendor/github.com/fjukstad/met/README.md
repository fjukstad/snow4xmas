# met
Get data from [data.met.no](http://data.met.no). Note that both the MET API and
this package are still in development. 

# Setup
Before you can use the met package you'll have to request client credentials
from data.met.no [here](https://data.met.no/auth/requestCredentials.html), and
store them in an environment variable `CLIENT_ID`. 

# Example
Get monthly precipitation at climate station Blindern for 2016 (until Dec. 21).  

```go
package main

import (
	"fmt"

	"github.com/fjukstad/met"
)

func main() {

	f := met.Filter{
		Sources:       []string{"SN18700"}, // Oslo Blindern
		ReferenceTime: "2016-01-01T00:00:00.000Z/2016-12-21T00:00:00.000Z",
		Elements:      []string{"sum(precipitation_amount 1M)"},
	}

	data, err := met.GetObservations(f)
	if err != nil {
		fmt.Println(err)
        return
	}
	for _, d := range data {
		obs := d.Observations[0]
		fmt.Println(d.ReferenceTime.Format("Jan"), obs.Value, obs.Unit)
	}

}
```

```
Jan 43 mm
Feb 50.4 mm
Mar 47.1 mm
Apr 76.5 mm
May 73.5 mm
Jun 61.1 mm
Jul 85.3 mm
Aug 143.7 mm
Sep 41 mm
Oct 13.2 mm
Nov 74.2 mm
Dec 16.4 mm
```

# Acknowledgements 
Data from the [Norwegian Meteorological Institute](http://met.no). 
