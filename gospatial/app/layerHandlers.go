package app

import (
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	"gospatial/utils"
	"net/http"
)

// ViewLayersHandler returns json containing customer layers
// @param apikey customer id
// @return json
func ViewLayersHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)
	// Get params
	//apikey := r.FormValue("apikey")
	// Check for apikey in request
	// if apikey == "" {
	// 	NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layers [401]")
	// 	http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
	// 	return
	// }

	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layers [401]")
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layers [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// return results
	js, err := json.Marshal(customer)
	if err != nil {
		NetworkLogger.Critical(r.RemoteAddr, " POST /api/v1/layers [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	NetworkLogger.Info(r.RemoteAddr, " POST /api/v1/layers [200]")
	NetworkLogger.Debug("[Out] ", string(js))
	SendJsonResponse(w, js)
}

// NewLayerHandler creates a new geojson layer. Saves layer to database and adds layer to customer
// @param apikey
// @return json
func NewLayerHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// // Get params
	// apikey := r.FormValue("apikey")

	// // Check for apikey in request
	// if apikey == "" {
	// 	NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer [401]")
	// 	http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
	// 	return
	// }

	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer [401]")
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create datasource
	ds, err := DB.NewLayer()
	if err != nil {
		NetworkLogger.Critical(r.RemoteAddr, " POST /api/v1/layer [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add datasource uuid to customer
	customer.Datasources = append(customer.Datasources, ds)
	DB.InsertCustomer(customer)

	// Generate message
	data := `{"status":"success","datasource":"` + ds + `"}`
	js, err := json.Marshal(data)
	if err != nil {
		NetworkLogger.Critical(r.RemoteAddr, " POST /api/v1/layer [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return results
	NetworkLogger.Error(r.RemoteAddr, " POST /api/v1/layer [200]")
	NetworkLogger.Debug("[Out] ", string(js))
	SendJsonResponse(w, js)
}

// ViewLayerHandler returns geojson of requested layer. Apikey/customer is checked for permissions to requested layer.
// @param ds
// @param apikey
// @return geojson
func ViewLayerHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// Get params
	//apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	// Check for apikey in request
	// if apikey == "" {
	// 	NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+" [401]")
	// 	http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
	// 	return
	// }

	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+" [401]")
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !utils.StringInSlice(ds, customer.Datasources) {
		NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+" [401]")
		http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
		return
	}

	// Get layer from database
	lyr, err := DB.GetLayer(ds)
	if err != nil {
		NetworkLogger.Error(r.RemoteAddr, " GET /api/v1/layer/"+ds+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Marshal datasource layer to json
	js, err := lyr.MarshalJSON()
	if err != nil {
		NetworkLogger.Critical(r.RemoteAddr, " GET /api/v1/layer/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return layer json
	NetworkLogger.Info(r.RemoteAddr, " GET /api/v1/layer/"+ds+" [200]")
	NetworkLogger.Debug("[Out] ", string(js))
	SendJsonResponse(w, js)
}

// DeleteLayerHandler deletes layer from database and removes it from customer list.
// @param ds
// @param apikey
// @return json
func DeleteLayerHandler(w http.ResponseWriter, r *http.Request) {
	NetworkLogger.Debug("[In] ", r)

	// Get params
	//apikey := r.FormValue("apikey")

	// Get ds from url path
	vars := mux.Vars(r)
	ds := vars["ds"]

	// Check for apikey in request
	// if apikey == "" {
	// 	NetworkLogger.Error(r.RemoteAddr, " DELETE /api/v1/layer/"+ds+" [401]")
	// 	http.Error(w, `{"status": "error", "result": "unauthorized"}`, http.StatusUnauthorized)
	// 	return
	// }

	apikey := GetApikeyFromRequest(w, r)
	if apikey == "" {
		NetworkLogger.Error(r.RemoteAddr, " DELETE /api/v1/layer/"+ds+" [401]")
		return
	}

	// Get customer from database
	customer, err := DB.GetCustomer(apikey)
	if err != nil {
		NetworkLogger.Error(r.RemoteAddr, " DELETE /api/v1/layer/"+ds+" [404]")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check customer datasource list
	if !utils.StringInSlice(ds, customer.Datasources) {
		NetworkLogger.Error(r.RemoteAddr, " DELETE /api/v1/layer/"+ds+" [401]")
		http.Error(w, `{"status": "error", "result": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// KEEP LAYER IN CASE OF RECOVERY
	// Delete layer from database
	// err = DB.DeleteLayer(ds)
	// if err != nil {
	// 	networkLoggerInfo.Println(r.RemoteAddr, "DELETE /api/v1/layer/"+ds+" [500]")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Delete layer from customer
	i := utils.SliceIndex(ds, customer.Datasources)
	customer.Datasources = append(customer.Datasources[:i], customer.Datasources[i+1:]...)
	DB.InsertCustomer(customer)

	// Generate message
	data := `{"status":"ok","datasource":"` + ds + `", "result":"datasource deleted"}`
	js, err := json.Marshal(data)
	if err != nil {
		NetworkLogger.Critical(r.RemoteAddr, " DELETE /api/v1/layer/"+ds+" [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Returns results
	NetworkLogger.Info(r.RemoteAddr, " DELETE /api/v1/layer/"+ds+" [200]")
	NetworkLogger.Debug("[Out] ", string(js))
	SendJsonResponse(w, js)
}
