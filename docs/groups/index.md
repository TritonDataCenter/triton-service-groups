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
curl -X GET -H "Content-Type: application/json" -H "" https://tsg.us-sw-1.svc.joyent.zone/v1/tsg/groups
```


#### Request Headers 

```

```

#### Sample Response

```

```

### GET `/v1/tsg/groups/{identifier}`

##### Inputs

##### Returns

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

```

#### Sample Response

```

```

### POST `/v1/tsg/groups`

##### Inputs

##### Returns

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

### PUT `/v1/tsg/groups/{identifier}/increment`

##### Inputs

| Name           | Type   | Description                                                                        |
| -------------- | ------ | ---------------------------------------------------------------------------------- |
| instance_count | number | The number of compute instances by which to increase the current compute capacity. |
| max_instance   | number | Maximum limit of compute instances to maintain.                                    |

##### Returns

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

### PUT `/v1/tsg/groups/{identifier}/decrement`

##### Inputs

| Name           | Type   | Description                                                                        |
| -------------- | ------ | ---------------------------------------------------------------------------------- |
| instance_count | number | The number of compute instances by which to decrease the current compute capacity. |
| min_instance   | number | Minimum limit of compute instances to maintain.                                    |

##### Returns

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

##### Inputs

##### Returns

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
