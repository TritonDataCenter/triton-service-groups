# Groups

| Name        | Type    | Description                                                                                                |
| ----------- | ------- | ---------------------------------------------------------------------------------------------------------- |
| id          | string  | The universal identifier (UUID) of the Group.                                                              |
| group_name  | string  | The name of the Group.                                                                                     |
| template_id | string  | A unique identifier for the Template that the Group is associated with.                                    |
| account_id  | string  | A unique identifier for the Account that the Group is associated with.                                     |
| capacity    | integer | The number of compute instances to run and maintain a specified number (the "desired count") of instances. |

### GET /v1/tsg/groups

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

### GET /v1/tsg/groups/{identifier}

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

### POST /v1/tsg/groups

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

### PUT /v1/tsg/groups/{identifier}

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

### DELETE /v1/tsg/groups/{identifier}

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
