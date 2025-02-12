# AdGuardHomeStats
A basic statistics proxy for [AdGuard Home](https://github.com/AdguardTeam/AdGuardHome)

Written in Go with zero external dependencies, and is intended to run on the same machine as the AdGuard Home instance. It uses basic authentication to serve a select number of statistics through a proxy.

Currently hosts a single endpoint (`/control/stats`) which contains select fields from the full AdGuard Home response.

## Example `/control/stats` response
```json
{
  "time_units": "hours",
  "num_dns_queries": 5629,
  "num_blocked_filtering": 1451,
  "avg_processing_time": 0.030035
}
```

## Homepage Usage
This proxy can be used as a [Homepage](https://github.com/gethomepage/homepage) target by using the existing [AdGuard Home Widget](https://gethomepage.dev/widgets/services/adguard-home/). Simply exclude the username and password field, while pointing to the proxy's host target.

## Usage
Supports the following arguments:  
```
  -host string  
        Host the proxy will bind with (default "127.0.0.1")  
  -p string  
        Password for authentication (required)  
  -port uint  
        Port the proxy will bind with (default 3001)  
  -u string  
        Username for authentication (required)  
  -url string  
        Base URL to request (default "http://127.0.0.1")  
```
