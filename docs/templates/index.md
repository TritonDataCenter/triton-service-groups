# Templates 

A template is a collection of configuration parameters that is used to launch a Triton instance. Templates are immutable,
therefore, once they are created, they cannot be changed, therefore a new template needs created. A template is made up
as follows:

| Field             | Type              |  Description                                                                  |
|----------         |-------------      |------                                                                         |
| ID                | string            | The UUID of the Template.                                                     |
| TemplateName      | string            | The name of the Template.                                                     |
| AccountID         | string            | The AccountID the Template is associated to.                                  | 
| Package           | string            | The ID of the Package to use when launching a Template.                       | 
| ImageID           | string            | The ID of the Package to use when launching a Template.                       | 
| FirewallEnabled   | bool              | Enable or disable the Firewall on the instances launched. Default is `false`. | 
| Networks          | []string          | A list of network IDs to attach to the instances launched.                    | 
| UserData          | string            | Data to be copied to the instances on boot.                                   | 
| MetaData          | map[string]string | A mapping of metadata to apply to the instances launched.                     | 
| Tags              | map[string]string | A mapping of tags to apply to the instances launched.                         |  

### GET /v1/tsg/templates

To list all of the templates associated with a specific Triton account, send a `GET` request to `/v1/tsg/templates`.

#### Example Request

```
curl -X GET -H "" "https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates"
```


#### Request Headers 

```

```

#### Sample Response

```

```

### GET /v1/tsg/templates/{identifier}

To show information about a specific template, send a `GET` request to `/v1/tsg/templates/{identifier}` 
using the request headers as detailed below.
 

#### Example Request

```
curl -X GET "https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates/319209784155176962"
```

#### Request Headers 

```

```

#### Sample Response

```

```

### POST /v1/tsg/templates

To create a new template, send a `POST` request to `/v1/tsg/templates`. The
request needs to include the headers as identified below. The attributes required to successfully create a template are
as follows:


 | Name             | Type              | Required  |
 |:----------:      |:-------------:    |:------:   |
 | Template Name    | string            | true      |
 | Package          | string            | true      |
 | ImageID          | string            | true      |
 | FirewallEnabled  | bool              |           |
 | Networks         | []string          |           |
 | UserData         | string            |           |
 | MetaData         | map[string]string |           |
 | Tags             | map[string]string |           |

#### Example Request

```
curl -X DELETE "https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates/319209784155176962"
```

#### Request Body

```

```

#### Request Headers 

```

```

#### Sample Response

```

```

### DELETE /v1/tsg/templates/{identifier}

Templates can be deleted by ID by sending a `DELETE` request to `/v1/tsg/templates/{identifier}`. The
request needs to include the headers as identified below.
 

#### Example Request

```
curl -X DELETE "https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates/319209784155176962"
```

#### Request Headers 

```

```

#### Sample Response

```

```
