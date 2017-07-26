window.Main = (function () {
    "use strict";
    function Main() {
        this.api = new Api();
        this.board = new Board(document.getElementById("game"), this);
        var main = this;
        this.api.onready = function () {
            main.start();
        };
        this.api.onerror = function (err) {
            console.log(err);
        };
        this.api.onclose = function () {
            console.log("API closed");
        };
        this.api.onupdate = function (method, params) {
            switch (method) {
            case "Game.Update":
                main.board.render(params.game.board);
                break;
            case "Game.Chat":
                console.log("Received: " + params.message);
                $("#chat").append(document.createTextNode(params.message + "\n"));
                break;
            default:
                console.log("Error: Unknown method: " + method);
            }
        };
    }
    Main.prototype.refresh = function () {
        var main = this;
        this.api.getGame()
            .then(function (result) {
                main.board.render(result.game.board);
            });
    };
    Main.prototype.start = function () {
        var main = this;
        this.api.hello()
            .then(function () {
                console.log("Connected to server!");
                main.refresh();
            });
    };
    return Main;
}());

$(document).ready(function () {
    "use strict";
    var main = new Main();
    $("#chatsend").click(function () {
        console.log($("#chatmsg").val());
        main.api.sendChat($("#chatmsg").val());
    });
});
