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
function HAMMER(max){
    var positive = true;
    var unit = max/500;
    var stopped = false;
    var min = max/20;
    var htmlTweaked = false;


    var setHTML = function(){
        if (!htmlTweaked){
            $("#loadcount").attr("min", min);
            $("#loadcount").attr("max", max);
            $("#loadcount").attr("step", unit);
            htmlTweaked = true;
        }
    };


    this.Move = function(){
        setHTML();
        var slider = document.querySelector("#loadcount");
        var loadcount = parseInt(slider.value);

        if (positive) {
            loadcount += unit;
            slider.value = loadcount;
            if (loadcount >= max) {
                positive = false;
                slider.value = max;
            }
        } else {
            loadcount -= unit;
            slider.value = loadcount;
            if (loadcount <= min) {
                positive = true;
                slider.value = min;
            }
        }

        if (stopped) {
            return;
        }

        if (loadcount > .8 * max) {
            setTimeout(this.Move.bind(this), 1)
        } else if (loadcount > .6 * max) {
            setTimeout(this.Move.bind(this), 4)
        } else if (loadcount > .4 * max) {
            setTimeout(this.Move.bind(this), 6)
        } else if (loadcount > .2 * max) {
            setTimeout(this.Move.bind(this), 8)
        } else {
            setTimeout(this.Move.bind(this), 10)
        }

        var output = document.querySelector("#loadoutput");
        output.value = slider.value;
    };

    this.Stop = function(){
        stopped = true;
    };

    this.Restart = function(){
        stopped = false;
        this.Move();
    };
}

function STATE(){
    var bellrung = false;
    var load_send = 0;
    var interval;

    this.IsBellRung = function(){
        return bellrung;
    };
    this.SetBellRung = function(){
        bellrung = true;;
    };

    this.SetLoadToSend = function(value){
        load_send = value;
    };

    this.GetLoadToSend = function(){
        return parseInt(load_send);
    };

    this.Reset = function(){
        bellrung = false;
        load_send = 0;
        this.EndPoll();
    };

    this.StartPoll = function(pollHandler){
        interval = setInterval(pollHandler, 100);
    };

    this.EndPoll = function(){
        clearInterval(interval);
    };

}

function RESULT(){

    this.SetTotal = function(content){
        $("#totalrequests span").html(content);
    };

    this.ClearHolder = function(){
        $("#gaelist").html("");
    };

    this.ShowResults = function(report){
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
                gaeInstanceDiv.addClass("lifter");
                gaeInstanceDiv.addClass(getRandomBinaryGender());
                $("#gaelist").append(gaeInstanceDiv);
            }

            gaeInstanceDiv.html("<span>" + report.instances[i].requests + "</span>");
        }
    }

    var getRandomBinaryGender = function (){
        var rando = Math.random() * 10;
        if (rando < 5 ){
            return "male";
        }
        return "female";
    };

}

var host = "https://hammer-strike-172817.appspot.com";
var hostname = "hammer-strike-172817.appspot.com";
var max = 10000;
var generator = new GENERATOR()
var token = generator.Get();
var distributor = new DISTRIBUTOR(hostname);
var hammer = new HAMMER(max);
var state = new STATE();
var scoreboard = new RESULT();
var bellstrike = new Audio('audio/bell.wav');
var explain_visible = false;


document.addEventListener('DOMContentLoaded', function() {
    resetGame();
    getGAElist();
    hammer.Move();
    document.querySelector("#flush").addEventListener("click", resetGame);
    document.querySelector("#load").addEventListener("click", handleStrike);
    document.querySelector("#load").focus();
    document.querySelector(".explain").addEventListener("click", toggle_explain);

});


function resetGame(e){
    $("#flush").prop("disabled", true);
    $("#load").prop("disabled", false);
    $(".alert").fadeOut('slow');
    $("#sentrequests span").html(0);
    $("#gaelist").html("");
    $("#gaelist").html("");
    $("#load").focus();

    state.Reset();
    token = generator.Get();
    hammer.Restart();

    getGAElist();

    if (explain_visible){
        toggle_explain();
    }

}

function handleStrike(e){
    sendLoad(document.querySelector("#loadcount").value)
}

function sendLoad(count){
    console.log("Pressed");
    hammer.Stop();
    state.StartPoll(getGAElist);
    state.SetLoadToSend(count);
    distributor.Load(token, count, handleInfo);

    $(".alert").fadeIn('slow');
    $("#gaelist-header").fadeIn('slow');
    $("#flush").prop("disabled", false);
    $("#load").prop("disabled", true);
    $("#sentrequests span").html(count);

    if (explain_visible){
        show_explanation();
    }

}

function getGAElist(){
    console.log("Polled.");
    distributor.Report(token, handleGAElist);
}

function handleGAElist(e){

    if (state.GetLoadToSend() == parseInt(e.request_count)){
        scoreboard.SetTotal("<br />All " + e.request_count);
        state.EndPoll();
    } else {
        scoreboard.SetTotal(e.request_count)
    }

    if (e.request_count == 0){
        scoreboard.ClearHolder();
    }

    scoreboard.ShowResults(e);

    var ratio = e.request_count / max;
    var adj = 255 * ratio;
    document.querySelector("#striker").style.bottom = (90 + adj) + "px";

    if ((adj > 250) && !state.IsBellRung() ) {
        bellstrike.play();
        state.SetBellRung();
    }
}

function handleError(e){
    console.log("Error:", e);
}

function handleInfo(e){
    console.log("Info:", e);
}

function toggle_explain(){
    var explain = $(".explain");
    var explain_text = $(".explain span")
    console.log("Clicked", explain_text.html());
    if (explain_text.html() == "Show Explanation"){
        console.log("Show");
        explain_visible = true;
        show_explanation();
        explain.addClass("explain-show");
        explain_text.html("Hide Explanation");
    } else{
        console.log("hide");
        explain_visible = false;
        hide_explanation();
        explain.removeClass("explain-show");
        explain_text.html("Show Explanation");
    }
}

function show_explanation(){
    $(".overlay-slider").fadeIn('slow');
    if (state.GetLoadToSend > 0){
        $(".overlay-gaelist").fadeIn('slow');
        $(".overlay-alert").fadeIn('slow');
    }
}

function hide_explanation(){
    $(".overlay").fadeOut('slow');
}