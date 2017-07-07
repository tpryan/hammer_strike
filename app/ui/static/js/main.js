var host = "https://hammer-strike-172817.appspot.com";
var token = generateToken();
var gaelist_endpoint = host + "/distributor/report";
var load_creator = host + "/distributor";
var bellstrike = new Audio('audio/bell.wav');
var bellrung = false;
var positive = true;
var load_send = 0;
var interval;
var explain_visible = false;

// config should be set in firebaseconfig.js
// firebase.initializeApp(firebaseConfig);
// var firebaseRef = firebase.database().ref('hammer-strike/strength');
// var firebaseInit = false;

document.addEventListener('DOMContentLoaded', function() {
    sendFlush();
    getGAElist();
    moveHammer();
    document.querySelector("#flush").addEventListener("click", sendFlush);
    document.querySelector("#load").addEventListener("click", handleStrike);
    document.querySelector("#load").focus();
    document.querySelector(".explain").addEventListener("click", toggle_explain);
    // firebaseRef.on("value", function(snap){
    //     if(firebaseInit){
    //         sendLoad(snap.val());
    //         document.querySelector("#loadcount").value = snap.val();
    //         document.querySelector("#loadoutput").value = snap.val();
    //     }
    //     firebaseInit = true;
    // });
});

function moveHammer(){
    var unit = 20;
    var slider = document.querySelector("#loadcount");
    var loadcount = parseInt(slider.value);


    if (positive){
        loadcount += unit;
        slider.value = loadcount;
        if (loadcount >= 10000){
            positive = false;
            slider.value = 10000;
        }
    } else{
        loadcount -= unit;
        slider.value = loadcount;
        if (loadcount <= 500){
            positive = true;
            slider.value = 500;
        }
    }

    if (load_send > 0){
        return;
    }

    if (loadcount > 8000){
        setTimeout(moveHammer, 1)
    } else if (loadcount > 6000){
        setTimeout(moveHammer, 4)
    } else if (loadcount > 4000){
        setTimeout(moveHammer, 6)
    } else if (loadcount > 2000){
        setTimeout(moveHammer, 8)
    } else {
        setTimeout(moveHammer, 10)
    }
    var output = document.querySelector("#loadoutput");
    output.value = slider.value;
}

function sendFlush(e){
    document.querySelector("#flush").disabled = true;
    document.querySelector("#load").disabled = false;
    clearInterval(interval);
    $(".alert").fadeOut('slow');
    load_send = 0;
    bellrung = false;
    document.querySelector("#sentrequests span").innerHTML = 0;
    token = generateToken();
    document.querySelector("#gaelist").innerHTML="";
    moveHammer();
    document.querySelector("#gaelist").innerHTML="";
    document.querySelector("#load").focus();
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
    interval = setInterval(getGAElist, 100);
    document.querySelector("#flush").disabled = false;
    document.querySelector("#load").disabled = true;
    $(".alert").fadeIn('slow');
    $("#gaelist-header").fadeIn('slow');
    load_send = count;

    if (explain_visible){
        show_explanation();
    }

    document.querySelector("#sentrequests span").innerHTML = load_send;
    console.log(load_send);
    console.log(load_creator + "?token=" + token + "&n=" +  load_send);
    $.ajax({
        url: load_creator + "?token=" + token + "&n=" +  load_send ,
        success: log,
        error: handleError

    });
}


function getGAElist(){
    console.log("Polled.");
    $.ajax({
        url: gaelist_endpoint + "?token=" + token,
        success: handleGAElist,
        error: handleError

    });
}


function handleGAElist(e){
    var loadcount = document.querySelector("#loadcount").value;
    var holder = document.querySelector("#gaelist");
    var totalHolder = document.querySelector("#totalrequests span");
    var instanceHolder = document.querySelector("#totalinstances span");
    var secondsHolder = document.querySelector("span.seconds");
    var qpsHolder = document.querySelector("span.qps");

    if (parseInt(load_send) == parseInt(e.request_count)){
        totalHolder.innerHTML = "<br />All " + e.request_count;
    } else {
        totalHolder.innerHTML = e.request_count;
    }

    if (e.request_count == 0){
        holder.innerHTML = "";
    }

    seconds = (e.end - e.start) / 1000000000;
    qps = e.request_count / seconds;

    if (seconds <= 0) {
        seconds = 0;
        qps = 0;
    }

    instanceHolder.innerHTML = e.instance_count;
    secondsHolder.innerHTML = Math.round(seconds * 100) / 100;
    qpsHolder.innerHTML = Math.round(qps * 100) / 100;;

    var ratio = e.request_count / 10000;
    var adj = 255 * ratio;
    document.querySelector("#striker").style.bottom = (90 + adj) + "px";

    if ((adj > 250) && !bellrung ) {
        bellstrike.play();
        bellrung = true;
    }

    if (e.instances == null){
        return;
    }
    for (var i = 0; i < e.instances.length; i++){
        var div = document.getElementById(e.instances[i].name)

        if (!div){

            var cssClass = "female";
            var rando = getRandomArbitrary(0, 10)
            if (rando < 5 ){
               cssClass = "male";
            }

            var div  = document.createElement("div");
            div.classList.add("lifter");
            div.classList.add(cssClass);
            div.id = e.instances[i].name;
            holder.appendChild(div);
        }
        div.innerHTML = "<span>" + e.instances[i].requests + "</span>";

    }

}

function getRandomArbitrary(min, max) {
    return Math.random() * (max - min) + min;
}

function handleError(e){
    console.log(e);
}

function log(e){
    console.log(e);
}


function generateToken(){
 return (Math.random().toString(36)+'00000000000000000').slice(2, 5+2)
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
    if (load_send > 0){
        $(".overlay-gaelist").fadeIn('slow');
        $(".overlay-alert").fadeIn('slow');
    }
}

function hide_explanation(){
    $(".overlay").fadeOut('slow');
}

