Go library for ANTLabs devices
==============================

This library provides communication with ANTLabs devices in the Go Language.
Specifically, it implements the API for innGates.  This ANTLabs InnGate API
is capable of mangaing accounts, authenticating users, retrieving information
about plans, and processing billing through the PMS module (if enabled).

Install this package:
---------------------

    go get github.com/secesh/gantlabs

Example:
--------

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

