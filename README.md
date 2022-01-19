# test-assignment

We have a two steam of data and we need to save it in the map:

![image](https://user-images.githubusercontent.com/97616185/150170255-c453c761-20c8-4adc-8b57-bea857159646.png)


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
