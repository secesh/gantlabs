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

// Package antlabs provides methods to communicate with an ANTLabs device.
// http://www.antlabs.com/
package antlabs

import (
	"net/http"
	"crypto/tls"
	// "fmt"
	"strconv"
	"io/ioutil"
	"regexp"
	// "strings"
	"errors"
)

type Host struct{
	Host      string
	Port      int
	ApiPass   string
	AdminPass string
	AdminName string
}

func basicURL(url string) (body []byte, err error){
	tr := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify : true} }
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if(err != nil){return nil, err}
	
	body, err = ioutil.ReadAll(resp.Body)
	if(err != nil){return nil, err}
	defer resp.Body.Close()
	
	return body, nil
}

//InnGateAPIRequest takes a querystring as the only argument and returns a map of values.
//All ANTLabs InnGate API requests work by a very simple webservice.  A URL is crafted
//according to the API to make the proper request.  The result is a plain-text file with
//lines that look like:
//field = value
//when a field has multiple values, they'll be delimited by pipes:
//field = value1|value2|value3|...
//The map of fields to values is produced by processing the body in parseBody().
func (ant *Host) InnGateApiRequest(queryString string) (parsed_body [][]string, err error){
	//we must ignore the certificate because it is self-signed to ezxcess.antlabs.com
	body, err := basicURL("https://"+ant.Host+":"+strconv.Itoa(ant.Port)+"/api/?"+queryString)
	if(err != nil){return nil, err}
	
	
	parsed_body, err = parseApiResponse(string(body))
	if(err != nil){return nil, err}
	
	return parsed_body, nil
}

//parseApiResponse is called by APIRequest and converts the plain-text response from the API into a map of fields and values.
//This is done using a very simple regular expression; data verification is performed elsewhere.
func parseApiResponse(body string) (matches [][]string, err error){
	assigned := regexp.MustCompile(`(?m)^(?:\s*)(\w+?)(?:\s*)=(?:\s*)(.+?)(?:\s*)$`)
	matches   = assigned.FindAllStringSubmatch(body, -1)
	err = nil
	//fmt.Println(body)
	if(len(matches) == 0){return nil, errors.New("Failed to parse body; 0 length")}
	return
}