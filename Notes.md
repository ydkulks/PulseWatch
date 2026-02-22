# Note

## TODO

- [x] Watch a process
- [ ] Watch multiple processes
- [ ] Notification
    - [ ] Discord
    - [ ] Email
- [ ] Front-end

## gRPC type for monitoring

Initially we think `server-side streaming` is the right and obvious choice for
monitoring. Thats how a lot of pub/sub or event systems feel at glance.

|Aspect                         |Server streaming (client → many from server)                                                                                        |Bidirectional streaming (both directions)                                                                                            |Winner for monitoring/alerting      |
|---                            |---                                                                                                                                 |---                                                                                                                                  |---                                 |
| Heartbeat sending             |Client must open new stream every ~30–60 s (or use separate unary RPCs)                                                             |Client sends heartbeats continuously on same long-lived stream                                                                       |Bidirectional                       |
| Detecting client death quickly|Server only notices when it tries to push and the stream is already broken (can take 10–90+ seconds depending on keepalive settings)|Server sees client stop reading / sending heartbeats immediately (or very quickly with gRPC keepalive pings)                         |Bidirectional                       |
| Sending commands back         |Not possible (stream is one-way from server)                                                                                        |"Server can push: rate-limit config, "please restart", "change heartbeat interval", "ack last error", emergency shutdown signal, etc.|Bidirectional                       |
| Reconnection & session state  |Every new stream is a new identity → need extra logic to correlate                                                                  |One long-lived stream = natural session. Easy to tie auth, process metadata, etc. to the stream itself                               |Bidirectional                       |
| Flow control & backpressure   |Server can be overwhelmed if many dead clients keep streams open                                                                    |Client can stop reading → server naturally pauses sending (built-in HTTP/2 flow control)                                             |Bidirectional                       |
| Implementation complexity     |Looks simpler at first, but you end up needing extra unary RPCs for heartbeats/config → fragmented                                  |One RPC method handles registration + heartbeats + alerts + control messages                                                         |Bidirectional (after learning curve)|
