window.Status = (function () {
    "use strict";
    function Status(main) {
        this.team = 0;
        this.p1available = false;
        this.p2available = false;
        this.p1ready = false;
        this.p2ready = false;
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
        $("#ctrlReset").hide();
        $("#ctrlSitP1").hide();
        $("#ctrlSitP2").hide();
    }
    Status.prototype.renderButtons = function () {
        if (this.team === 0) {
            $("ctrlReset").hide();
            if (!this.p1available) {
                $("#ctrlSitP1").show();
            } else {
                $("#ctrlSitP1").hide();
            }
            if (!this.p2available) {
                $("#ctrlSitP2").show();
            } else {
                $("#ctrlSitP2").hide();
            }
        } else {
            $("#ctrlReset").show();
            $("#ctrlSitP1").hide();
            $("#ctrlSitP2").hide();
        }
    };
    Status.prototype.updateRole = function (update) {
        this.team = update.team;
        var ctrlStatusRole = "Your role: ".bold() + update.role;
        if (update.team > 0) {
            ctrlStatusRole += ("(Team " + update.team + ")").italics();
        }
        $("#ctrlStatusRole").html(ctrlStatusRole);
        this.renderButtons();
    };
    Status.prototype.updatePlayerReady = function (update) {
        console.log("updatePlayerReady");
        console.log(update);
        this.p1available = update.p1available;
        this.p2available = update.p2available;
        this.p1ready = update.p1ready;
        this.p2ready = update.p2ready;
        var p1 = (update.p1available ? (update.p1ready ? "Ready" : "Thinking...") : "Offline"),
            p2 = (update.p2available ? (update.p2ready ? "Ready" : "Thinking...") : "Offline"),
            ctrlStatusP1 = "Player 1: ".bold() + p1,
            ctrlStatusP2 = "Player 2: ".bold() + p2;
        $("#ctrlStatusP1").html(ctrlStatusP1);
        $("#ctrlStatusP2").html(ctrlStatusP2);
        this.renderButtons();
    };
    return Status;
}());
