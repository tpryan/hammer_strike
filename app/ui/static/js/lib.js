// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

function DISTRIBUTOR(hostname){

    this.debug = false;
    var apihostname = hostname;
    this.timeout = 1000;
    var apiprotocol = "https://"

    var uri_report = "/distributor/report";
    var uri_load = "/distributor";
    var uri_list = "/distributor/list";

    var ajaxProxy = function(url, successHandler, errorHandler, timeout) {
        timeout = typeof timeout !== 'undefined' ? timeout : this.timeout;
        $.ajax({
            url: url,
            success: successHandler,
            error: errorHandler,
            timeout: timeout

        });
        if (this.debug){
            console.log("Called: ", url);
        }
    };

    var defaultErrorHandler = function(e){
        console.log("Error in distributor api:", e);
    };

    var getReportURI = function(){
        return apiprotocol + apihostname + uri_report;
    };

    var getListURI = function(){
        return apiprotocol + apihostname + uri_list;
    };

    var getLoadURI = function(){
        return apiprotocol + apihostname + uri_load;
    };

    this.Report = function(token, successHandler){
        ajaxProxy(getReportURI() +"?token=" + token, successHandler,
                    defaultErrorHandler, this.timeout);
    };

    this.List = function(successHandler){
        ajaxProxy(getListURI(), successHandler,
                    defaultErrorHandler, this.timeout);
    };

    this.Load = function(token, count, successHandler){
        console.log(getLoadURI() +"?token=" + token + "&n=" + count);
        ajaxProxy(getLoadURI() +"?token=" + token + "&n=" + count,
                    successHandler, defaultErrorHandler, this.timeout);
    };

}

function GENERATOR(){
    this.Get = function(){
        return (Math.random().toString(36)+'00000000000000000').slice(2, 5+2)
    };
}
