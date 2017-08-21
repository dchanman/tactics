window.Jsonrpc = (function () {
    "use strict";
    function ajax(data) {
        var d = JSON.stringify(data);
        return new Promise(function (resolve, reject) {
            $.ajax({
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                type: "POST",
                async: true,
                'url': '/data/',
                'data': d,
                'dataType': 'json',
                'success': function (result) {
                    resolve(result);
                },
                'error': function (err) {
                    reject(err);
                }
            });
        });
    }
    function getGameIds() {
        var jsonrpc = {
                "id": 1,
                "jsonrpc": "2.0",
                "method": "Server.GetGameIds",
                "params": []
            };
        return ajax(jsonrpc);
    }
    function createGame(gameSize) {
        var jsonrpc = {
                "id": 1,
                "jsonrpc": "2.0",
                "method": "Server.CreateGame",
                "params": [{"gameType": gameSize}]
            };
        return ajax(jsonrpc);
    }
    return {
        getGameIds: getGameIds,
        createGame: createGame
    };
}());