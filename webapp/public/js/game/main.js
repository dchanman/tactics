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
                main.board.runEngine(params);
                break;
            case "Game.Chat":
                main.chat.receiveMessage(params);
                break;
            case "Game.Over":
                main.handleGameOver(params);
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
    Main.prototype.handleGameOver = function (data) {
        console.log("Game over received");
        console.log(data);
        if (data.team === 0) {
            this.chat.notification("Draw game!");
        } else {
            var winner = (data.team === 1 ? "White" : "Black") + " is victorious!";
            this.chat.notification(winner);
        }
        this.board.handleGameOver(data.team);
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
    var main = new Main(),
        resizeTimer = null;
    $("#chatsend").click(function () {
        main.chat.sendMessage();
    });
    $("#chatmsg").keydown(function (event) {
        if (event.keyCode === 13) {
            main.chat.sendMessage();
        }
    });
    function setOverlaySettings() {
        var showLastMove = $("#overlay-setting-ctrl-showLastMove").is(':checked'),
            showLastUnit = $("#overlay-setting-ctrl-showLastUnit").is(':checked'),
            showCollisions = $("#overlay-setting-ctrl-showCollisions").is(':checked');
        main.board.setOverlaySettings(showLastMove, showLastUnit, showCollisions);
    }
    setOverlaySettings();
    $(".overlay-setting-ctrl").click(function () {
        setOverlaySettings();
    });
    $(window).resize(function () {
        if (resizeTimer !== null) {
            window.clearTimeout(resizeTimer);
        }
        resizeTimer = window.setTimeout(function () {
            main.board.onWindowResize();
        }, 200);
    });
});
