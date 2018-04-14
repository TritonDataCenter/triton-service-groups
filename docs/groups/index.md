# Groups

A group contains a collection of compute instances that share similar characteristics and are treated as a
logical grouping for the purposes of instance scaling and management. A group starts by launching enough
compute instances to meet its desired capacity. A group object contains the following fields:

| Name        | Type   | Description                                                                                                |
| ----------- | ------ | ---------------------------------------------------------------------------------------------------------- |
| id          | string | The universal identifier (UUID) of the group.                                                              |
| group_name  | string | The name of the group. The group name is limited to a maximum of 182 alphanumeric characters.              |
| template_id | string | A unique identifier for the template that the group is associated with.                                    |
| capacity    | number | The number of compute instances to run and maintain a specified number (the "desired count") of instances. |
| created_at  | string | When this group was created. ISO8601 date format.                                                          |
| updated_at  | string | When this group's details were last updated. ISO8601 date format.                                          |

### POST `/v1/tsg/groups`

To create a new group, send a `POST` request to `/v1/tsg/groups`. The request must include the authentication
headers. The attributes required to successfully create a group are as follows:

| Name        | Type   | Description                                                                                                | Required |
| ----------- | ------ | ---------------------------------------------------------------------------------------------------------- | -------- |
| group_name  | string | The name of the group. The group name is limited to a maximum of 182 alphanumeric characters.              | Yes      |
| template_id | string | A unique identifier for the template that the group is associated with.                                    | Yes      |
| capacity    | string | The number of compute instances to run and maintain a specified number (the "desired count") of instances. | Yes      |

A successful creation will return a `200 OK` HTTP response code, and object representing newly created
group in the response body.

#### Example request

```
curl -X POST -H 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups

```

#### Example request body

```
{
    "group_name": "jolly-jelly",
    "account_id": "6f873d02-172c-418f-8416-4da2b50d5c53",
    "template_id": "437c560d-b1a9-4dae-b3b3-6dbabb7d23a7",
    "capacity": 5
}
```

#### Example request headers

```
Date: Sat, 14 Apr 2018 15:24:18 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Example response

```
{
    "id": "bc351939-48a1-4f87-af62-ae8ea9f0acf6",
    "group_name": "jolly-jelly",
    "template_id": "ebec1e0c-9caa-47d9-97e2-3e31d277a35f",
    "capacity": 5,
    "created_at": "2018-04-14T15:24:20.205784Z",
    "updated_at": "2018-04-14T15:24:20.205784Z"
}
```

### DELETE `/v1/tsg/groups/{UUID}`

To delete a group, send a `DELETE` request to `/v1/tsg/groups/{UUID}`, where the `{UUID}` is the unique
identifier (UUID) of the group. The request must include the authentication headers.

A successful deletion will return a `204 No Content` HTTP status code, and no body will be included in the response.

#### Example request

```
curl -X DELETE -H 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/bc351939-48a1-4f87-af62-ae8ea9f0acf6
```

#### Example request headers

```
Date: Sat, 14 Apr 2018 15:35:09 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Example response

```
204 No Content
```

### GET `/v1/tsg/groups`

To list all of the groups, send a `GET` request to `/v1/tsg/groups`. The request must include the authentication
headers.

A successful request to list groups will return a `200 OK` HTTP status code, and a list of objects representing
a group in the response body.

#### Example request

```
curl -X GET -H 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups
```

#### Example request headers

```
Date: Sat, 14 Apr 2018 15:54:01 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Example response

```
[
    {
        "id": "0a774d9e-7c76-4740-8ecc-20c3846956c7",
        "group_name": "cuddly-cat",
        "template_id": "ebec1e0c-9caa-47d9-97e2-3e31d277a35f",
        "capacity": 5,
        "created_at": "2018-04-14T16:02:04.032525Z",
        "updated_at": "2018-04-14T16:02:04.032525Z"
    },
    {
        "id": "722d25ed-f32a-4944-9861-8990e204850e",
        "group_name": "jolly-jelly",
        "template_id": "ebec1e0c-9caa-47d9-97e2-3e31d277a35f",
        "capacity": 5,
        "created_at": "2018-04-14T15:50:08.872758Z",
        "updated_at": "2018-04-14T15:50:08.872758Z"
    }
]
```

### GET `/v1/tsg/groups/{UUID}`

To show information about a specific group, send a `GET` request to `/v1/tsg/groups/{UUID}`, where the `{UUID}`
is the unique identifier (UUID) of the group. The request must include the authentication headers.

A successful creation will return a `200 OK` HTTP response code, and object representing a group in the response
body.

#### Example request

```
curl -X GET 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/722d25ed-f32a-4944-9861-8990e204850e
```

#### Example request headers

```
Date: Sat, 14 Apr 2018 16:07:45 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Example response

```
{
    "id": "722d25ed-f32a-4944-9861-8990e204850e",
    "group_name": "jolly-jelly",
    "template_id": "ebec1e0c-9caa-47d9-97e2-3e31d277a35f",
    "capacity": 5,
    "created_at": "2018-04-14T15:50:08.872758Z",
    "updated_at": "2018-04-14T15:50:08.872758Z"
}
```

### PUT `/v1/tsg/groups/{UUID}`

To update a group (e.g. capacity, etc.), send a `PUT` request to `/v1/tsg/groups/{UUID}`, where the `{UUID}`
is the unique identifier (UUID) of the group. The request must include the authentication headers. The
attributes required to successfully create a group are as follows:

| Name        | Type   | Description                                                                                                | Required |
| ----------- | ------ | ---------------------------------------------------------------------------------------------------------- | -------- |
| group_name  | string | The name of the group. The group name is limited to a maximum of 182 alphanumeric characters.              | Yes      |
| template_id | string | A unique identifier for the template that the group is associated with.                                    | Yes      |
| capacity    | string | The number of compute instances to run and maintain a specified number (the "desired count") of instances. | Yes      |

A successful update will return a `200 OK` HTTP response code, and object representing a group in the response
body.

#### Example request

```
curl -X PUT 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/722d25ed-f32a-4944-9861-8990e204850e
```

#### Example request body

```
{
	"group_name": "jolly-jelly",
	"template_id": "ebec1e0c-9caa-47d9-97e2-3e31d277a35f",
    "capacity": 10
}
```

#### Example request headers

```
Date: Sat, 14 Apr 2018 16:14:07 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date"
```

#### Example response

```
{
    "id": "722d25ed-f32a-4944-9861-8990e204850e",
    "group_name": "jolly-jelly",
    "template_id": "ebec1e0c-9caa-47d9-97e2-3e31d277a35f",
    "capacity": 10,
    "created_at": "2018-04-14T15:50:08.872758Z",
    "updated_at": "2018-04-14T16:14:08.70981Z"
}
```

### GET `/v1/tsg/groups/{UUID}/instances`

To list all of the compute instances associated with a specific group, send a `GET` request to `/v1/tsg/groups/{UUID}/instances`,
where the `{UUID}` is the unique identifier (UUID) of the group. The request must include the authentication headers.

A successful request to list compute instances associated with a given group will return a `200 OK` HTTP status code,
and a list of objects representing a compute instance. The details about an object representing a compute instance
can be found in the [Joyent CloudAPI](https://apidocs.joyent.com/cloudapi/) documentation in the
[instances](https://apidocs.joyent.com/cloudapi/#instances) section.

#### Example request

```
curl -X GET 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/722d25ed-f32a-4944-9861-8990e204850e/instances
```

#### Example request headers

```
Date: Sat, 14 Apr 2018 16:28:27 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Example response

```
[
    {
        "id": "77df263d-2278-e820-928f-de07ff070122",
        "name": "tsg-722d25ed-77df263d",
        "type": "virtualmachine",
        "brand": "kvm",
        "state": "running",
        "image": "342045ce-6af1-4adf-9ef1-e5bfaf9de28c",
        "memory": 3840,
        "disk": 51200,
        "metadata":{
            "root_authorized_keys": "..."
        },
        "tags":{
            "owner": "user",
            "tsg.name": "jolly-jelly"
        },
        "created":"2018-04-14T16:26:26.985Z",
        "updated":"2018-04-14T16:26:39Z",
        "docker": false,
        "ips":[
            "10.0.0.1"
        ],
        "networks":[
            "0106d35d-90fa-48dc-b0dd-2ebcbaa64ca0"
        ],
        "primaryIp": "10.0.0.1",
        "firewall_enabled": false,
        "compute_node": "44454c4c-5300-1048-804a-b8c04f524432",
        "package": "k4-general-kvm-3.75G",
        "dns_names": null,
        "deletion_protection": false,
        "CNS":{
            "Disable": false,
            "ReversePTR": "",
            "Services": null
        }
    }
]
```

### PUT `/v1/tsg/groups/{UUID}/increment`

To a number of new compute instances to a group while maintaining the maximum limit, send a `PUT` request to `/v1/tsg/groups/{UUID}/increment`,
where the `{UUID}` is the unique identifier (UUID) of the group. The request must include the authentication headers.
The attributes required to successfully create a group are as follows:

| Name           | Type   | Description                                                           | Required |
| -------------- | ------ | --------------------------------------------------------------------- | -------- |
| instance_count | number | The number of compute instances to add to the current group capacity. | Yes      |
| max_instance   | number | Maximum number of compute instances allowed in the group.             | Yes      |

A successful request to add a number of compute instances to a group will return a `202 Accepted` HTTP status code,
and no body will be included in the response.

#### Example request

```
curl -X PUT -H 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/722d25ed-f32a-4944-9861-8990e204850e/increment
```

#### Example request body

```
{
    "instance_count": 1,
    "max_instance": 5
}
```

#### Example request headers

```
Date: Sat, 14 Apr 2018 19:58:05 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Example response

```
202 Accepted
```

### PUT `/v1/tsg/groups/{UUID}/decrement`

To remove a numnber of compute instances from a group while maintaining the minimum limit, send a `PUT` request to `/v1/tsg/groups/{UUID}/decrement`,
where the `{UUID}` is the unique identifier (UUID) of the group. The request must include the authentication headers.
The attributes required to successfully create a group are as follows:

| Name           | Type   | Description                                                         | Required |
| -------------- | ------ | ------------------------------------------------------------------- | -------- |
| instance_count | number | The number of compute instances to remove from the group capacity.  | Yes      |
| min_instance   | number | Minimum number of compute instances allowed in the group.           | Yes      |

A successful request to remove a number of compute instances from a group will return a `202 Accepted` HTTP status code,
and no body will be included in the response.

#### Example request

```
curl -X PUT -H 'Content-Type: application/json' https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/722d25ed-f32a-4944-9861-8990e204850e/decrement
```

#### Example request body

```
{
    "instance_count": 1,
    "min_instance": 5
}
```

#### Example request headers

```
Date: Sat, 14 Apr 2018 20:02:18 GMT
Content-Type: application/json
Authorization: Signature keyId="/user/keys/32:98:8a:b8:b3:a3:cb:f4:3c:42:24:d8:44:b8:0b:63",algorithm="rsa-sha256",headers="date" ...
```

#### Example response

```
202 Accepted
```
