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

// Package innGateApi implements the ANTLabs InnGate API, a structured
// method of communication with an ANTLabs InnGate. 
// http://www.antlabs.com/
//
// Example:
//   import("antlabs/innGate")
//   func main(){
//     innGate := innGateApi.Host{Host : "ant.example.com"}
//     resp, _ := innGate.ApiVersion()
//     fmt.Println("\n\nAPI_Version:", resp.ApiVersion)
//   }
package innGateApi

type Host struct{
	Host, Pass string
	Port int
}