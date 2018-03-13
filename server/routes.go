package server

import (
	"net/http"

	"github.com/joyent/triton-service-groups/groups"
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
	router.Route{
		"CreateTemplate",
		http.MethodPost,
		"/v1/tsg/templates",
		templates_v1.Create,
	},
	router.Route{
		"DeleteTemplate",
		http.MethodDelete,
		"/v1/tsg/templates/{identifier}",
		templates_v1.Delete,
	},
}

var groupRoutes = router.Routes{
	router.Route{
		"GetGroup",
		http.MethodGet,
		"/v1/tsg/{identifier}",
		groups_v1.Get,
	},
	router.Route{
		"CreateGroup",
		http.MethodPost,
		"/v1/tsg",
		groups_v1.Create,
	},
	router.Route{
		"UpdateGroup",
		http.MethodPut,
		"/v1/tsg/{identifier}",
		groups_v1.Update,
	},
	router.Route{
		"DeleteGroup",
		http.MethodDelete,
		"/v1/tsg/{identifier}",
		groups_v1.Delete,
	},
	router.Route{
		"ListGroups",
		http.MethodGet,
		"/v1/tsg",
		groups_v1.List,
	},
}

var RoutingTable = router.RouteTable{
	templateRoutes,
	groupRoutes,
}
