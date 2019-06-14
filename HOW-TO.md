# Configuration of the bulk connector
There are 3 steps to successful bulk connector integration.  
1. Retrieve a CyberGRX API Token.
1. Download and configure the bulk-connector binary or docker image.
1. Protect the bulk-connector by reverse proxying requests from a HTTPS enabled gateway

## Retrieve a CyberGRX API Token
You must have an active CyberGRX account that has the ability to manage users. 

The workflow is pretty simple:
 - Enter `Manage my company user accounts` using the top right navigational element
 - Click the tab to `Manage Access Tokens`
 - Click the button to `Add a new token`
 - Accept the promt creating a token
 - Show the secret and copy it, this is the only time that this token will be available over the API!
 - Close the dialog and configure your token in the Splunk connector

 ### CyberGRX Management Workflow
 ![enter-user-management]
 ![add-a-token]
 ![confirm-new-token]
 ![view-token]
 ![copy-secret]


[enter-user-management]: /docs/enter-user-management.png "Click top right icon and enter `Manage my company user accounts`"

[add-a-token]: /docs/add-a-token.png "Click the tab to `Manage Access Tokens` and Add a new token"

[confirm-new-token]: /docs/confirm-new-token.png "Accept the promt creating a token"

[view-token]: /docs/make-sure-you-view.png "Before leaving view the token secret"

[copy-secret]: /docs/copy-secret.png "Show the secret and copy it, this is the only time that this token will be available over the API!"