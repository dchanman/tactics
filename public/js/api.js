window.Api = (function () {
    "use strict";
    function createNewWebsocket(api) {
        var uri = (window.location.protocol === "https:") ? "wss://" : "ws://",
            ws = new WebSocket(uri + window.location.host + "/ws");
        ws.onopen = function () {
            api.onready();
        };
        ws.onmessage = function (event) {
            var data = JSON.parse(event.data),
                callback;
            if (data !== null) {
                if (data.id !== undefined) {
                    // Handle RPC response
                    callback = api.sentmsgs[data.id];
                    api.sentmsgs[data.id] = null;
                    if (callback) {
                        callback(data.result, data.error);
                    }
                } else {
                    // Handle notification
                    api.onupdate(data.method, data.params);
                }
            }
        };
        ws.onclose = function () {
            api.onclose();
        };
        ws.onerror = function (err) {
            api.onerror(err);
        };
        return ws;
    }
    function Api() {
        this.ws = createNewWebsocket(this);
        this.sentmsgs = {};
        this.currid = 1;
        // Handler functions
        this.onready = function () { return; };
        this.onerror = function () { return; };
        this.onclose = function () { return; };
        this.onupdate = function () { return; };
    }
    function sendmsg(api, method, params, callback) {
        var call = {
            "id": api.currid,
            "method": method,
            "params": params
        };
        api.sentmsgs[api.currid] = callback;
        api.currid += 1;
        api.ws.send(JSON.stringify(call));
    }
    Api.prototype.hello = function () {
        var api = this;
        return new Promise(function (resolve, reject) {
            sendmsg(api, "TacticsApi.Hello", [], function (result, err) {
                if (err) {
                    reject(err);
                } else {
                    resolve(result);
                }
            });
        });
    };
    Api.prototype.getGame = function () {
        var api = this;
        return new Promise(function (resolve, reject) {
            sendmsg(api, "TacticsApi.GetGame", [], function (result, err) {
                if (err) {
                    reject(err);
                } else {
                    resolve(result);
                }
            });
        });
    };
    Api.prototype.addUnit = function (x, y, unit) {
        var api = this;
        return new Promise(function (resolve, reject) {
            var params = {
                "x": x,
                "y": y,
                "unit": unit
            };
            sendmsg(api, "TacticsApi.AddUnit", [params], function (result, err) {
                if (err) {
                    reject(err);
                } else {
                    resolve(result);
                }
            });
        });
    };
    return Api;
}());