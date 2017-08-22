window.Api = (function () {
    "use strict";
    function createNewWebsocket(api) {
        var uri = (window.location.protocol === "https:") ? "wss://" : "ws://",
            url,
            ws;
        url = uri + window.location.host + "/ws";
        ws = new WebSocket(url);
        console.log(url);
        ws.onopen = function () {
            var gameid = /g\/([0-9]{6})/.exec(window.location),
                gameidint;
            if (gameid.length < 2) {
                throw "invalid URL: gameid is bad";
            }
            gameidint = parseInt(gameid[1], 10);
            console.log(gameidint);
            api.subscribeGame(gameidint)
                .then(function () {
                    api.onready();
                })
                .catch(function (err) {
                    throw err;
                });
            api.subscribeChat(gameidint)
                .catch(function (err) {
                    throw err;
                });
        };
        ws.onmessage = function (event) {
            var data, callback;
            try {
                data = JSON.parse(event.data);
            } catch (err) {
                data = null;
                console.log(err);
                console.log(event.data);
            }
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
        try {
            api.ws.send(JSON.stringify(call));
        } catch (err) {
            callback(null, err);
        }
    }
    Api.prototype.heartbeat = function () {
        var api = this;
        return new Promise(function (resolve, reject) {
            sendmsg(api, "TacticsApi.Heartbeat", [], function (result, err) {
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
    Api.prototype.getRole = function () {
        var api = this;
        return new Promise(function (resolve, reject) {
            sendmsg(api, "TacticsApi.GetRole", [], function (result, err) {
                if (err) {
                    reject(err);
                } else {
                    resolve(result);
                }
            });
        });
    };
    Api.prototype.joinGame = function (playerNumber) {
        var api = this;
        return new Promise(function (resolve, reject) {
            var params = {
                "playerNumber": playerNumber
            };
            sendmsg(api, "TacticsApi.JoinGame", [params], function (result, err) {
                if (err) {
                    reject(err);
                } else {
                    resolve(result);
                }
            });
        });
    };
    Api.prototype.sendChat = function (msg) {
        var api = this;
        return new Promise(function (resolve, reject) {
            var params = {
                "message": msg
            };
            sendmsg(api, "TacticsApi.SendChat", [params], function (result, err) {
                if (err) {
                    reject(err);
                } else {
                    resolve(result);
                }
            });
        });
    };
    Api.prototype.commitMove = function (fromX, fromY, toX, toY) {
        var api = this;
        return new Promise(function (resolve, reject) {
            var params = {
                "fromX": fromX,
                "fromY": fromY,
                "toX": toX,
                "toY": toY
            };
            sendmsg(api, "TacticsApi.CommitMove", [params], function (result, err) {
                if (err) {
                    reject(err);
                } else {
                    resolve(result);
                }
            });
        });
    };
    Api.prototype.resetBoard = function () {
        var api = this;
        return new Promise(function (resolve, reject) {
            sendmsg(api, "TacticsApi.ResetBoard", [], function (result, err) {
                if (err) {
                    reject(err);
                } else {
                    resolve(result);
                }
            });
        });
    };
    Api.prototype.subscribeGame = function (id) {
        var api = this;
        return new Promise(function (resolve, reject) {
            var params = {
                "id": id
            };
            sendmsg(api, "TacticsApi.SubscribeGame", [params], function (result, err) {
                if (err) {
                    reject(err);
                } else {
                    resolve(result);
                }
            });
        });
    };
    Api.prototype.subscribeChat = function (id) {
        var api = this;
        return new Promise(function (resolve, reject) {
            var params = {
                "id": id
            };
            sendmsg(api, "TacticsApi.SubscribeChat", [params], function (result, err) {
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