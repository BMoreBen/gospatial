package app

import (
	"encoding/json"
	// "github.com/gorilla/mux"
	"gospatial/utils"
	"net/http"
	"runtime"
	"time"
)

import mylogger "gospatial/logs"

var startTime = time.Now()

// SuperuserKey api servers superuser key
var SuperuserKey string = "su"

// PingHandler provides an api route for server health check
func PingHandler(w http.ResponseWriter, r *http.Request) {
	mylogger.Network.Debug("[In] ",r)
	data := `{"status": "success", "data": {"result": "pong"}}`
	js, err := json.Marshal(data)
	if err != nil {
		mylogger.Network.Critical(r.RemoteAddr, " GET /ping [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mylogger.Network.Info(r.RemoteAddr, " GET /ping [200]")
	w.Header().Set("Content-Type", "application/json")
	mylogger.Network.Debug("[Out] ",string(js))
	w.Write(js)
}

// ServerProfile returns basic server stats
func ServerProfile(w http.ResponseWriter, r *http.Request) {
	mylogger.Network.Debug("[In] ",r)
	var data map[string]interface{}
	data = make(map[string]interface{})
	data["registered"] = startTime.UTC()
	data["uptime"] = time.Since(startTime).Seconds()
	// data["status"] = AppMode // debug, static, standard
	data["num_cores"] = runtime.NumCPU()
	// data["free_mem"] = runtime.MemStats()
	js, err := json.Marshal(data)
	if err != nil {
		mylogger.Network.Critical(r.RemoteAddr, " GET /management/profile [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mylogger.Network.Info(r.RemoteAddr, " GET /management/profile [200]")
	w.Header().Set("Content-Type", "application/json")
	mylogger.Network.Debug("[Out] ",string(js))
	w.Write(js)
}

// NewCustomerHandler superuser route to create new api customers/apikeys
func NewCustomerHandler(w http.ResponseWriter, r *http.Request) {
	mylogger.Network.Debug("[In] ",r)
	// Check auth key
	if SuperuserKey != r.FormValue("authkey") {
		mylogger.Network.Error(r.RemoteAddr, " POST /management/customer [401]")
		http.Error(w, `{"status": "fail", "data": {"error": "unauthorized"}}`, http.StatusUnauthorized)
		return
	}
	// new customer
	apikey := utils.NewAPIKey(12)
	customer := Customer{Apikey: apikey}
	err := DB.InsertCustomer(customer)
	if err != nil {
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return results
	data := `{"status":"success","apikey":"` + apikey + `", "result":"customer created"}`
	js, err := json.Marshal(data)
	if err != nil {
		mylogger.Network.Critical(r.RemoteAddr, " POST /management/customer [500]")
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mylogger.Network.Info(r.RemoteAddr, " POST /management/customer [200]")
	mylogger.Network.Debug("[Out] ",string(js))
	w.Write(js)
}

// ShareLayerHandler gives customer access to an existing datasource.
// @param apikey - customer to give access
// @param authkey
// @return json
// func ShareLayerHandler(w http.ResponseWriter, r *http.Request) {

// 	// Get url params
// 	apikey := r.FormValue("apikey")
// 	authkey := r.FormValue("authkey")

// 	// Get ds from url path
// 	vars := mux.Vars(r)
// 	ds := vars["ds"]

// 	// superuser access
// 	if SuperuserKey != authkey {
// 		http.Error(w, "unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	if apikey == "" {
// 		networkLoggerError.Println(r.RemoteAddr, "PUT /api/v1/layer/{ds} [401]")
// 		http.Error(w, "bad request", http.StatusBadRequest)
// 		return
// 	}

// 	// Get customer from database
// 	customer, err := DB.GetCustomer(apikey)
// 	if err != nil {
// 		networkLoggerWarning.Println(r.RemoteAddr, "PUT /api/v1/layer/{ds} [404]")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	// Add datasource uuid to customer
// 	customer.Datasources = append(customer.Datasources, ds)
// 	DB.InsertCustomer(customer)

// 	// Generate message
// 	data := `{"status":"ok","datasource":"` + ds + `"}`
// 	js, err := json.Marshal(data)
// 	if err != nil {
// 		networkLoggerError.Println(r.RemoteAddr, "PUT /api/v1/layer [500]")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Return results
// 	networkLoggerInfo.Println(r.RemoteAddr, "PUT /api/v1/layer [200]")
// 	w.Header().Set("Content-Type", "application/json")
// 	// allow cross domain AJAX requests
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Write(js)

// }
