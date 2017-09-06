window.Main = (function () {
    "use strict";
    function Main() {
        console.log("Starting");
    }
    function redirectToGameIDFunc(gameID) {
        return function () {
            window.location.href = "/g/" + Util.zeropad(gameID);
        };
    }
    Main.prototype.createGameList = function (gameIDs) {
        var i;
        for (i = 0; i < gameIDs.length; i += 1) {
            gameIDs[i] = Util.zeropad(gameIDs[i]);
        }
        gameIDs.sort();
        for (i = 0; i < gameIDs.length; i += 1) {
            $("<div></div>")
                .addClass("card card-inverse card-joingame")
                .append(
                    $("<div></div>")
                        .addClass("card-block")
                        .append(
                            $("<h5></h5>")
                                .addClass("card-title")
                                .html(gameIDs[i])
                        )
                )
                .appendTo("#gameList")
                .click(redirectToGameIDFunc(gameIDs[i]));
        }
    };
    Main.prototype.run = function () {
        var main = this;
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
