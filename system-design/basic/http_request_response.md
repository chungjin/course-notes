# Http Request/Response protocol



## Short polling

With the traditional or "short polling" technique, a client sends regular requests to the server and each request attempts to "pull" any available events or data.  If there are no events or data available, the server returns an empty response and the client waits for some time before sending another poll request.



Pros: 

- Implemetation is simple

Cons:

- Unacceptable burden on the server, the network, or both.

## Long polling

The server achieves these efficiencies by responding to a request only when a particular event, status, or timeout has occurred.  Once the server sends a long poll response, typically the client immediately sends a new long poll request. 

Pros:

- save traffic

Cons:

- allocate resources to keep the connection open



messenger这种应用不适合的原因在于:

- 如果使用了LB，那么收到message的server不一定保存有long-pulling connection to the client.
- 由于不是bidirection, 无法得知client是否已经disconnect

## Websocket

Reference: [[*WebSocket* 教程 - 阮一峰的网络日志](https://www.ruanyifeng.com/blog/2017/05/websocket.html)]





## Reference

- [RFC 6202 Bidirectional HTTP](https://datatracker.ietf.org/doc/html/rfc6202#section-2.2)

