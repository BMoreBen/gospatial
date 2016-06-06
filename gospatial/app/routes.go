package app

import (
	"net/http"
)

type apiRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type apiRoutes []apiRoute

var routes = apiRoutes{
	// Health check
	apiRoute{"Ping", "GET", "/ping", PingHandler},

	// Layers
	apiRoute{"ViewLayers", "GET", "/api/v1/layers", ViewLayersHandler},
	apiRoute{"ViewLayer", "GET", "/api/v1/layer/{ds}", ViewLayerHandler},
	apiRoute{"NewLayer", "POST", "/api/v1/layer", NewLayerHandler},
	apiRoute{"DeleteLayer", "DELETE", "/api/v1/layer/{ds}", DeleteLayerHandler},
	apiRoute{"ShareLayerHandler", "PUT", "/api/v1/layer/{ds}", ShareLayerHandler},

	//
	apiRoute{"NewFeature", "POST", "/api/v1/layer/{ds}/feature", NewFeatureHandler},
	apiRoute{"ViewFeature", "GET", "/api/v1/layer/{ds}/feature/{k}", ViewFeatureHandler},

	// Superuser apiRoutes
	apiRoute{"NewCustomerHandler", "POST", "/api/v1/customer", NewCustomerHandler},

	// Web Client apiRoutes
	apiRoute{"Index", "GET", "/", IndexHandler},
	apiRoute{"MapNew", "GET", "/map", MapHandler},
	apiRoute{"CustomerManagement", "GET", "/management", CustomerManagementHandler},

	// Web Socket apiRoute
	apiRoute{"Socket", "GET", "/ws/{ds}", serveWs},

	// Experimental
	apiRoute{"UnloadLayer", "GET", "/management/unload/{ds}", UnloadLayer},
	apiRoute{"LoadedLayers", "GET", "/management/loaded", LoadedLayers},
	apiRoute{"LoadedLayers", "GET", "/management/profile", ServerProfile},
}
