# Groups

| Name        | Type   | Description                                                                                                |
| ----------- | ------ | ---------------------------------------------------------------------------------------------------------- |
| id          | string | The universal identifier (UUID) of the group.                                                              |
| group_name  | string | The name of the group. The group name is limited to a maximum of 182 alphanumeric characters.              |
| template_id | string | A unique identifier for the template that the group is associated with.                                    |
| account_id  | string | A unique identifier for the account that the group is associated with.                                     |
| capacity    | number | The number of compute instances to run and maintain a specified number (the "desired count") of instances. |
| created_at  | string | When this group was created. ISO8601 date	format.                                                         |
| updated_at  | string | When this group's details were last updated. ISO8601 date format.                                          |

### GET `/v1/tsg/groups`

#### Example Request

```
curl -X GET -H "Content-Type: application/json" -H "" https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups
```


#### Request Headers 

```

```

#### Sample Response

```

```

### GET `/v1/tsg/groups/{identifier}`

#### Example Request

```
curl -X GET -H "Content-Type: application/json" -H "" https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}
```


#### Request Headers 

```

```

#### Sample Response

```

```

### GET `/v1/tsg/groups/{identifier}/instances`

#### Example Request

```
curl -X GET -H "Content-Type: application/json" -H "" https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}/instances
```


#### Request Headers 

```

```

#### Sample Response

```

```

### POST `/v1/tsg/groups`

#### Example Request

```
curl -X POST -H "Content-Type: application/json" -H "" https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups -d '{}'

```


#### Request Headers 

```

```

#### Sample Response

```

```

### PUT `/v1/tsg/groups/{identifier}`

#### Example Request

```
curl -X PUT -H "Content-Type: application/json" -H "" https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier} -d '{}'
```


#### Request Headers 

```

```

#### Sample Response

```

```

| Name           | Type   | Description |
| -------------- | ------ | ----------- |
| instance_count | number |             |
| max_instance   | number |             |
| min_instance   | number |             |


### PUT `/v1/tsg/groups/{identifier}/increment`

#### Example Request

```
curl -X PUT -H "Content-Type: application/json" -H "" https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}/increment -d '{}'
```


#### Request Headers 

```

```

#### Sample Response

```

```

| Name           | Type   | Description |
| -------------- | ------ | ----------- |
| instance_count | number |             |
| max_instance   | number |             |
| min_instance   | number |             |

### PUT `/v1/tsg/groups/{identifier}/decrement`

#### Example Request

```
curl -X PUT -H "Content-Type: application/json" -H "" https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}/decrement -d '{}'
```


#### Request Headers 

```

```

#### Sample Response

```

```

### DELETE `/v1/tsg/groups/{identifier}`

#### Example Request

```
curl -X DELETE -H "Content-Type: application/json" -H "" https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups/{identifier}
```


#### Request Headers 

```

```

#### Sample Response

```

```
