# reverse proxy written in go

## features
- [x] http/https support
- [x] websocket support
- [x] configurable from a file
- [x] chunk encoding
- [x] handle subdomains
- [x] http2 support
- [ ] load balancer to avoid DDoS attacks
- [ ] basic redirect if url not valid
...

## configurable file
``` js
{
	"servers": [
		{
			"port": int,
			"server_name": string,          // must be the domain
			"ssl_to_client": bool,          // activate https from proxy to client
			"ssl_certificate": string,      // path to ssl_certificate file
			"ssl_certificate_key": string,  // path to ssl_certificate_key file
			"max_redirect": int,            // number of max redirect to follow
			"chunk_encoding": bool,         // response will be encoded in chunks (http/1.1)
			"chunk_size": int,              // Kb (min: 8kb)
			"chunk_timeout": int,           // time in ms before chunk is sent (min 30ms)
			"location": [                   // this needs to be an array
				{
					"subdomain": string,    // define which subdomain redirect to "to"
					"to": string            // server to call when subdomain used
				},
			]
		}
	]
}
```
