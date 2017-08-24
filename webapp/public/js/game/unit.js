window.Unit = (function () {
    "use strict";
    function Unit() {
        this.stack = 0;
        this.team = 0;
        this.exists = true;
    }
    Unit.fromJSON = function (json) {
        var u = new Unit();
        u.stack = json.stack;
        u.team = json.team;
        u.exists = json.exists;
        return u;
    };
    Unit.prototype.getRenderHtml = function (width) {
        var num = (this.stack > 1 ? this.stack : ""),
            cls = "piece";
        cls += " " + (this.team === 1 ? "piece-1" : "piece-2");
        return '<svg class="' + cls + '"><circle cx="50%" cy="50%" r="' + (width * 0.4) + '"></circle><text x="50%" y="50%">' + num + '</text></svg>';
    };
    return Unit;
}());
