# SSLB (Super Simple Load Balance)

Hey, it's a super simples, it's just a toy project that makes me laugh, just to comparing with NGINX here on my OSX El Capitan and dreaming to achieve some kind of performance.

## Configuration file (example)
```
{
    "general": {
        "maxProcs": 4
    },
    "frontends" : [{
        "name" : "Front1",
        "host" : "127.0.0.1",
        "port" : 9000,
        "route" : "/",
        "timeout" : 5000,
        "workerPoolSize": 10,
        "dispatcherPoolSize": 10,
        "backends" : [{
            "name" : "Back1",
            "address" : "http://127.0.0.1:9001",
            "heartbeat" : "http://127.0.0.1:9001",
            "inactiveAfter" : 3,
            "heartbeatTime" : 5000,
            "retryTime" : 5000
        },{
            "name" : "Back2",
            "address" : "http://127.0.0.1:9002",
            "heartbeat" : "http://127.0.0.1:9002",
            "inactiveAfter" : 3,
            "heartbeatTime" : 5000,
            "retryTime" : 5000
        }]
    }]
}
```

* general:
	* maxProcs: Number of processors used
	
* frontends:
	* name: Just a identifier
	* host: Host that serves the HTTP front
	* port: Port that serves the HTTP front
	* route: Route to receive the traffic
	* timeout: How long can wait for the result (ms) from the backend
	* workerPoolSize: Number of workers for processing request
	* dispatcherPoolSize: Number of dispatchers for send the requests to the backends

* backends:
	* name: Just a identifier
	* address: Address (URL) for your backend
	* hearbeat: Addres to send Head request to test if it's ok
	* inactiveAfter: Consider the backend inactive after
	* heartbeatTime: The interval to send a "ping"
	* retryTime: The interval to send a "ping" after the first failed "ping"
	