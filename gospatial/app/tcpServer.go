package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paulmach/go.geojson"
	"gospatial/utils"
	"io/ioutil"
	"net"
	"net/textproto"
	"os"
	"os/exec"
	"strings"
)

const (
	TCP_DEFAULT_CONN_HOST = "localhost"
	TCP_DEFAULT_CONN_PORT = "3333"
	TCP_DEFAULT_CONN_TYPE = "tcp"
)

type TcpServer struct {
	Host string
	Port string
}

func (self TcpServer) Start() {
	go func() {
		// Check settings and apply defaults
		host := self.Host
		if host == "" {
			host = TCP_DEFAULT_CONN_HOST
		}

		port := self.Port
		if port == "" {
			port = TCP_DEFAULT_CONN_PORT
		}

		// Listen for incoming connections.
		l, err := net.Listen(TCP_DEFAULT_CONN_TYPE, host+":"+port)
		if err != nil {
			ServerLogger.Error("Error listening:", err.Error())
			panic(err)
		}

		// Close the listener when the application closes.
		defer l.Close()

		ServerLogger.Info("Tcp Listening on " + host + ":" + port)

		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				NetworkLogger.Error("Error accepting: ", err.Error())
				return
			}

			NetworkLogger.Info("Connection open ", conn.RemoteAddr().String(), " [TCP]")

			// check for local connection
			if strings.Contains(conn.RemoteAddr().String(), "127.0.0.1") {
				// Handle connections in a new goroutine.
				go self.tcpClientHandler(conn)
			} else {
				conn.Close()
			}

		}
	}()
}

// Handles incoming requests.
func (self TcpServer) tcpClientHandler(conn net.Conn) {

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	defer conn.Close()

	// DEBUGGING
	//	authenticated := false
	authenticated := true

	for {

		// will listen for message to process ending in newline (\n)
		//message, _ := bufio.NewReader(conn).ReadString('\n') // sometimes read partial messages
		message, _ := tp.ReadLine()

		// output message received
		NetworkLogger.Info("[TCP] Message Received: ", string([]byte(message)))

		// json parse message
		req := TcpMessage{}
		err := json.Unmarshal([]byte(message), &req)
		if err != nil {

			// invalid message
			// close connection
			NetworkLogger.Warn("error:", err)
			resp := `{"status": "error", "error": "` + fmt.Sprintf("%v", err) + `",""}`
			conn.Write([]byte(resp + "\n"))
			NetworkLogger.Info("Connection closed", " [TCP]")
			return

		} else {

			success := false
			switch {

			case req.Method == "ping":
				resp := `{"status": "ok", "data": { "message": "pong", "version": "` + VERSION + `"}}`
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "help":
				conn.Write([]byte("Methods:\n"))
				conn.Write([]byte("\t ping\n"))
				conn.Write([]byte("\t assign_datasource\n"))
				conn.Write([]byte("\t create_apikey\n"))
				conn.Write([]byte("\t insert_apikey\n"))
				conn.Write([]byte("\t insert_feature\n"))
				conn.Write([]byte("\t edit_feature\n"))
				conn.Write([]byte("\t create_datasource\n"))
				conn.Write([]byte("\t export_apikeys\n"))
				conn.Write([]byte("\t export_apikey\n"))
				conn.Write([]byte("\t export_datasources\n"))
				conn.Write([]byte("\t export_datasource\n"))
				conn.Write([]byte("\t import_file\n"))
				success = true

			case req.Method == "authenticate":
				// {"method":"authenticate", "authkey": "7q1qcqmsxnvw"}
				authenticated = SuperuserKey == req.Authkey
				if authenticated {
					resp := `{"status": "ok", "data": {}}`
					conn.Write([]byte(resp + "\n"))
				} else {
					NetworkLogger.Warn("error: incorrect authkey", " [TCP]")
					resp := `{"status": "error", "error": "incorrect authkey"}`
					conn.Write([]byte(resp + "\n"))
				}
				success = true

			case req.Method == "assign_datasource" && authenticated:
				// {"method":"assign_datasource"}
				resp := `{"status": "ok", "data": {}}`
				datasource_id := req.Datasource
				apikey := req.Apikey

				if "" == datasource_id || "" == apikey {
					err := errors.New("Missing required parameters")
					resp = `{"status": "error", "error": "` + err.Error() + `"}`
				} else {

					customer, err := DB.GetCustomer(apikey)
					resp = `{"status": "ok", "data": {}}`
					if err != nil {
						resp = `{"status": "error", "error": "` + err.Error() + `"}`
					}

					_, err = DB.GetLayer(datasource_id)
					if err != nil {
						resp = `{"status": "error", "error": "` + err.Error() + `"}`
					} else {
						customer.Datasources = append(customer.Datasources, datasource_id)
						DB.InsertCustomer(customer)
					}
				}

				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "create_apikey" && authenticated:
				// {"method":"create_apikey"}
				apikey := utils.NewAPIKey(12)
				customer := Customer{Apikey: apikey}
				resp := `{"status": "ok", "data": {"apikey": "` + apikey + `"}}`
				err := DB.InsertCustomer(customer)
				if err != nil {
					fmt.Println(err)
					resp = `{"status": "error", "error": "` + err.Error() + `"}`
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "insert_apikey" && authenticated:
				// {"method": "insert_apikey"}
				resp := `{"status": "ok", "data": {}}`
				if "" == req.Data.Apikey {
					err := errors.New("Missing required parameters")
					resp = `{"status": "error", "error": "` + err.Error() + `"}`
				} else {

					customer := Customer{Apikey: req.Data.Apikey, Datasources: req.Data.Datasources}
					resp = `{"status": "ok", "data": {"apikey": "` + req.Data.Apikey + `"}}`
					err := DB.InsertCustomer(customer)
					if err != nil {
						fmt.Println(err)
						resp = `{"status": "error", "error": "` + err.Error() + `"}`
					}
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "insert_feature" && authenticated:
				// {"method":"insert_feature"}
				resp := `{"status":"ok","data": {"datasource_id":"` + req.Data.Datasource + `", "message":"feature added"}}`
				if "" == req.Data.Datasource {
					err := errors.New("Missing required parameters")
					resp = `{"status": "error", "error": "` + err.Error() + `"}`
				} else {
					err = DB.InsertFeature(req.Data.Datasource, req.Data.Feature)
					if err != nil {
						resp = `{"status": "error", "error": "` + err.Error() + `"}`
					}
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "edit_feature" && authenticated:
				// {"method":"edit_feature"}
				resp := `{"status":"ok","data": {"datasource_id":"` + req.Data.Datasource + `", "message":"feature edited"}}`
				if "" == req.Data.Datasource {
					err := errors.New("Missing required parameters")
					resp = `{"status": "error", "error": "` + err.Error() + `"}`
				} else {
					err = DB.EditFeature(req.Data.Datasource, req.Data.GeoId, req.Data.Feature)
					if err != nil {
						fmt.Println(err)
						resp = `{"status": "error", "error": "` + err.Error() + `"}`
					}
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "create_datasource" && authenticated:
				// {"method":"create_datasource"}
				resp := `{"status":"ok","data":{}}`
				if "" != req.Data.Datasource {
					resp = `{"status":"ok","data": {"datasource_id":"` + req.Data.Datasource + `"}}`
					err = DB.InsertLayer(req.Data.Datasource, req.Data.Layer)
					if err != nil {
						resp = `{"status": "error", "error": "` + err.Error() + `"}`
					}
				} else {
					datasource_id, err := DB.NewLayer()
					resp = `{"status":"ok","data": {"datasource_id":"` + datasource_id + `"}}`
					if err != nil {
						resp = `{"status": "error", "error": "` + err.Error() + `"}`
					}
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "export_apikeys" && authenticated:
				// {"method":"export_apikeys"}
				resp := `{"status":"ok","data":{}}`
				apikeys, err := DB.SelectAll("apikeys")
				if err != nil {
					resp = `{"status":"error", "error":"` + err.Error() + `"}`
				} else {
					js, err := json.Marshal(apikeys)
					resp = `{"status":"ok","data":` + string(js) + `}`
					if err != nil {
						resp = `{"status":"error", "error":"` + err.Error() + `"}`
					}
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "export_apikey" && authenticated:
				// {"method":"export_apikey","apikey":"12dB6BlenIeB"}
				resp := `{"status":"ok","data":{}}`
				apikey, err := DB.GetCustomer(req.Apikey)
				if err != nil {
					resp = `{"status":"error", "error":"` + err.Error() + `"}`
				} else {
					js, err := json.Marshal(apikey)
					resp = `{"status":"ok","data":` + string(js) + `}`
					if err != nil {
						resp = `{"status":"error", "error":"` + err.Error() + `"}`
					}
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			// TODO: ERROR HANDLING
			case req.Method == "export_datasources" && authenticated:
				// {"method":"export_datasources"}
				resp := `{"status":"ok","data":{}}`
				layers, err := DB.SelectAll("layers")
				if err != nil {
					resp = `{"status":"error", "error":"` + err.Error() + `"}`
				} else {
					js, err := json.Marshal(layers)
					resp = `{"status":"ok","data":"` + string(js) + `"}`
					if err != nil {
						resp = `{"status":"error", "error":"` + err.Error() + `"}`
					}
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "export_datasource" && authenticated:
				// {"method":"export_datasource","datasource":"3b1f5d633d884b9499adfc9b49c45236"}
				resp := `{"status":"ok","data":{}}`
				layer, err := DB.GetLayer(req.Datasource)
				if err != nil {
					resp = `{"status":"error", "error":"` + err.Error() + `"}`
				} else {
					js, err := json.Marshal(layer)
					resp = `{"status":"ok","data":"` + string(js) + `"}`
					if err != nil {
						resp = `{"status":"error", "error":"` + err.Error() + `"}`
					}
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "import_file" && authenticated:
				// {"method":"import_file","file":"springfield_projects_edit.geojson"}
				resp := `{"status":"ok","data":{}}`
				result, err := importDatasource(req.File)
				resp = `{"status":"ok","data": {"datasource": "` + result + `"}}`
				if err != nil {
					resp = `{"status":"error", "error":"` + err.Error() + `"}`
				}
				conn.Write([]byte(resp + "\n"))
				success = true
			}

			if !authenticated {
				resp := `{"status": "error", "error": "connection not authenticated"}`
				conn.Write([]byte(resp + "\n"))
			} else if !success {
				resp := `{"status": "error", "error": "method not found"}`
				conn.Write([]byte(resp + "\n"))
			}

		}

	}
}

func importDatasource(importFile string) (string, error) {
	//fmt.Println("Importing", importFile)
	// get geojson file
	var geojsonFile string
	// check if file exists
	if _, err := os.Stat(importFile); os.IsNotExist(err) {
		return "", err
	}
	// ERROR
	// CRASHES IF NO "." character FOUND
	ext := strings.Split(importFile, ".")[1]
	// convert shapefile
	if ext == "shp" {
		// Convert .shp to .geojson
		geojsonFile := strings.Replace(importFile, ".shp", ".geojson", -1)
		fmt.Println("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojsonFile, importFile)
		out, err := exec.Command("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", geojsonFile, importFile).Output()
		if err != nil {
			return fmt.Sprintf("%v", out), err
		}
	} else if ext == "geojson" {
		geojsonFile = importFile
	} else {
		return fmt.Sprintf("Unsupported file type: %v", ext), errors.New(fmt.Sprintf("Unsupported file type: %v", ext))
	}
	// Read .geojson file
	file, err := ioutil.ReadFile(geojsonFile)
	if err != nil {
		return "", err
	}
	// Unmarshal to geojson struct
	geojs, err := geojson.UnmarshalFeatureCollection(file)
	if err != nil {
		return "", err
	}
	// Create datasource
	ds, _ := utils.NewUUID()
	DB.InsertLayer(ds, geojs)
	// Cleanup artifacts
	if geojsonFile != importFile {
		os.Remove(geojsonFile)
	}
	return ds, nil
}

/*

"insert_layer"
"delete_layer"

*/
