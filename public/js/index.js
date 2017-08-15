window.Main = (function () {
    "use strict";
    function Main() {
        this.api = new Api();
        this.chat = new Chat(this);
        this.board = new Board(document.getElementById("game"), this);
        this.status = new Status(this);
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
                main.handleGameInfo(params);
                break;
            case "Game.Chat":
                main.chat.receiveMessage(params);
                break;
            default:
                console.log("Error: Unknown method: " + method);
            }
        };
    }
    Main.prototype.handleGameInfo = function (data) {
        var main = this;
        console.log("Received update!");
        console.log(data.history);
        this.board.render(data.board);
        this.board.renderHistory(data.history);
        this.status.updatePlayerReady(data);
        this.api.getRole()
            .then(function (result) {
                main.status.updateRole(result);
                main.board.setPlayerTeam(result.team);
            });
    };
    Main.prototype.refresh = function () {
        var main = this;
        this.api.getGame()
            .then(function (result) {
                main.handleGameInfo(result);
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
});
