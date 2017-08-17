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
    Main.prototype.getGameIDs = function () {
        var main = this;
        $.ajax({
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            type: "POST",
            async: true,
            'url': '/data/',
            'data': '{ "jsonrpc": "2.0", "method": "Server.GetGameIds", "id": 1 , "params": []}',
            'dataType': 'json',
            'success': function (result) {
                main.createGameList(result.result.gameids);
            },
            'error': function (data) {
                console.log(data);
            }
        });
    };
    Main.prototype.createGameList = function (gameIDs) {
        var i, li;
        gameIDs.sort();
        for (i = 0; i < gameIDs.length; i += 1) {
            li = document.createElement("li");
            $(li).html('<a href="/g/' + zeropad(gameIDs[i]) + '">' + zeropad(gameIDs[i]) + '</a>');
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
    main.run();
});