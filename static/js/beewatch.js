var wsUri = "ws://" + window.location.hostname + ":" + window.location.port + "/beewatch";
var output;
var connected = false;
var websocket = new WebSocket(wsUri);

function init() {
    output = document.getElementById("output");
    setupWebSocket();
}

function setupWebSocket() {
    websocket.onopen = function (evt) {
        onOpen(evt)
    };
    websocket.onclose = function (evt) {
        onClose(evt)
    };
    websocket.onmessage = function (evt) {
        onMessage(evt)
    };
    websocket.onerror = function (evt) {
        onError(evt)
    };
}

function onOpen(evt) {
    connected = true;
    //document.getElementById("disconnect").className = "buttonEnabled";
    writeToScreen("Connection has established.", "label label-funky", "INFO", "");
    sendConnected();
}

function onClose(evt) {
    //handleDisconnected();
}

function onMessage(evt) {
    try {
        var cmd = JSON.parse(evt.data);
    } catch (e) {
        console.log('[hopwatch] failed to read valid JSON: ', message.data);
        return;
    }
}

function onError(evt) {
    //writeToScreen(evt, "err mono");
}

function writeToScreen(title, cls, level, msg) {
    var logdiv = document.createElement("div");
    addTime(logdiv, cls, level);
    addTitle(logdiv, title);
    addMessage(logdiv, msg);
    logdiv.scrollIntoView();
    output.appendChild(logdiv);
}

function addTime(logdiv, cls, level) {
    var stamp = document.createElement("span");
    stamp.innerHTML = timeHHMMSS() + " " + level;
    stamp.className = cls;
    logdiv.appendChild(stamp);
}
function timeHHMMSS() {
    return new Date().toTimeString().replace(/.*(\d{2}:\d{2}:\d{2}).*/, "$1");
}

function addTitle(logdiv, title) {
    var name = document.createElement("span");
    name.innerHTML = " " + title;
    logdiv.appendChild(name);
}


function addMessage(logdiv, msg) {
    var txt = document.createElement("span");
    var msgcls;

    switch (msg.substr(1, 4)) {
        case "INIT", "INFO":
            txt.className = "text-success";
    }

    txt.innerHTML = " " + msg;
    logdiv.appendChild(txt);
}

function actionDisconnect() {
    if (!connected) return;
    connected = false;
    //document.getElementById("disconnect").className = "buttonDisabled";
    sendQuit();
    writeToScreen("Disconnected.", "label label-funky", "INFO", "");
    websocket.close();  // seems not to trigger close on Go-side ; so handleDisconnected cannot be used here.
}

function sendConnected() {
    doSend('{"Action":"CONNECTED"}');
}

function sendQuit() {
    doSend('{"Action":"QUIT"}');
}

function doSend(message) {
    // console.log("[hopwatch] send: " + message);
    websocket.send(message);
}

window.addEventListener("load", init, false);