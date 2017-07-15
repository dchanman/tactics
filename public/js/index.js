window.Main = (function () {
    "use strict";
    console.log("Hello World!");
    var Main = function () {
        var uri = (window.location.protocol === "https:") ? "wss://" : "ws://";
        this.board = new Board(document.getElementById("game"));
        this.ws = new WebSocket(uri + window.location.host + "/ws");
        this.ws.onopen = function () {
            this.send('{"id": 1, "method": "TacticsApi.Hello", "params": []}');
        };
        this.ws.onmessage = function (event) {
            console.log("\nReceived new message!!!");
            console.log(event);
        };
        this.ws.onclose = function () {
            console.log("Websocket closed");
        };
        this.ws.onerror = function (err) {
            console.log("Error: " + err);
        };
    };
    Main.prototype.refresh = function () {
        var xmlHttp = new XMLHttpRequest(),
            board = this.board;
        xmlHttp.onreadystatechange = function () {
            if (xmlHttp.readyState === 4 && xmlHttp.status === 200) {
                board.render(xmlHttp.responseText);
            }
        };
        xmlHttp.open("GET", "/game", true);
        xmlHttp.send(null);
    };
    Main.prototype.start = function () {
        this.refresh();
    };
    return Main;
}());
var main = new Main();
main.start();
