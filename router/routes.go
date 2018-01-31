package router

import (
	"net/http"

	"github.com/joyent/triton-service-groups/groups"
	"github.com/joyent/triton-service-groups/templates"
)

var templateRoutes = Routes{
	Route{
		"GetTemplate",
		http.MethodGet,
		"/v1/tsg/templates/{name}",
		templates_v1.Get,
	},
	Route{
		"CreateTemplate",
		http.MethodPost,
		"/v1/tsg/templates",
		templates_v1.Create,
	},
	Route{
		"UpdateTemplate",
		http.MethodPut,
		"/v1/tsg/templates/{name}",
		templates_v1.Update,
	},
	Route{
		"DeleteTemplate",
		http.MethodDelete,
		"/v1/tsg/templates/{name}",
		templates_v1.Delete,
	},
	Route{
		"ListTemplates",
		http.MethodGet,
		"/v1/tsg/templates",
		templates_v1.List,
	},
}

var groupRoutes = Routes{
	Route{
		"GetGroup",
		http.MethodGet,
		"/v1/tsg/{name}",
		groups_v1.Get,
	},
	Route{
		"CreateGroup",
		http.MethodPost,
		"/v1/tsg/",
		groups_v1.Create,
	},
	Route{
		"UpdateGroup",
		http.MethodPut,
		"/v1/tsg/{name}",
		groups_v1.Update,
	},
	Route{
		"DeleteGroup",
		http.MethodDelete,
		"/v1/tsg/{name}",
		groups_v1.Delete,
	},
	Route{
		"ListGroups",
		http.MethodGet,
		"/v1/tsg/",
		groups_v1.List,
	},
}
