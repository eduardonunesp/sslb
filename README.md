# SSLB (Super Simple Load Balancer) ver 0.1.0

It's a Super Simples Load Balancer, just a little project to achieve some kind of performance.

## Features
 * High availability (improving with time the speed)
 * Support to WebSockets
 * Monitoring the internal state (improving)
 * Really easy to configure, just a little JSON file

## Next features
 * Manage configurations in runtime without downtime
 * Complete internal status and diagnostics
 * HTTP/2 support
 * Cache 
 * HTTPS support
 
 If you have any suggestion don't hesitate to open an issue, pull requests are welcome too.

## Install

To install type:

```
go get github.com/eduardonunesp/sslb
```

Don't forget to create your configuration file `config.json` at the same directory of project and run it. You can use the command `sslb -c` to create an example of configuration file.


## Usage
Type `sslb -h` for the command line help

```
sslb -h                                                                                                                                                              
NAME:
   SSLB (github.com/eduardonunesp/sslb) - sslb

USAGE:
   sslb [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   status, s	Return the internal status
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, -b	activate the verbose output
   --filename, -f 	set the filename as the configuration
   --help, -h		show help
   --version, -v	print the version
```

After the configuration file completed you can type only `sslb -b` to start SSLB with verbose mode, that command will log the output from SSLB in console. That will print something like that:

```
sslb -b                                                                                                                                                               
2015/10/25 22:58:33 Start SSLB (Server)
2015/10/25 22:58:33 Create worker pool with [1000]
2015/10/25 22:58:33 Prepare to run server ...
2015/10/25 22:58:33 Setup and check configuration
2015/10/25 22:58:33 Setup ok ...
2015/10/25 22:58:33 Run frontend server [Front1] at [0.0.0.0:80]
2015/10/25 22:58:34 Backend active [Backend 1]
2015/10/25 22:58:34 Backend active [Backend 2]
2015/10/25 22:58:34 Backend active [Backend 3]
```

## Configuration options

* general:
	* maxProcs: Number of processors used by Go runtime (default: Number of CPUS)
	* workerPoolSize: Number of workers for processing request (default: 10)
	* gracefulShutdown: Wait for the last connection closed, before shutdown (default: true)
	* websocket: Ready for respond websocket connections (default: true)
	* rpchost: Address to expose the internal state (default: 127.0.0.1)
	* rpcport: Port to expose the internal state (default: 42555)
	
* frontends:
	* name: Just a identifier to your front server (required)
	* host: Host address that serves the HTTP front (required)
	* port: Port address that serves the HTTP front (required)
	* route: Route to receive the traffic (required)
	* timeout: How long can wait for the result (ms) from the backend (default: 30000ms)

* backends:
	* name: Just a identifier (required)
	* address: Address (URL) for your backend (required)
	* hearbeat: Addres to send Head request to test if it's ok (required)
	* hbmethod: Method used in request to check the heartbeat (default: HEAD)
	* inactiveAfter: Consider the backend inactive after the number of checks (default: 3)
	* activeAfter: COnsider the backend active after the number of checks (default: 1)
	* heartbeatTime: The interval to send a "ping" (default: 30000ms)
	* retryTime: The interval to send a "ping" after the first failed "ping" (default: 5000ms)
	
### Example (config.json)

```
{
    "general": {
        "maxProcs": 4,
        "workerPoolSize": 10,
    },
    
    "frontends" : [
        {
            "name" : "Front1",
            "host" : "127.0.0.1",
            "port" : 9000,
            "route" : "/",
            "timeout" : 5000,
            
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