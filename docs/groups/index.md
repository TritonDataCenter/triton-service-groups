# Groups

A group is made up as follows:

| Name        | Type   | Description                                                                                                |
| ----------- | ------ | ---------------------------------------------------------------------------------------------------------- |
| id          | string | The universal identifier (UUID) of the group.                                                              |
| group_name  | string | The name of the group. The group name is limited to a maximum of 182 alphanumeric characters.              |
| template_id | string | A unique identifier for the template that the group is associated with.                                    |
| account_id  | string | A unique identifier for the account that the group is associated with.                                     |
| capacity    | number | The number of compute instances to run and maintain a specified number (the "desired count") of instances. |
| created_at  | string | When this group was created. ISO8601 date format.                                                          |
| updated_at  | string | When this group's details were last updated. ISO8601 date format.                                          |

### GET `/v1/tsg/groups`

To list all of the groups associated with a specific Triton account, send a `GET`
request to `/v1/tsg/groups` with the request headers detailed below.

#### Example Request

```
curl -X GET https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
```

#### Sample Response

```
```

### GET `/v1/tsg/groups/{UUID}`

To show information about a specific group, send a `GET` request to `/v1/tsg/groups/{UUID}`
using the request headers as detailed below.

#### Example Request

```
curl -X GET https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{UUID}
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
```

#### Sample Response

```
```

### GET `/v1/tsg/groups/{UUID}/instances`

#### Example Request

```
curl -X GET -H "Content-Type: application/json" -H "" https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}/instances
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
```

#### Sample Response

```
```

### POST `/v1/tsg/groups`


To create a new group, send a `POST` request to /v1/tsg/groups. The request needs
to include the headers as identified below. A successful creation will return a
`201 Created` HTTP Response Code. The attributes required to successfully create
a group are as follows:

| Name        | Type   | Required |
| ----------- | ------ | -------- |
| group_name  | string | Yes      |
| account_id  | string | Yes      |
| template_id | string | Yes      |
| capacity    | string | Yes      |


#### Example Request

```
curl -X POST https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups

```

### Request Body

```json
{
    "group_name": "jolly-jelly",
    "account_id": "6f873d02-172c-418f-8416-4da2b50d5c53",
    "template_id": "437c560d-b1a9-4dae-b3b3-6dbabb7d23a7",
    "capacity": 20
}
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
Content-Type: application/json

```

#### Sample Response

```
```

### PUT `/v1/tsg/groups/{UUID}`

#### Example Request

```
curl -X PUT https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{UUID}
```

#### Request Body

```
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
Content-Type: application/json
```

#### Sample Response

```
```

### PUT `/v1/tsg/groups/{UUID}/increment`

| Name           | Type   | Description                                                                        | Required |
| -------------- | ------ | ---------------------------------------------------------------------------------- | -------- |
| instance_count | number | The number of compute instances by which to increase the current compute capacity. | Yes      |
| max_instance   | number | Maximum limit of compute instances to maintain.                                    | Yes      |

#### Example Request

```
curl -X PUT https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{UUID}/increment
```

#### Request Body

```json
{
    "instance_count": 25,
    "max_instance": 50
}
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
Content-Type: application/json
```

#### Sample Response

```
```

### PUT `/v1/tsg/groups/{UUID}/decrement`

| Name           | Type   | Description                                                                        | Required |
| -------------- | ------ | ---------------------------------------------------------------------------------- | -------- |
| instance_count | number | The number of compute instances by which to decrease the current compute capacity. | Yes      |
| min_instance   | number | Minimum limit of compute instances to maintain.                                    | Yes      |

#### Example Request

```
curl -X PUT https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{UUID}/decrement
```

#### Request Body

```json
{
    "instance_count": 15,
    "min_instance": 5
}
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
Content-Type: application/json
```

#### Sample Response

```
```

### DELETE `/v1/tsg/groups/{UUID}`

Groups can be deleted by sending a `DELETE` request to `/v1/tsg/groups/{UUID}`,
where the `{UUID}` is the group ID. The request must include the headers as detailed
below. A successful delete will return a HTTP status code of `204 No Content`.

#### Example Request

```
curl -X DELETE "https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{UUID}"
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
