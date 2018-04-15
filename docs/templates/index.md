# Templates

A template is a collection of configuration parameters that are used to launch a compute instance.

Templates are immutable, therefore, once they are created they cannot be changed. If you need to
make changes, you must create a new template.

A template object contains the following fields:

| Field            | Type             | Description                                                                              |
| ---------------- | ---------------- | ---------------------------------------------------------------------------------------- |
| id               | string           | The universal identifier (UUID) of the template.                                         |
| template_name    | string           | The name of the template.                                                                |
| package          | string           | The unique identifier (UUID) of the package to use when launching compute instances.     |
| image_id         | string           | The unique identifier (UUID) of the image to use when launching compute instances.       |
| firewall_enabled | boolean          | Whether to enable or disable the firewall on the instances launched. Default is `false`. |
| networks         | array of strings | A list of unique network identifiers to attach to the compute instances launched.        |
| user_data        | string           | An arbitrary data to be copied to the instances on boot. This will not be executed.      |
| meta_data        | object           | A mapping of metadata (a key-value pairs) to apply to the instances launched.            |
| tags             | object           | A mapping of tags (a key-value pairs) to apply to the instances launched.                |
| created_at       | string           | When this template was created. ISO8601 date format.                                     |

The template object shares attributes with the compute instance object as found in the
[Joyent CloudAPI][1] documentation in the [instances][2] section.

### POST `/v1/tsg/templates`

To create a new template, send a `POST` request to `/v1/tsg/templates`. The request must include
the authentication headers. The attributes required to successfully create a template are as
follows:

| Name             | Type             | Description                                                                          | Required   |
| ---------------- | ---------------- | ------------------------------------------------------------------------------------ | :--------: |
| template_name    | string           | The name of the template.                                                            | Yes        |
| package          | string           | The unique identifier (UUID) of the package to use when launching compute instances. | Yes        |
| image_id         | string           | The unique identifier (UUID) of the image to use when launching compute instances.   | Yes        |
| firewall_enabled | boolean          | Whether to enable or disable the firewall on the instances launched.                 | No         |
| networks         | array of strings | A list of unique network identifiers to attach to the compute instances launched.    | No         |
| user_data        | string           | An arbitrary data to be copied to the instances on boot. This will not be executed.  | No         |
| meta_data        | object           | A mapping of metadata (a key-value pairs) to apply to the instances launched.        | No         |
| tags             | object           | A mapping of tags (a key-value pairs) to apply to the instances launched.            | No         |

A successful request will return a `201 Created` HTTP response code, and object representing newly
created template in the response body.

#### Example Request

```
curl -X POST -H 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates
```

#### Request Body

```
{
    "template_name": "jolly-jelly",
    "package": "14aba044-d0f8-11e5-8c88-eb339a5da5d0",
    "image_id": "342045ce-6af1-4adf-9ef1-e5bfaf9de28c",
    "firewall_enabled": false,
    "networks": [
        "27ea1d5f-df02-410e-843a-c60dba9ec5ca"
    ],
    "tags": {
        "owner": "user"
    },
    "metadata": {
        "user-script": "#!/bin/bash\ndate\n"
    }
}
```

#### Request Headers

```
Date: Sun, 15 Apr 2018 20:24:06 GMT
ontent-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Sample Response

```
{
    "id": "29a08459-1a41-4ec9-bbb7-5c737f17a463",
    "template_name": "jolly-jelly",
    "package": "14aba044-d0f8-11e5-8c88-eb339a5da5d0",
    "image_id": "342045ce-6af1-4adf-9ef1-e5bfaf9de28c",
    "firewall_enabled": false,
    "networks": [
        "27ea1d5f-df02-410e-843a-c60dba9ec5ca"
    ],
    "userdata": "",
    "metadata": {
        "user-script": "#!/bin/bash\ndate\n"
    },
    "tags": {
        "owner": "user"
    },
    "created_at": "2018-04-15T20:24:07.481363Z"
}
```

### DELETE `/v1/tsg/templates/{UUID}`

To delete a template, send a `DELETE` request to `/v1/tsg/templates/{UUID}`, where the `{UUID}` is the unique
identifier (UUID) of the template. The request must include the authentication headers.

A successful request will return a `204 No Content` HTTP status code, and no body will be
included in the response.

#### Example Request

```
curl -X DELETE -H 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates/29a08459-1a41-4ec9-bbb7-5c737f17a463
```

#### Request Headers

```
Date: Sun, 15 Apr 2018 20:34:56 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Sample Response

```
204 No Content
```

### GET `/v1/tsg/templates`

To list all of the templates, send a `GET` request to `/v1/tsg/templates`. The request must include
the authentication headers.

A successful request will return a `200 OK` HTTP status code, and a list of objects representing
a template in the response body.

#### Example Request

```
curl -X GET -H 'Content-Type: application/json'  https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates
```

#### Request Headers

```
Date: Sun, 15 Apr 2018 20:30:16 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Sample Response

```
[
    {
        "id": "29a08459-1a41-4ec9-bbb7-5c737f17a463",
        "template_name": "jolly-jelly",
        "package": "14aba044-d0f8-11e5-8c88-eb339a5da5d0",
        "image_id": "342045ce-6af1-4adf-9ef1-e5bfaf9de28c",
        "firewall_enabled": false,
        "networks": [
            "27ea1d5f-df02-410e-843a-c60dba9ec5ca"
        ],
        "userdata": "",
        "metadata": {
            "user-script": "#!/bin/bash\ndate\n"
        },
        "tags": {
            "owner": "user"
        },
        "created_at": "2018-04-15T20:24:07.481363Z"
    }
]
```

### GET `/v1/tsg/templates/{UUID}`

To show information about a specific template, send a `GET` request to `/v1/tsg/templates/{UUID}`,
where the `{UUID}` is the unique identifier (UUID) of the template. The request must include the
authentication headers.

A successful request will return a `200 OK` HTTP response code, and an object representing
a template in the response body.

#### Example request

```
curl -X GET -H 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/templates/29a08459-1a41-4ec9-bbb7-5c737f17a463
```

#### Example request headers

```
Date: Sun, 15 Apr 2018 20:32:47 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Example response

```
{
    "id": "29a08459-1a41-4ec9-bbb7-5c737f17a463",
    "template_name": "jolly-jelly",
    "package": "14aba044-d0f8-11e5-8c88-eb339a5da5d0",
    "image_id": "342045ce-6af1-4adf-9ef1-e5bfaf9de28c",
    "firewall_enabled": false,
    "networks": [
        "27ea1d5f-df02-410e-843a-c60dba9ec5ca"
    ],
    "userdata": "",
    "metadata": {
        "user-script": "#!/bin/bash\ndate\n"
    },
    "tags": {
        "owner": "user"
    },
    "created_at": "2018-04-15T20:24:07.481363Z"
}
```

[1]: https://apidocs.joyent.com/cloudapi
[2]: https://apidocs.joyent.com/cloudapi/#instances
