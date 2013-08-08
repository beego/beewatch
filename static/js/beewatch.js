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

function handleDisconnected() {
    //connected = false;
    //document.getElementById("resume").className = "buttonDisabled";
    //document.getElementById("disconnect").className = "buttonDisabled";
    //writeToScreen("Disconnected.", "label label-funky", "INFO", "");
}

function onMessage(evt) {
    try {
        var cmd = JSON.parse(evt.data);
    } catch (e) {
        console.log('[ERRO] Failed to read valid JSON: ', e.message.data);
        return;
    }

    switch (cmd.Action) {
        case "DISPLAY":
            writeToScreen(getTitle(cmd), "label label-info", "INFO", watchParametersToHtml(cmd.Parameters));
            sendResume();
            return;
        case "DONE":
            actionDisconnect();
            return;
    }
}

function getTitle(cmd) {
    var i = cmd.Parameters["go.file"].lastIndexOf("/") + 1;
    return cmd.Parameters["go.file"].substring(i, cmd.Parameters["go.file"].length) + ":" + cmd.Parameters["go.line"];
}

function onError(evt) {
    //writeToScreen(evt, "err mono");
}

function writeToScreen(title, cls, level, msg) {
    var logdiv = document.createElement("div");
    addTime(logdiv, cls, level);
    addTitle(logdiv, title);
    addMessage(logdiv, msg, level);
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


function addMessage(logdiv, msg, level) {
    var txt = document.createElement("span");

    if (msg.substr(0, 1) == "[") {
        // Debugger messages.
        txt.className = getMsgClass(msg.substr(1, 4));
    } else {
        // App messages.
        txt.className = getMsgClass(level);
    }

    txt.innerHTML = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;" + msg;
    logdiv.appendChild(txt);
}

function getMsgClass(level) {
    switch (level) {
        case "INIT", "INFO":
            return "text-success";
    }
}

function watchParametersToHtml(parameters) {
    var line = "";
    var multiline = false;
    for (var prop in parameters) {
        if (prop.slice(0, 3) != "go.") {
            if (multiline) {
                line = line + ", ";
            }
            line = line + prop + "=" + parameters[prop];
            multiline = true;
        }
    }
    return line
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

function sendResume() {
    doSend('{"Action":"RESUME"}');
}

function sendQuit() {
    doSend('{"Action":"QUIT"}');
}

function doSend(message) {
    // console.log("[hopwatch] send: " + message);
    websocket.send(message);
}

window.addEventListener("load", init, false);