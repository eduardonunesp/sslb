# SSLB (Super Simple Load Balancer)

It's a Super Simples Load Balancer, just a little project to achieve some kind of performance.

## Install

To install type:

```
go get github.com/eduardonunesp/sslb
```

Don't forget to create your configuration file `config.json` at the same directory of project and run it

```
go run main.go
```

## Configuration file (example)
```
{
    "general": {
        "maxProcs": 4
    },
    
    "frontends" : [
        {
            "name" : "Front1",
            "host" : "127.0.0.1",
            "port" : 9000,
            "route" : "/",
            "timeout" : 5000,
            "workerPoolSize": 10,
            "dispatcherPoolSize": 10,
            
            "backends" : [
                {
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
                }
            ]
        }
    ]
}
```

* general:
	* maxProcs: Number of processors used by Go runtime
	
* frontends:
	* name: Just a identifier to your front server
	* host: Host address that serves the HTTP front
	* port: Port address that serves the HTTP front
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
	

## LICENSE
Copyright (c) 2015, Eduardo Nunes Pereira
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of sslb nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.