package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"net/http"
	"strconv"
)

/*=======================================*/
// Method: NewFeatureHandler
// Description:
//		Adds a new feature to a layer
//		Saves layer to database
// @param apikey customer id
// @oaram ds datasource uuid
// @return json
/*=======================================*/
func NewFeatureHandler(w http.ResponseWriter, r *http.Request) {
	// Get request body
	// If this id done later in this function an EOF error occurs
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		Warning.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	/*=======================================*/

	// Get layer from database
	featCollection, err := DB.GetLayer(ds)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshal feature
	feat, err := geojson.UnmarshalFeature(body)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [500]")
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add new feature to layer
	featCollection.AddFeature(feat)
	DB.InsertLayer(ds, featCollection)

	// Generate message
	data := `{"status":"ok","datasource":"` + ds + `", "message":"feature added"}`
	js, err := json.Marshal(data)
	if err != nil {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update websockets
	conn := connection{ds: ds, ip: r.RemoteAddr}
	Hub.broadcast(true, &conn)

	// Return results
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	Info.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [200]")
	w.Write(js)

}

/*=======================================*/
// Method: ViewFeatureHandler
// Description:
//		Finds feature from layer
// @param apikey customer id
// @oaram ds datasource uuid
// @return feature geojson
/*=======================================*/
func ViewFeatureHandler(w http.ResponseWriter, r *http.Request) {

	// Get params
	apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	k, err := strconv.Atoi(vars["k"])
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [400]")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*=======================================*/
	// Check for apikey in request
	if apikey == "" {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		Warning.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !stringInSlice(ds, customer.Datasources) {
		Error.Println(r.RemoteAddr, "POST /api/v1/layer/"+ds+"/feature [401]")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	/*=======================================*/

	// Get layer from database
	data, err := DB.GetLayer(ds)
	if err != nil {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check for feature
	if k > len(data.Features) {
		Warning.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [404]")
		err := fmt.Errorf("Not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Marshal feature to json
	js, err := data.Features[k].MarshalJSON()
	if err != nil {
		Error.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//
	Info.Println(r.RemoteAddr, "GET /api/v1/layer/"+ds+"/feature/"+vars["k"]+" [200]")
	w.Write(js)

}
