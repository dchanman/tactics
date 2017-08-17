window.Main = (function () {
    "use strict";
    function Main() {
        console.log("Starting");
    }
    function zeropad(id) {
        var s = id.toString();
        while (s.length < 6) {
            s = "0" + s;
        }
        return s;
    }
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
    Main.prototype.getGameIDs = function () {
        var main = this,
            jsonrpc = {
                "id": 1,
                "jsonrpc": "2.0",
                "method": "Server.GetGameIds",
                "params": []
            };
        ajax(jsonrpc)
            .then(function (result) {
                main.createGameList(result.result.gameids);
            })
            .catch(function (err) {
                console.log(err);
            });
    };
    Main.prototype.createGame = function (gameSize) {
        var jsonrpc = {
                "id": 1,
                "jsonrpc": "2.0",
                "method": "Server.CreateGame",
                "params": [{"gameType": gameSize}]
            };
        ajax(jsonrpc)
            .then(function (result) {
                window.location.href = "/g/" + zeropad(result.result.gameid);
            })
            .catch(function (err) {
                console.log(err);
            });
    };
    Main.prototype.createGameList = function (gameIDs) {
        var i, li;
        for (i = 0; i < gameIDs.length; i += 1) {
            gameIDs[i] = zeropad(gameIDs[i]);
        }
        gameIDs.sort();
        for (i = 0; i < gameIDs.length; i += 1) {
            li = document.createElement("li");
            $(li).html('<a href="/g/' + gameIDs[i] + '">' + gameIDs[i] + '</a>');
            $("#gameList").append(li);
        }
    };
    Main.prototype.run = function () {
        console.log("Running");
        this.getGameIDs();
    };
    return Main;
}());

$(document).ready(function () {
    "use strict";
    var main = new Main();
    $("#gameCreateSmall").click(function () {
        main.createGame("small");
    });
    $("#gameCreateLarge").click(function () {
        main.createGame("large");
    });
    main.run();
});
