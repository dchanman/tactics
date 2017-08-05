window.Status = (function () {
    "use strict";
    var Status = {};
    Status.update = function (update) {
        var ctrlStatusRole = "Your role: ".bold() + update.role;
        $("#ctrlStatusRole").html(ctrlStatusRole);
    };
    Status.updatePlayerReady = function (update) {
        console.log("updatePlayerReady");
        console.log(update);
        var p1 = (update.p1available ? (update.p1ready ? "Ready" : "Thinking...") : "Offline"),
            p2 = (update.p2available ? (update.p2ready ? "Ready" : "Thinking...") : "Offline"),
            ctrlStatusP1 = "Player 1: ".bold() + p1,
            ctrlStatusP2 = "Player 2: ".bold() + p2;
        $("#ctrlStatusP1").html(ctrlStatusP1);
        $("#ctrlStatusP2").html(ctrlStatusP2);
    };
    return Status;
}());
