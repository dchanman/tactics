$(document).ready(function () {
    "use strict";
    function createGame(gameSize) {
        Jsonrpc.createGame(gameSize)
            .then(function (result) {
                window.location.href = "/g/" + Util.zeropad(result.result.gameid);
            })
            .catch(function (err) {
                console.log(err);
            });
    }
    $("#navCreateGameSmall").click(function () {
        createGame("small");
    });
    $("#navCreateGameLarge").click(function () {
        createGame("large");
    });
});