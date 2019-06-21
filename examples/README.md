All of these sample integrations default to communicating with the bulk connector at http://127.0.0.1:8080.  For local development it is generally safe to run the bulk-api if you do not expose it to the network at large.  The way to do this is to set the following environmental properties for the bulk connector `HOST=127.0.0.1 PORT=8080 ./api-connector-bulk`

If you have the bulk connector running on another system (don't forget to wrap the application with HTTPS) you can set the environmental parameter `CYBERGRX_BULK_API=https://hostname-of-bulk-connector`.

# A set of sample integrations with the bulk connector
- [Bulk export to an excel file](./excel-export/README.md)
