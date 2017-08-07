

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
