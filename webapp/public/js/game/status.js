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
        this.renderButtons();
        $("#status-score-p1").removeClass("status-score-current-player");
        $("#status-score-p2").removeClass("status-score-current-enemy");
        if (update.team === 1) {
            $("#status-score-p1").addClass("status-score-current-player");
        } else if (update.team === 2) {
            $("#status-score-p2").addClass("status-score-current-player");
        }
    };
    Status.prototype.updatePlayerReady = function (update) {
        this.p1available = update.p1available;
        this.p2available = update.p2available;
        this.p1ready = update.p1ready;
        this.p2ready = update.p2ready;
        var p1 = (update.p1available ? (update.p1ready ? "Ready" : "Thinking...") : "Offline"),
            p2 = (update.p2available ? (update.p2ready ? "Ready" : "Thinking...") : "Offline");
        this.renderButtons();
        $("#status-score-p1").html(p1);
        $("#status-score-p2").html(p2);
    };
    return Status;
}());
