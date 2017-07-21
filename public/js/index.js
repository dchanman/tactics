window.Main = (function () {
    "use strict";
    console.log("Hello World!");
    function Main() {
        this.board = new Board(document.getElementById("game"));
        this.ws = this.createNewWebsocket();
        this.sentmsgs = {};
    }
    Main.prototype.createNewWebsocket = function () {
        var main = this,
            uri = (window.location.protocol === "https:") ? "wss://" : "ws://",
            ws = new WebSocket(uri + window.location.host + "/ws"),
            currid = 0,
            sentmsgs = {};
        ws.onopen = function () {
            main.start();
        };
        ws.onmessage = function (event) {
            var data = JSON.parse(event.data),
                callback;
            if (data && data.id) {
                callback = sentmsgs[data.id];
                if (callback) {
                    callback(data.result, data.error);
                }
            }
        };
        ws.onclose = function () {
            console.log("Websocket closed");
        };
        ws.onerror = function (err) {
            console.log("Error: " + err);
        };
        ws.sendmsg = function (method, params, callback) {
            var call = {
                "id": currid,
                "method": method,
                "params": params
            };
            sentmsgs[currid] = callback;
            currid += 1;
            this.send(JSON.stringify(call));
        };
        return ws;
    };
    Main.prototype.refresh = function () {
        var main = this;
        this.ws.sendmsg("TacticsApi.GetGame", [], function (result, err) {
            if (!err) {
                main.board.render(result.game);
            }
        });
    };
    Main.prototype.start = function () {
        this.ws.sendmsg("TacticsApi.Hello", [], function () {
            console.log("Received resonse to Hello!");
        });
        this.refresh();
    };
    return Main;
}());
var main = new Main();
