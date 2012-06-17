Go library for ANTLabs devices
==============================

This library provides communication with ANTLabs devices in the Go Language.
Specifically, it implements the API for innGates.  The ANTLabs InnGate API
is capable of managing accounts, authenticating users, retrieving information
about plans and processing billing through the PMS module (if enabled).

Install this package:
---------------------

    go get github.com/secesh/gantlabs

Example:
--------

````go
package main

import (
    "fmt"
    "github.com/secesh/gantlabs/innGate"
)
func main(){    
    ant := innGateApi.Host{ Host : "ant.example.com" }
    
    resp, _ := ant.ApiVersion()
    fmt.Println("API_Version:", resp.ApiVersion)
    
    resp2, _ := ant.AccountGetAll(innGateApi.AccountGetAllRequest{Creator : "admin"})
    fmt.Println("Got", resp2.Count, "accounts")

    resp3, _ := ant.AccountDelete(innGateApi.AccountDeleteRequest{Code : []string{"abc123", "def456"}})
    fmt.Println("Deleted", resp3.Deleted, "accounts")
}
````
InnGate API Status:
-------
Below is a list of API modules supported by the ANTLabs InnGate.
  * **Bold** modules are implemented.
  * Non-bold modules are not yet implemented.

Account
  * account_add
  * **account_delete**
  * **account_get**
  * **account_get_all**
  * account_update
  
API
  * **api_module**
  * **api_modules**
  * api_password_get (only available to native PHP API)
  * **api_version**
  
Authentication
  * **auth_authenticate** (not tested)
  * auth_init
  * **auth_login** (not tested)
  * auth_logout
  * auth_update
  * sid_get
  * publicip_get
  
Plan
  * plan_get_all
  * plan_get_id
  
Data
  * data_get
  * data_set
  * data_get_keys
  * data_get_names
  * data_delete
  
Property Management System (PMS)
  * pms_billing_log
  * pms_guest_status
  * pms_post_check
  * pms_post
  * pms_room_status
  
Network
  * vlan_get
  * vlan_update
  * device_status
  
Credit Card
  * cc_payflowpro_post
  
Miscellaneous
  * browser

Contact:
--------

gantlabs@chasefox.net