window.Main = (function () {
    "use strict";
    console.log("Hello World!");
    var Main = function () {
        var uri = (window.location.protocol === "https:") ? "wss://" : "ws://",
            board;
        this.board = new Board(document.getElementById("game"));
        board = this.board;
        this.ws = new WebSocket(uri + window.location.host + "/ws");
        this.ws.onopen = function () {
            this.send('{"id": 1, "method": "TacticsApi.Hello", "params": []}');
            this.send('{"id": 2, "method": "TacticsApi.GetGame", "params": []}');
        };
        this.ws.onmessage = function (event) {
            console.log("\nReceived new message!!!");
            console.log(event);
            console.log(event.data);
            var data = JSON.parse(event.data);
            if (data) {
                if (data.error !== null) {
                    console.log(data.error);
                    return;
                }
                console.log(data);
                board.render(data.result.game);
            }
        };
        this.ws.onclose = function () {
            console.log("Websocket closed");
        };
        this.ws.onerror = function (err) {
            console.log("Error: " + err);
        };
    };
    Main.prototype.refresh = function () {
        // TODO: fill me in
    };
    Main.prototype.start = function () {
        this.refresh();
    };
    return Main;
}());
var main = new Main();
main.start();
