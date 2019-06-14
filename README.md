# CyberGRX API Bulk Connector

This is a general purpose reverse proxy for the [CyberGRX API](https://api.cybergrx.com/v1/swagger/) that eliminates the need to handle pagination from inside a connected application.  

**Notes:**
- This service maintains no internal state and only logs to system out/err.
- You will need to pass a valid API token in the `Authorization` header to gain access to the CyberGRX API.  
- **You are responsible** for placing this service behind a HTTPS gateway, failing to do so may result in the capture of an access token if an attacker is able to sniff the `Authorization` header.
 
**Key features:**
- Self documenting using the formal API swagger specification (navigate to https://HOST:PORT/)
- Pull an entire ecosystem with a sinle API request.
- Retrieve all information for each third party, this includes the latest authorized residual_risk (gaps/findings) and control scores.
- The CyberGRX API is versioned and will maintain backwards compatibility, that will cascade to this service's implementation as well.

## Usage
```
Execute the webservice using GIN_MODE=release ./api-connector-bulk
```

## Documentation

If you are taking a look at configuring the CyberGRX Bulk Connector, please take a look at the [Install Guide (How To)](./HOW-TO.md).

# Development Workflows
- Make sure you have golang >= 1.12
- Make modifications to the source
 - `go build`
 - `CYBERGRX_API=${INTERNAL_API} ./api-connector-bulk`
 - Once you are satisfied update VERSION, commit and push your code
- Make a release
 - `make release`
 - Tag a release in GitHub attaching the versioned binaries (VERSION must match)
- Push docker image
 - `make docker`