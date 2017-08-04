window.Main = (function () {
    "use strict";
    function Main() {
        this.api = new Api();
        this.chat = new Chat(this);
        this.board = new Board(document.getElementById("game"), this);
        this.status = Status;
        this.heartbeat = null;
        var main = this;
        this.api.onready = function () {
            main.start();
        };
        this.api.onerror = function (err) {
            console.log(err);
        };
        this.api.onclose = function () {
            console.log("API closed");
            clearInterval(main.heartbeat);
        };
        this.api.onupdate = function (method, params) {
            switch (method) {
            case "Game.Update":
                console.log("Received update!");
                main.board.render(params.game.board);
                main.status.updatePlayerReady(params.p1ready, params.p2ready);
                break;
            case "Game.Chat":
                main.chat.receiveMessage(params);
                break;
            case "Game.Status":
                main.status.update(params);
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
        this.api.getStatus()
            .then(function (result) {
                console.log(result);
                main.status.update(result);
            });
    };
    Main.prototype.start = function () {
        var main = this;
        main.refresh();
        main.heartbeat = setInterval(function () {
            main.api.heartbeat();
        }, 10000);
    };
    return Main;
}());

$(document).ready(function () {
    "use strict";
    var main = new Main();
    $("#chatsend").click(function () {
        main.chat.sendMessage();
    });
    $("#chatmsg").keydown(function (event) {
        if (event.keyCode === 13) {
            main.chat.sendMessage();
        }
    });
    $("#ctrlReset").click(function () {
        main.api.resetBoard();
    });
    $("#ctrlSitP1").click(function () {
        main.api.joinGame(1)
            .catch(function (err) {
                console.log(err);
            });
    });
    $("#ctrlSitP2").click(function () {
        main.api.joinGame(2);
    });
});
