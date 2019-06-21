All of these sample integrations default to communicating with the bulk connector at http://127.0.0.1:8080.  For local development it is generally safe to run the bulk-api if you do not expose it to the network at large.  In this context; `generally safe` means that the only way for an attacker to sniff your Authorization token is to be activley monitoring network communications on your local system.  Acknowlege that this is a security risk but chances of exploit are low.  The safest best practice is to experiment with the API locally and then delete the token from the [CyberGRX token management](../HOW-TO.md#cybergrx-management-workflow) page.

The way to run the bulk connector in this mode of operation is to set the following environmental properties `HOST=127.0.0.1 PORT=8080 GIN_MODE=debug ./api-connector-bulk`.  

# Interface with production bulk connector
If you have the bulk connector running on a production system (don't forget to wrap the application with HTTPS).  You can run an example application and set the environmental parameter `CYBERGRX_BULK_API=https://hostname-of-bulk-connector`.

# A set of sample integrations with the bulk connector
- [Bulk export to an excel file](./excel-export/README.md)
