//  Copyright 2012 ChaseFox (Matthew R Chase)
//  
//  This file is part of gantlabs, a go library for communicating with
//  ANTLabs devices. http://www.antlabs.com/
//  
//  gantlabs is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published
//  by the Free Software Foundation, either version 3 of the License,
//  or (at your option) any later version.
//  
//  gantlabs is distributed in the hope that it will be useful, but
//  WITHOUT ANY WARRANTY; without even the implied warranty of 
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//  
//  You should have received a copy of the GNU General Public License
//  along with gantlabs.  If not, see <http://www.gnu.org/licenses/>.

package innGateApi

import (
	"antlabs"
	//"fmt"
	"strconv"
	"regexp"
	"strings"
	"errors"
	"time"
	"html"
)


type Account struct{
	Type    string
	Creator string
	UserId  string
	Code    string
	Description string
	Enable      bool
	ValidFrom   time.Time
	ValidUntil  time.Time
	LoginLimit  bool
	LoginMax, LoginCount, SharingMax int64
	UserGroupName string
	CreateTime    string
	UpdateTime    string
	Accounting    string
	BillingId     string
}
type responseCommon struct{
	Op            string
	Result        string
	Resultcode    int64
	Error         string
	ModuleVersion float64
}

type requestCommon struct{
	op         string
}
//These structs are used when op=account_get
type accountGetResponse struct{
	responseCommon
	UserId       []string
	Code         []string
	SharingIndex []int64
	ClientMac    []string
	Description  []string
	Enabled      []bool
	ValidFrom    []time.Time
	ValidUntil   []time.Time
	LoginLimit   []bool
	LoginMax     []int64
	LoginCount   []int64
	SharingMax   []int64
	Plan         []string
	DurationBalance []string
	VolumeBalance   []string
	CreateTime      []time.Time
	UpdateTime      []time.Time
}
type AccountGetRequest struct{
	requestCommon
	UserId, Code, ClientMac string
}
//These structs are used when op=account_get_all
type accountGetAllResponse struct{
	responseCommon
	Count int64
	Header []string
	Accounts []Account
}
type AccountGetAllRequest struct{
	requestCommon
	ValidFromStart, ValidFromEnd, ValidUntilStart, ValidUntilEnd time.Time
	Creator, Description, Type string
	CreatedStart, CreatedEnd, PlanName string
}

//These structs are used when op=account_delete
type accountDeleteResponse struct{
	responseCommon
	Deleted int64
}
type AccountDeleteRequest struct{
	requestCommon
	UserId, Code interface{}
}
//These structs are used when op=api_module
type moduleResponse struct{
	responseCommon
	Version float64
}
type ModuleRequest struct{
	requestCommon
	Module string
}

//These structs are used when op=api_modules
type modulesResponse struct{
	responseCommon
	Count   int64
	Modules map[string]float64
}
type modulesRequest struct{
	requestCommon
}

//These structs are used when op=api_version
type versionResponse struct{
	responseCommon
	ApiVersion float64
}
type versionRequest struct{
	requestCommon
}


//findCommoners scans the field:value map for API elements common to every API response.
func (common *responseCommon) findCommoners(parsed_body [][]string) (err error){
	for _, v := range parsed_body{
		switch v[1]{
		case "op":
			common.Op = v[2]
		case "result":
			common.Result = v[2]
		case "resultcode":
			common.Resultcode, err = strconv.ParseInt(v[2], 10, 64)
			if(err != nil){ return }
		case "error":
			common.Error = v[2]
		case "version":
			//This will not be valid if the op=api_module.  The API says it
			//assigns version twice, but in practice it only sends version once.
			//The version sent by the API refers to the module about which the
			//inquiry is made (not of api_module itself)
			common.ModuleVersion, err = strconv.ParseFloat(v[2], 64)
			if(err != nil){ return }
		}
	}
	
	if(len(common.Error)    >0){err = errors.New("Error: " + common.Error)}
	if(len(common.Op)      ==0){err = errors.New("Missing expected field in reply (op).")}
	if(len(common.Result)  ==0){err = errors.New("Missing expected field in reply (result).")}
	if(common.Resultcode   !=0){err = errors.New("Resultcode is not OK (" + strconv.FormatInt(common.Resultcode, 10) + ").")}
    if(common.ModuleVersion==0){err = errors.New("Missing expected field in reply (version).")}
	
	return
}



//  Module performs an API request for op=api_module
//  
//  This method requires one argument of type innGateApi.ModuleRequest.
//  See the ANTLabs API for more information regarding elements of the argument.  
//  The example below demonstrates how to send a request.
//
//Example: 
//  ant := innGateApi.Host{ 
// 	   Host : "ant.example.com", //can be an IP or hostname
//  }
//  resp, err := ant.Module(innGateApi.ModuleRequest{Module : "api_modules"})
//  if(err != nil){ panic(err) }
//  fmt.Println("\n\nAPI Module:", resp.Version)
func (api *Host) Module(request ModuleRequest) (result *moduleResponse, err error){
	ant := api.ant()
	request.op     = "api_module"
	result         = &moduleResponse{}
	
	parsed_body, err := ant.InnGateApiRequest("api_password="+api.Pass+"&op="+request.op+"&module="+request.Module)
	if( err != nil){ return nil, err }
	
	err = result.findCommoners(parsed_body)
	//if( err != nil){ return nil, err }
	
	for _, v := range parsed_body{
		switch v[1]{
		case "version":
			result.Version, err = strconv.ParseFloat(v[2], 64)
			if(err != nil){ return nil, err }
		}
	}
	
	return result, nil
}

//  Modules performs an API request for op=api_modules
//  
//  This module does not require or accept any arguments.
//
//Example: 
//  ant := innGateApi.Host{ 
// 	   Host : "ant.example.com", //can be an IP or hostname
//  }
//  resp, err := ant.Modules()
//  if(err != nil){ panic(err) }
//  fmt.Println("\n\nAPI_Version:", resp.ApiModules)
func (api *Host) Modules() (result *modulesResponse, err error){
	ant := api.ant()
	request     := modulesRequest{}
	request.op   = "api_modules" 
	
	result = &modulesResponse{}
	result.Modules = make(map[string]float64) //initalize the map so we can assign values in it later.
	
	parsed_body, err := ant.InnGateApiRequest("api_password="+api.Pass+"&op="+request.op)
	if( err != nil){ return nil, err }
	
	err = result.findCommoners(parsed_body)
	//if( err != nil){ return nil, err }
	
	for _, v := range parsed_body{
		switch v[1]{
		case "count":
			result.Count, err = strconv.ParseInt(v[2], 10, 64)
			if(err != nil){ return nil, err }
	 	case "modules":
	 		modules := strings.Split(v[2], "|")
	 		for _, v := range modules{
	 			module := strings.Split(v, " ")
	 			moduleName         := module[0]
	 			moduleVersion, err := strconv.ParseFloat(module[1], 64)
	 			if(err != nil){ return nil, err }
	 			
	 			result.Modules[moduleName] = moduleVersion
	 		}
		}
	}
	
	return result, nil
}

//  AccountGet performs an API request for op=account_get
//  
//  This method requires one argument of type innGateApi.AccountGetRequest.
//  See the ANTLabs API for more information regarding elements of the argument.  
//  The example below demonstrates how to send a request.
//  
//Example: 
//  ant := innGateApi.Host{ 
// 	   Host : "ant.example.com", //can be an IP or hostname
//  }
//  resp, err := ant.AccountGet(innGateApi.AccountGetRequest{Code : "abc123"})
//  if(err != nil){ panic(err) }
//  fmt.Println("\n\nAccount:", resp)
func (api *Host) AccountGet(request AccountGetRequest) (result *accountGetResponse, err error){
	ant := api.ant()
	request.op   = "account_get"
	result       = &accountGetResponse{}
	
	query := "api_password="+api.Pass+"&op="+request.op
	if(request.Code != ""){ query += "&code=" + html.EscapeString(request.Code)}
	if(request.UserId != ""){ query += "&userid=" + html.EscapeString(request.UserId)}
	if(request.ClientMac != ""){ query += "&client_mac=" + html.EscapeString(request.ClientMac)}
	
	parsed_body, err := ant.InnGateApiRequest(query)
	if( err != nil){ return nil, err }
	
	err = result.findCommoners(parsed_body)
	//if( err != nil){ return nil, err }
	
	for _, v := range parsed_body{
		switch v[1]{
		case "userid":
			result.UserId = strings.Split(v[2], "|")
		case "code": 
			result.Code = strings.Split(v[2], "|")
		case "sharing_index":
			for _,v := range strings.Split(v[2], "|"){
				sharingIndex, err := strconv.ParseInt(v, 10, 64)
				if(err != nil){ return nil, err }
				result.SharingIndex = append(result.SharingIndex, sharingIndex)
			}
		case "client_mac":
			result.Code = strings.Split(v[2], "|")
		case "description":
			result.Description = strings.Split(v[2], "|")
		case "enabled":
			for _,v := range strings.Split(v[2], "|"){
				if(v == "yes"){
					result.Enabled = append(result.Enabled, true)
				}else{
					result.Enabled = append(result.Enabled, true)
				}
			}
		case "valid_from":
			for _,v := range strings.Split(v[2], "|"){
				validFrom, err := time.Parse(time.RFC1123Z, v)
				if(err != nil){ return nil, err }
				result.ValidFrom = append(result.ValidFrom, validFrom)
			}
		case "valid_until":
			for _,v := range strings.Split(v[2], "|"){
				validUntil, err := time.Parse(time.RFC1123Z, v)
				if(err != nil){ return nil, err }
				result.ValidUntil = append(result.ValidUntil, validUntil)
			}
		case "login_limit":
			for _,v := range strings.Split(v[2], "|"){
				if(v == "on"){
					result.LoginLimit = append(result.LoginLimit, true)
				}else{
					result.LoginLimit = append(result.LoginLimit, false)
				}
			}
		case "login_max":
			for _,v := range strings.Split(v[2], "|"){
				loginMax, err := strconv.ParseInt(v, 10, 64)
				if(err != nil){ return nil, err }
				result.LoginMax = append(result.LoginMax, loginMax)
			}
		case "login_count":
			for _,v := range strings.Split(v[2], "|"){
				loginCount, err := strconv.ParseInt(v, 10, 64)
				if(err != nil){ return nil, err }
				result.LoginCount = append(result.LoginCount, loginCount)
			}
		case "sharing_max":
			for _,v := range strings.Split(v[2], "|"){
				sharingMax, err := strconv.ParseInt(v, 10, 64)
				if(err != nil){ return nil, err }
				result.SharingMax = append(result.SharingMax, sharingMax)
			}
		case "plan":
			result.Plan = strings.Split(v[2], "|")
		case "duration_balance":
			result.DurationBalance = strings.Split(v[2], "|")
		case "volume_balance":
			result.VolumeBalance = strings.Split(v[2], "|")
		case "create_time":
			for _,v := range strings.Split(v[2], "|"){
				createTime, err := time.Parse(time.RFC1123Z, v)
				if(err != nil){ return nil, err }
				result.CreateTime = append(result.CreateTime, createTime)
			}
		case "update_time":
			for _,v := range strings.Split(v[2], "|"){
				updateTime, err := time.Parse(time.RFC1123Z, v)
				if(err != nil){ return nil, err }
				result.UpdateTime = append(result.UpdateTime, updateTime)
			}
		}
	}
	
	return result, nil
}

//  AccountGetAll performs an API request for op=account_get_all
//  
//  This method requires one argument of type innGateApi.AccountGetAllRequest.
//  See the ANTLabs API for more information regarding elements of the argument.  
//  The example below demonstrates how to send a request.
//
//Example: 
//   ant := innGateApi.Host{ 
// 	   Host : "ant.example.com", //can be an IP or hostname
//   }
//   resp, err := ant.AccountGetAll(nil)
//   if(err != nil){ panic(err) }
//   fmt.Println("\n\nAccounts (", resp.Count, "):\n", resp.Header)
//   fmt.Println(resp.Accounts)
//   
//   resp2, err := ant.AccountGetAll(innGateApi.AccountGetAllRequest{Creator : "admin"})
//   if(err != nil){ panic(err) }
//   fmt.Println("\n\n", resp2.Count, "accounts created by admin.")
//   fmt.Println(resp2.Accounts)
//
//NOTICE:
//   If you submit something that returns an empty result (like specifying a "creator"
//   that didn't create any accounts), the API might return an error 90.  This is a
//   bug in the API that has not been worked-around in this package.
func (api *Host) AccountGetAll(arg interface{}) (result *accountGetAllResponse, err error){
	ant := api.ant()
	request     := AccountGetAllRequest{}
	request, _   = arg.(AccountGetAllRequest) //fail silently in case we got sent a nil.  Otherwise assume we got a good argument.
	request.op   = "account_get_all" 
	result       = &accountGetAllResponse{}
	
	query := "api_password="+api.Pass+"&op="+request.op
	if(len(request.Creator)>0){ query += "&creator=" + html.EscapeString(request.Creator)}
	if(len(request.Type)>0){query += "&type=" + html.EscapeString(request.Type)}
	if(request.ValidFromStart  != time.Time{}){query += "&type=valid_from_start"  + strconv.FormatInt(request.ValidFromStart.Unix(), 10)}
	if(request.ValidFromEnd    != time.Time{}){query += "&type=valid_from_end"    + strconv.FormatInt(request.ValidFromEnd.Unix(), 10)}
	if(request.ValidUntilStart != time.Time{}){query += "&type=valid_until_start" + strconv.FormatInt(request.ValidUntilStart.Unix(), 10)}
	if(request.ValidUntilEnd   != time.Time{}){query += "&type=valid_until_end"   + strconv.FormatInt(request.ValidUntilEnd.Unix(), 10)}
	if(len(request.Description)>0){query += "&description=" + html.EscapeString(request.Description)}
	if(len(request.CreatedStart)>0){query += "&created_start=" + html.EscapeString(request.CreatedStart)}
	if(len(request.CreatedEnd)>0){query += "&created_end=" + html.EscapeString(request.CreatedEnd)}
	if(len(request.PlanName)>0){query += "&plan_name=" + html.EscapeString(request.PlanName)}
	
	parsed_body, err := ant.InnGateApiRequest(query)
	if( err != nil){ return nil, err }
	
	err = result.findCommoners(parsed_body)
	//if( err != nil){ return nil, err }
	
	recordIdentifier := regexp.MustCompile(`record_(\d+)`)
	records          := make([]Account, 0, 0)
	for _, v := range parsed_body{
		switch {
		case v[1] == "header":
			//fmt.Println(v[2])
			result.Header = strings.Split(v[2], "|")
		case recordIdentifier.MatchString(v[1]):
			account := Account{}
			line    := strings.Split(v[2], "|")
			if(len(line) != 17){ return nil, errors.New("Unknown account information (unexpected array length" + strconv.Itoa(len(line)) +").")}
			account.Type          = line[ 0]
			account.Creator       = line[ 1]
			account.UserId        = line[ 2]
			account.Code          = line[ 3]
			account.Description   = line[ 4]
			
			switch line[5]{
			case "yes":
				account.Enable = true
			default :
				account.Enable = false
			}
			
			validFrom, err := strconv.ParseInt(line[ 6], 10, 64)
			if(err != nil){ return nil, err }
			account.ValidFrom = time.Unix(validFrom, 0)
			
			validUntil, err := strconv.ParseInt(line[ 7], 10, 64)
			if(err != nil){ return nil, err }
			account.ValidUntil = time.Unix(validUntil, 0)
			
	        switch line[8]{
	        case "on":
	        	account.LoginLimit = true
	        default:
	        	account.LoginLimit = false
	        }
			
			account.LoginMax, err = strconv.ParseInt(line[ 9], 10, 64)
			if(err != nil){ return nil, err }
			
			account.LoginCount, err = strconv.ParseInt(line[10], 10, 64)
			if(err != nil){ return nil, err }
			
			account.SharingMax, err = strconv.ParseInt(line[11], 10, 64)
			if(err != nil){ return nil, err }
			
			account.UserGroupName = line[12]
			
			account.CreateTime    = line[13]
			account.UpdateTime    = line[14]
			account.Accounting    = line[15]
			account.BillingId     = line[16]
			
			records = append(records, account)
	 	case v[1] == "count":
			result.Count, err = strconv.ParseInt(v[2], 10, 64)
			if(err != nil){ return nil, err }
		default:
			//fmt.Println(v[1])
		}
	}
	result.Accounts = records
	return result, nil
}

//  AccountDelete performs the an API request for op=account_delete
//  
//  This method requires one argument of type innGateApi.AccountDeleteRequest.
//  See the ANTLabs API for more information regarding elements of the argument.  
//  The example below demonstrates how to send a request.
//
//Example: 
//  ant := innGateApi.Host{ 
// 	   Host : "ant.example.com", //can be an IP or hostname
//  }
//  resp, err := ant.AccountDelete(innGateApi.AccountDeleteRequest{Code : "abc123"})
//  if(err != nil){ panic(err) }
//  fmt.Println("\n\nDeleted:", resp.Deleted)
//
//  resp2, err := ant.AccountDelete(innGateApi.AccountDeleteRequest{Code : []string{"abc123", "xyz456"}})
//  if(err != nil){ panic(err) }
//  fmt.Println("\n\nDeleted:", resp2.Deleted)
//
//NOTICE:
//  Database errors can be common.  If the API cannot find any of the accounts requested for deletion in
//  the database, it will return an error with resultcode 98 (database error).  If it finds at least one
//  match, the request should reply with success.  Furthermore, the ANTLabs database/API seems bugarrific; 
//  frequently an account can be seen through the admin portal, but not found when making an API request.
//  If a database error occurs with the API, that result will be passed along.
func (api *Host) AccountDelete(request AccountDeleteRequest) (result *accountDeleteResponse, err error){
	ant := api.ant()
	request.op = "account_delete" 
	result     = &accountDeleteResponse{}
	
	query := "api_password="+api.Pass+"&op="+request.op
	switch request.Code.(type){
	case string:
		if(request.Code != ""){query += "&code=" + request.Code.(string)}
	case []string:
		if(len(request.Code.([]string)) > 0){query += "&code=" + strings.Join(request.Code.([]string), "|")}
	}
	switch request.UserId.(type){
	case string:
		if(request.UserId != ""){query += "&userid=" + request.UserId.(string)}
	case []string:
		if(len(request.UserId.([]string)) > 0){query += "&userid=" + strings.Join(request.UserId.([]string), "|")}
	}
	
	parsed_body, err := ant.InnGateApiRequest(query)
	if( err != nil){ return nil, err }
	
	err = result.findCommoners(parsed_body)
	//if( err != nil){ return nil, err }
	
	for _, v := range parsed_body{
		switch v[1]{
	 	case "deleted":
			result.Deleted, err = strconv.ParseInt(v[2], 10, 64)
			if(err != nil){ return nil, err }
		}
	}
	
	return result, nil
}

//  ApiVersion performs the an API request for op=api_version
//  
//  No optional or required arguments.
//
//Example: 
//  ant := innGateApi.Host{ 
// 	   Host : "ant.example.com", //can be an IP or hostname
//  }
//  resp, err := ant.ApiVersion()
//  if(err != nil){ panic(err) }
//  fmt.Println("\n\nAPI_Version:", resp.ApiVersion)
//
//  we understand it might be misleading to call this method "ApiVersion"
//  instead of just "Version."  But this comes from the op, and we chose
//  to keep it in long form so the result.ApiVersion is distinct from
//  the common version (of the op, not the API).
func (api *Host) ApiVersion() (result *versionResponse, err error){
	ant := api.ant()
	request     := versionRequest{}
	request.op   = "api_version" 
	result       = &versionResponse{}
	
	parsed_body, err := ant.InnGateApiRequest("api_password="+api.Pass+"&op="+request.op)
	if( err != nil){ return nil, err }
	
	err = result.findCommoners(parsed_body)
	//if( err != nil){ return nil, err }
	
	for _, v := range parsed_body{
		switch v[1]{
	 	case "api_version":
			result.ApiVersion, err = strconv.ParseFloat(v[2], 64)
			if(err != nil){ return nil, err }
		}
	}
	
	return result, nil
}

//ant = api.ant() is called at the beginning of every API method, which
//provides invisible glue between the innGateApi package and the antlabs package.
func (api *Host) ant() (ant *antlabs.Host){
	//Let's start by defining defaults in case these arguments weren't specified.
	port := api.Port; if(port == 0 ){ port = 443; api.Port = port }
	pass := api.Pass; if(pass == ""){ pass = "admin"; api.Pass = pass }
	
	//And now we prepare our return.
	ant = &antlabs.Host{
		Host    : api.Host,
		Port    : port,
		ApiPass : pass,
	}
	return
}