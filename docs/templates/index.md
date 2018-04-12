# Templates 

A template is a collection of configuration parameters that are used to launch a Triton instance. Templates are immutable,  therefore, once they are created they cannot be changed. If you need to make changes, you must create a new template.

A template is made up as follows:

| Field            | Type             | Description                                                                   |
| ---------------- | ---------------- | ----------------------------------------------------------------------------- |
| id               | string           | The UUID of the template.                                                     |
| template_name    | string           | The name of the template.                                                     |
| account_id       | string           | The account ID the template is associated to.                                 |
| package          | string           | The ID of the package to use when launching a template.                       |
| image_id         | string           | The ID of the image to use when launching a template.                         |
| firewall_enabled | boolean          | Enable or disable the firewall on the instances launched. Default is `false`. |
| networks         | array of strings | A list of network IDs to attach to the instances launched.                    |
| user_data        | string           | Data to be copied to the instances on boot.                                   |
| meta_data        | object           | A mapping of metadata to apply to the instances launched.                     |
| tags             | object           | A mapping of tags to apply to the instances launched.                         |
| created_at       | string           | When this template was created. ISO8601 date format.                          |

### GET `/v1/tsg/templates`

To list all of the templates associated with a specific Triton account, send a `GET` request to `/v1/tsg/templates` with
the request headers detailed below.

#### Example Request

```
curl -X GET "https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates"
```


#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
```

#### Sample Response

```
[
    {
        "id": "437c560d-b1a9-4dae-b3b3-6dbabb7d23a7",
        "template_name": "test-template-6",
        "account_id": "6f873d02-172c-418f-8416-4da2b50d5c53",
        "package": "test-package",
        "image_id": "49b22aec-0c8a-11e6-8807-a3eb4db576ba",
        "firewall_enabled": false,
        "networks": [
            "f7ed95d3-faaf-43ef-9346-15644403b963"
        ],
        "metadata": {
    	    "root_pw": "s8v9kuht5e"
    	}, 
        "tags": {
            "role": "web",
            "owner": "api-team"
        },
        "created_at": "2018-04-12T15:59:08.098244Z"
    }
}
```

### GET `/v1/tsg/templates/{UUID}`

To show information about a specific template, send a `GET` request to `/v1/tsg/templates/{UUID}` 
using the request headers as detailed below.
 

#### Example Request

```
curl -X GET "https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates/319209784155176962"
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
```

#### Sample Response

```
{
    "id": "437c560d-b1a9-4dae-b3b3-6dbabb7d23a7",
    "template_name": "test-template-6",
    "account_id": "6f873d02-172c-418f-8416-4da2b50d5c53",
    "package": "test-package",
    "image_id": "49b22aec-0c8a-11e6-8807-a3eb4db576ba",
    "firewall_enabled": false,
    "networks": [
        "f7ed95d3-faaf-43ef-9346-15644403b963"
    ],
    "metadata": {
	    "root_pw": "s8v9kuht5e"
	}, 
    "tags": {
        "role": "web",
        "owner": "api-team"
    },
    "created_at": "2018-04-12T15:59:08.098244Z"
}
```

### POST `/v1/tsg/templates`

To create a new template, send a `POST` request to `/v1/tsg/templates`. The
request needs to include the headers as identified below. A successful creation will return 
a `201 created` HTTP Response Code. The attributes required to successfully create a template are as follows:


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
curl -X POST "https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates"
```

#### Request Body

```
{
    "template_name": "silly-salmon",
    "package": "7b17343c-94af-6266-e0e8-893a3b9993d0",
    "image_id": "49b22aec-0c8a-11e6-8807-a3eb4db576ba",
    "firewall_enabled": false,
    "networks": [
        "f7ed95d3-faaf-43ef-9346-15644403b963"
    ],
    "metadata": {
	    "root_pw": "s8v9kuht5e"
	},
    "tags": {
    	"role": "api",
    	"owner": "design"
    },
    "created_at": "2018-04-12T15:59:08.098244Z"
}
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
```

#### Sample Response

```
{
    "id": "5ffdfc6a-42ad-40c3-aa9f-e5e6d0c33003",
    "template_name": "silly-salmon",
    "account_id": "6f873d02-172c-418f-8416-4da2b50d5c53",
    "package": "7b17343c-94af-6266-e0e8-893a3b9993d0",
    "image_id": "49b22aec-0c8a-11e6-8807-a3eb4db576ba",
    "firewall_enabled": false,
    "networks": [
        "f7ed95d3-faaf-43ef-9346-15644403b963"
    ],
    "userdata": "",
    "metadata": {
        "root_pw": "s8v9kuht5e"
    },
    "tags": {
        "owner": "design",
        "role": "api"
    },
    "created_at": "2018-04-12T15:59:08.098244Z"
}
```

### DELETE `/v1/tsg/templates/{UUID}`

Templates can be deleted by sending a `DELETE` request to `/v1/tsg/templates/{UUID}`, where the `{UUID}` is the template ID. The request must include the headers as identified below. A successful delete will return a HTTP status code of `204 No Content`.
 

#### Example Request

```
curl -X DELETE "https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates/319209784155176962"
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
```

#### Sample Response

```
204 No Content
```
