/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	"flag"
	"fmt"
	"gospatial/app"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

var (
	port     int
	database string
	bind     string
	debug    bool
	version  bool
)

const (
	VERSION string = "1.6.4"
)

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		app.Error.Fatal(err)
	}
	// app.Info.Println(dir)
	db := strings.Replace(dir, "bin", "bolt", -1)
	app.Info.Println(db)
	flag.IntVar(&port, "p", 8080, "server port")
	// flag.StringVar(&database, "db", "bolt", "app database")
	flag.StringVar(&database, "db", db, "app database")
	flag.StringVar(&app.SuperuserKey, "s", "7q1qcqmsxnvw", "superuser key")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.BoolVar(&version, "v", false, "App Version")
	flag.Parse()
	if version {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}
	if debug {
		app.DebugMode()
	}
}

func main() {

	// Graceful shut down
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		for sig := range sigs {
			// sig is a ^C, handle it
			fmt.Printf("%s \n", sig)
			app.Info.Println("Gracefulling shutting down")
			app.Info.Println("Waiting for sockets to close...")
			for {
				if len(app.Hub.Sockets) == 0 {
					// auto backup
					app.DB.Backup("backup")
					app.Info.Println("Shutting down...")
					os.Exit(0)
				}
			}
		}
	}()

	// Initiate Database
	// app.DB = app.Database{File: "./" + database + ".db"}
	app.DB = app.Database{File: database + ".db"}
	app.DB.Init()
	// auto backup
	app.DB.Backup("backup")

	router := app.NewRouter()

	// Server static folder
	// router.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Start server
	app.Info.Printf("Magic happens on port %v...\n", port)
	if app.AppMode == "debug" {
		fmt.Printf("Magic happens on port %v...\n", port)
	}

	// https://golang.org/pkg/net/http/pprof/
	go func() {
		app.Info.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	bind := fmt.Sprintf(":%v", port)
	// bind := fmt.Sprintf("0.0.0.0:%v", port)
	// ListenAndServeTLS(bind, certFile, keyFile, router)
	// flag for certFile
	// flag for keyFile
	// if both there run TLS
	err := http.ListenAndServe(bind, router)
	if err != nil {
		panic(err)
	}
}
