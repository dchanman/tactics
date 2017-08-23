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
    function clearStatusIconColor(statusIconId) {
        $(statusIconId).removeClass("status-icon-current-player");
        $(statusIconId).removeClass("status-icon-thinking");
        $(statusIconId).removeClass("status-icon-ready");
        $(statusIconId).removeClass("status-icon-offline");
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
        var color,
            ctrlStatusRole = "Your role: ".bold() + update.role;
        if (update.team > 0) {
            color = (update.team === 1 ? " White" : " Black");
            ctrlStatusRole += color.italics();
        }
        $("#ctrlStatusRole").html(ctrlStatusRole);
        this.renderButtons();
        clearStatusIconColor("#status-icon-role-p1");
        clearStatusIconColor("#status-icon-role-p2");
        if (update.team === 1) {
            $("#status-icon-role-p1").addClass("status-icon-current-player");
        } else if (update.team === 2) {
            $("#status-icon-role-p2").addClass("status-icon-current-player");
        }
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
            ctrlStatusP1 = "White: ".bold() + p1,
            ctrlStatusP2 = "Black: ".bold() + p2,
            classp1 = (update.p1available ? (update.p1ready ? "status-icon-ready" : "status-icon-thinking") : "status-icon-offline"),
            classp2 = (update.p2available ? (update.p2ready ? "status-icon-ready" : "status-icon-thinking") : "status-icon-offline");
        $("#ctrlStatusP1").html(ctrlStatusP1);
        $("#ctrlStatusP2").html(ctrlStatusP2);
        this.renderButtons();
        $("#status-score-p1").html(p1);
        $("#status-score-p2").html(p2);
        clearStatusIconColor("#status-icon-status-p1");
        clearStatusIconColor("#status-icon-status-p2");
        $("#status-icon-status-p1").addClass(classp1);
        $("#status-icon-status-p2").addClass(classp2);
    };
    return Status;
}());
