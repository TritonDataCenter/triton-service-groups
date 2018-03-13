package server

import (
	"net/http"

	"github.com/joyent/triton-service-groups/server/router"
	"github.com/joyent/triton-service-groups/templates"
)

var templateRoutes = router.Routes{
	router.Route{
		"ListTemplates",
		http.MethodGet,
		"/v1/tsg/templates",
		templates_v1.List,
	},
	router.Route{
		"GetTemplate",
		http.MethodGet,
		"/v1/tsg/templates/{identifier}",
		templates_v1.Get,
	},
	// router.Route{
	// 	"CreateTemplate",
	// 	http.MethodPost,
	// 	"/templates",
	// 	templates_v1.CreateHandler,
	// },
	// router.Route{
	// 	"DeleteTemplate",
	// 	http.MethodDelete,
	// 	"/templates/{identifier}",
	// 	templates_v1.DeleteHandler,
	// },
}

var routingTable = []router.Routes{templateRoutes}
