# test-assignment

We have a two steam of data and we need to save it in the map:

[image:F249E0F5-014E-4E79-8EB4-FDA408DC2541-3213-00004D8212C79A29/62AA9010-02F4-4759-9681-47EEAB2AA678.png]

Incoming data:
userid, message
[event.go](https://github.com/hostpanda/test-assignment/blob/master/internal/domain/event/event.go#L8)

- POST request: /addHTTP  
- WS path /addWS

What we want to do during the interview:
- talk regarding architectures and how to keep your code base clean and readable
- create an ability to add new messages to the map. We want to have this structure:

```json
[{
  "userID":  jsonMessage
}, {
  "userID":  jsonMessage
}]	 
```

- use test_ws_client and test_http to make sure code will work in the concurrently enviroment
