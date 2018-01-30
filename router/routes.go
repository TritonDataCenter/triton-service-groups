package router

import (
	"net/http"

<<<<<<< HEAD
=======
	"path"

>>>>>>> ba2ab54... Changes after PR Review
	"github.com/joyent/triton-service-groups/groups"
	"github.com/joyent/triton-service-groups/templates"
)

<<<<<<< HEAD
=======
var (
	urlPrefix = "/v1/tsg"
)

>>>>>>> ba2ab54... Changes after PR Review
var templateRoutes = Routes{
	Route{
		"GetTemplate",
		http.MethodGet,
<<<<<<< HEAD
		"/v1/tsg/templates/{name}",
=======
		path.Join(urlPrefix, "templates", "{name}"),
>>>>>>> ba2ab54... Changes after PR Review
		templates_v1.Get,
	},
	Route{
		"CreateTemplate",
		http.MethodPost,
<<<<<<< HEAD
		"/v1/tsg/templates",
=======
		path.Join(urlPrefix, "templates"),
>>>>>>> ba2ab54... Changes after PR Review
		templates_v1.Create,
	},
	Route{
		"UpdateTemplate",
		http.MethodPut,
<<<<<<< HEAD
		"/v1/tsg/templates/{name}",
=======
		path.Join(urlPrefix, "templates", "{name}"),
>>>>>>> ba2ab54... Changes after PR Review
		templates_v1.Update,
	},
	Route{
		"DeleteTemplate",
		http.MethodDelete,
<<<<<<< HEAD
		"/v1/tsg/templates/{name}",
=======
		path.Join(urlPrefix, "templates", "{name}"),
>>>>>>> ba2ab54... Changes after PR Review
		templates_v1.Delete,
	},
	Route{
		"ListTemplates",
		http.MethodGet,
<<<<<<< HEAD
		"/v1/tsg/templates",
=======
		path.Join(urlPrefix, "templates"),
>>>>>>> ba2ab54... Changes after PR Review
		templates_v1.List,
	},
}

var groupRoutes = Routes{
	Route{
		"GetGroup",
		http.MethodGet,
<<<<<<< HEAD
		"/v1/tsg/{name}",
=======
		path.Join(urlPrefix, "{name}"),
>>>>>>> ba2ab54... Changes after PR Review
		groups_v1.Get,
	},
	Route{
		"CreateGroup",
		http.MethodPost,
<<<<<<< HEAD
		"/v1/tsg/",
=======
		path.Join(urlPrefix),
>>>>>>> ba2ab54... Changes after PR Review
		groups_v1.Create,
	},
	Route{
		"UpdateGroup",
		http.MethodPut,
<<<<<<< HEAD
		"/v1/tsg/{name}",
=======
		path.Join(urlPrefix, "{name}"),
>>>>>>> ba2ab54... Changes after PR Review
		groups_v1.Update,
	},
	Route{
		"DeleteGroup",
		http.MethodDelete,
<<<<<<< HEAD
		"/v1/tsg/{name}",
=======
		path.Join(urlPrefix, "{name}"),
>>>>>>> ba2ab54... Changes after PR Review
		groups_v1.Delete,
	},
	Route{
		"ListGroups",
		http.MethodGet,
<<<<<<< HEAD
		"/v1/tsg/",
=======
		path.Join(urlPrefix),
>>>>>>> ba2ab54... Changes after PR Review
		groups_v1.List,
	},
}
