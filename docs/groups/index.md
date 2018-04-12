# Groups

### GET `/v1/tsg/groups`

##### Inputs

| Name        | Type   | Description                                                                                                |
| ----------- | ------ | ---------------------------------------------------------------------------------------------------------- |
| id          | string | The universal identifier (UUID) of the group.                                                              |
| group_name  | string | The name of the group. The group name is limited to a maximum of 182 alphanumeric characters.              |
| template_id | string | A unique identifier for the template that the group is associated with.                                    |
| account_id  | string | A unique identifier for the account that the group is associated with.                                     |
| capacity    | number | The number of compute instances to run and maintain a specified number (the "desired count") of instances. |
| created_at  | string | When this group was created. ISO8601 date format.                                                          |
| updated_at  | string | When this group's details were last updated. ISO8601 date format.                                          |

##### Returns

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

### GET `/v1/tsg/groups/{identifier}`

##### Inputs

##### Returns

#### Example Request

```
curl -X GET https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}
```


#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
```

#### Sample Response

```
```

### GET `/v1/tsg/groups/{identifier}/instances`

##### Inputs

##### Returns

An array of objects, which contain:

| Name | Type | Description |
| ---- | ---- | ----------- |

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

##### Inputs

##### Returns

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

### PUT `/v1/tsg/groups/{identifier}`

#### Example Request

```
curl -X PUT https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}
```

#### Request Body

```json
{
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

### PUT `/v1/tsg/groups/{identifier}/increment`

##### Inputs

| Name           | Type   | Description                                                                        |
| -------------- | ------ | ---------------------------------------------------------------------------------- |
| instance_count | number | The number of compute instances by which to increase the current compute capacity. |
| max_instance   | number | Maximum limit of compute instances to maintain.                                    |

##### Returns

#### Example Request

```
curl -X PUT https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}/increment
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

### PUT `/v1/tsg/groups/{identifier}/decrement`

##### Inputs

| Name           | Type   | Description                                                                        |
| -------------- | ------ | ---------------------------------------------------------------------------------- |
| instance_count | number | The number of compute instances by which to decrease the current compute capacity. |
| min_instance   | number | Minimum limit of compute instances to maintain.                                    |

##### Returns

#### Example Request

```
curl -X PUT https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}/decrement
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

### DELETE `/v1/tsg/groups/{identifier}`

##### Inputs

##### Returns

#### Example Request

```
curl -X DELETE https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}
```

#### Request Headers 

```
Authorization: Signature keyId="/test-user-name/keys/aa:bb:cc:dd:9c:54:e9:78:3f:80:0d:ba:6b:c6:ff:44",algorithm="rsa-sha1",headers="date",signature="..."
Date: Fri, 06 Apr 2018 18:33:38 UTC
```

#### Sample Response

```
```
