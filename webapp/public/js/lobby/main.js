window.Main = (function () {
    "use strict";
    function Main() {
        console.log("Starting");
    }
    Main.prototype.createGameList = function (gameIDs) {
        var i, li;
        for (i = 0; i < gameIDs.length; i += 1) {
            gameIDs[i] = Util.zeropad(gameIDs[i]);
        }
        gameIDs.sort();
        for (i = 0; i < gameIDs.length; i += 1) {
            li = document.createElement("li");
            $(li).html('<a href="/g/' + gameIDs[i] + '">' + gameIDs[i] + '</a>');
            $("#gameList").append(li);
        }
    };
    Main.prototype.run = function () {
        var main = this;
        console.log("Running");
        Jsonrpc.getGameIds()
            .then(function (result) {
                main.createGameList(result.result.gameids);
            })
            .catch(function (err) {
                console.log(err);
            });
    };
    return Main;
}());

$(document).ready(function () {
    "use strict";
    var main = new Main();
    $("#gameCreateSmall").click(function () {
        Jsonrpc.createGame("small")
            .then(function (result) {
                window.location.href = "/g/" + Util.zeropad(result.result.gameid);
            })
            .catch(function (err) {
                console.log(err);
            });
    });
    $("#gameCreateLarge").click(function () {
        Jsonrpc.createGame("large")
            .then(function (result) {
                window.location.href = "/g/" + Util.zeropad(result.result.gameid);
            })
            .catch(function (err) {
                console.log(err);
            });
    });
    main.run();
});
