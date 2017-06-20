window.Board = (function () {
    "use strict";
    var Board = function (rootDiv) {
        this.cols = 0;
        this.rows = 0;
        this.rootDiv = rootDiv;
    };
    Board.prototype.render = function (json) {
        console.log(json);
    };
    return Board;
}());