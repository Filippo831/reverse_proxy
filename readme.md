# reverse proxy written in go

## features
- [x] http/https support
- [x] websocket support
- [x] configurable from a file
- [x] chunk encoding
- [x] handle subdomains
- [x] http2 support
- [x] http3-ish support
- [x] logging (log_file.txt)


# usage
## docker
- clone the repository
```
git clone https://github.com/Filippo831/reverse_proxy.git
```
- generate a locally trusted development certificate via **mkcert**
```
mkcert -cert-file reverse_proxy.com.pem -key-file reverse_proxy.com-key.pem localhost "*.localhost" ::1
```
- create the configuration file *configuration.json*. Write the configuration using the structure below (example after the structure)
```
touch configuration.json
```
### configurable file
``` js
{
    "servers": [
        {
            "port": int,
            "server_name": string,          // must be the domain
            "http3": bool                   // if true use http3 otherwise http1 & http2
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
### example configuration
This is a simple configuration that creates 2 servers, 
``` js
{
  "servers": [
    {
      "port": 8081,
      "ssl_to_client": false,
      "server_name": "localhost",
      "chunk_encoding": true,
      "chunk_size": 64,
      "chunk_timeout": 200,
      "location": [
        {
          "domain": "localhost",
          "to": "http://127.0.0.1:8080"
        }
      ]
    },
    {
      "port": 8082,
      "ssl_to_client": true,
      "ssl_certificate": "reverse_proxy.com.pem",
      "ssl_certificate_key": "reverse_proxy.com-key.pem",
      "server_name": "localhost",
      "chunk_encoding": true,
      "chunk_size": 64,
      "chunk_timeout": 200,
      "location": [
        {
          "domain": "localhost",
          "to": "http://127.0.0.1:8080"
        }
      ]
    }
  ]
}
```
