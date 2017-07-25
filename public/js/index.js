window.Main = (function () {
    "use strict";
    function Main() {
        this.board = new Board(document.getElementById("game"));
        this.api = new Api();
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
var main = new Main();
