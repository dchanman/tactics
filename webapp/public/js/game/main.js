window.Main = (function () {
    "use strict";
    function Main() {
        this.api = new Api();
        this.chat = new Chat(this);
        this.board = new Renderer(this, document.getElementById("game"), document.getElementById("history"));
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
            case "Game.Status":
                main.handleStatus(params);
                break;
            case "Game.Turn":
                main.board.handleGameTurn(params);
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
        this.board.engineInit(data);
        this.status.updatePlayerReady(data);
    };
    Main.prototype.handleStatus = function (data) {
        var main = this;
        this.status.updatePlayerReady(data);
        // Why are we getting role every time??
        this.api.getRole()
            .then(function (result) {
                main.status.updateRole(result);
                main.board.setPlayerTeam(result.team);
            });
    };
    Main.prototype.handleGameOver = function (data) {
        if (data.team === 0) {
            this.chat.notification("Draw game!");
        } else {
            var winner = (data.team === 1 ? "White" : "Black") + " is victorious!";
            this.chat.notification(winner);
        }
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
            showCollisions = $("#overlay-setting-ctrl-showCollisions").is(':checked'),
            moveConfirmation = $("#overlay-setting-ctrl-moveConfirmation").is(':checked');
        main.board.setOverlaySettings(showLastMove, showCollisions, moveConfirmation);
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
