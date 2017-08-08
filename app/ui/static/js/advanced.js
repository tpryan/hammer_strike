function RESULT(){

    var loadservers = 0;

    this.CountLoadServers = function(){
        return loadservers;
    }

    this.SetTotal = function(content){
        $("#totalrequests span").html(content);
    };

    this.ClearHolder = function(){
        $("#gaelist").html("");
    };

    this.ShowGAEResults = function(report){
        var seconds = (report.end - report.start) / 1000000000;
        var qps = report.request_count / seconds;

        if (seconds <= 0) {
            seconds = 0;
            qps = 0;
        }
        $("#totalinstances span").html(report.instance_count);
        $("span.seconds").html(Math.round(seconds * 100) / 100);
        $("span.qps").html(Math.round(qps * 100) / 100);

        if (report.instances == null){
            return;
        }


        for (var i = 0; i < report.instances.length; i++){
            var gaeInstanceDiv = $("#"+ report.instances[i].name);

            if (gaeInstanceDiv.length == 0){
                var gaeInstanceDiv  = $("<div></div>");
                gaeInstanceDiv.attr("id", report.instances[i].name);
                gaeInstanceDiv.addClass("gae");
                gaeInstanceDiv.addClass("bigEntrance");
                $("#gaelist").append(gaeInstanceDiv);

            }

            gaeInstanceDiv.html("<span>" + report.instances[i].requests + "</span>");
        }
    }

    this.ShowGKEResults = function(report){
        $("#gkelist").show();

        for (var i = 0; i < report.length; i++){
            var gkeInstanceDiv = $("#"+ report[i].Name);

            if (gkeInstanceDiv.length == 0){
                var gkeInstanceDiv  = $("<a></a>");
                gkeInstanceDiv.attr("id", report[i].name);
                gkeInstanceDiv.addClass("gke");
                gkeInstanceDiv.attr("href", "http://" + report[i].ip + ":30000/log?token=" + token);
                gkeInstanceDiv.attr("target","_new");
                $("#gkelist").append(gkeInstanceDiv);

            }
        }
        loadservers = report.length;
    }


}

function openAllGKE(e){
    $gke = $(".gke");   
    for (var i=0; i < $gke.length; i++){
        window.open($gke[i].href, '_new' + i);
    }

}

function STATE(){
    var load_send = 0;
    var interval;


    this.SetLoadToSend = function(value){
        load_send = value;
    };

    this.ShowLoad = function(){
        return load_send > 0;
    };

    this.GetLoadToSend = function(){
        return parseInt(load_send);
    };

    this.Reset = function(){
        clearInterval(interval);
        load_send = 0;

    };

    this.StartPoll = function(pollHandler){
        interval = setInterval(pollHandler, 100);
    };

    this.SlowPoll = function(pollHandler){
        clearInterval(interval);
        interval = setInterval(pollHandler, 1000);
    };

}

var host = "https://hammer-strike-172817.appspot.com";
var hostname = "hammer-strike-172817.appspot.com";
var generator = new GENERATOR();
var token = generator.Get();
var distributor = new DISTRIBUTOR(hostname);
var scoreboard = new RESULT();
var state = new STATE();


document.addEventListener('DOMContentLoaded', function() {
    $("#loadcount").change(setLoadOutput);
    $("#load").click(handleStrike);
    $("#reset").click(resetView);
    $("#analyze").click(analyizeResults);
    setLoadOutput();
    distributor.List(scoreboard.ShowGKEResults);
    $("#showgke").click(openAllGKE);
});


function setLoadOutput(e){
    var slider = $("#loadcount");
    var loadcount = parseInt(slider.val());
    var output = $("#loadoutput");
    output.val(loadcount);
}


function handleStrike(e){
    sendLoad($("#loadcount").val())
}

function resetView(){
    state.Reset();
    $("#totalinstances span").html(0);
    $("span.seconds").html(0);
    $("span.qps").html(0);
    $(".alert").hide();
    $("#gaelist").html("");
    $("#load").prop("disabled", false);
    $("#reset").prop("disabled", true);
    $("#load_url").html('http://[a totally real app].appspot.com');
    $("#gaelist").hide();
    $("#loaddetails").html("");
    token = generator.Get();

}

function sendLoad(count){
    console.log("Pressed");
    state.StartPoll(getGAElist);
    distributor.Load(token, count, handleInfo);
    state.SetLoadToSend(count);
    $("#load").prop("disabled", true);
    $("#sentrequests span").html(count);
    $("#gkelist div").addClass("floating");
    $("#gkelist div").removeClass("bigEntrance");
    $("#gaelist").show();
    $(".alert").show();
    $("#reset").prop("disabled", false);
    $("#loadtodistribute").html(count);
    divideCount(count);
}

function divideCount(count){
    var loads = scoreboard.CountLoadServers()
    var nodeCount = count / loads;

    for (var i= 0; i < loads; i++){
        var nodeCountDiv  = $("<div>"+nodeCount+"</div>");
        nodeCountDiv.addClass("nodecount");
        $(".gke").html(nodeCountDiv);
    }

}


function getGAElist(){
    console.log("Polled.");
    distributor.Report(token, handleGAElist);
}

function handleGAElist(e){

    if (!state.ShowLoad()){
        return;
    }

    if (state.GetLoadToSend() == parseInt(e.request_count)){
        scoreboard.SetTotal(" All " + e.request_count);
        state.SlowPoll(getGAElist);
        var urlToDisplay = host + "/load/?token=" + token;
        $("#load_url").html('<a href="'+urlToDisplay+'">'+urlToDisplay+'</a>');
        $("#gkelist div").removeClass("floating");
    } else {
        scoreboard.SetTotal(e.request_count)
    }

    if (e.request_count == 0){
        scoreboard.ClearHolder();
    }

    scoreboard.ShowGAEResults(e);

}

function handleInfo(e){
    console.log("Info:", e);
}

function analyizeResults(){
   var total = 0;
   var caches = [];
   var instances = [];
   var others = [];
   $("#gaelist div").each(function(){
       instances.push($(this).attr("id"));
        var list = token + "_instances_" + $(this).attr("id").slice(-1)
        if (caches.indexOf(list) < 0){
            caches.push(list);
        }
   });

   others.push(token + "_start");
   others.push(token + "_end");
   others.push(token + "_total");
   others.push(token + "_totalInstances")
   others.push("LoadNodeList")

   caches.sort();
   console.log(caches)
   console.log(instances)
   var totalKeys = caches.length + instances.length + others.length;
   var totalReport = caches.length + parseInt($("#totalinstances span").html()) + others.length;
   console.log("Memcache count should be:", totalKeys);
   console.log("Total comes back as:", totalReport);




}
