window.Overlay = (function () {
    "use strict";
    function Overlay(htmlTable, rows, cols) {
        this.rows = rows;
        this.cols = cols;
        this.htmlTable = htmlTable;
        this.container = document.createElement("div");
        $(this.container).addClass("overlay-container");
        this.resize();
        htmlTable.appendChild(this.container);
    }
    Overlay.prototype.clear = function () {
        $(this.container).html("");
        $(this.container).removeClass("overlay-container-win");
        $(this.container).removeClass("overlay-container-lose");
        $(this.container).removeClass("overlay-container-draw");
    };
    Overlay.prototype.resize = function () {
        this.width = $(this.htmlTable).width();
        this.height = $(this.htmlTable).height();
        $(this.container).width(this.width);
        $(this.container).height(this.height);
    };
    Overlay.prototype.displayWinScreen = function () {
        $(this.container).addClass("overlay-container-win");
    };
    Overlay.prototype.displayLoseScreen = function () {
        $(this.container).addClass("overlay-container-lose");
    };
    Overlay.prototype.displayDrawScreen = function () {
        $(this.container).addClass("overlay-container-draw");
    };
    Overlay.prototype.renderMove = function (fromX, fromY, toX, toY, team, playerTeam) {
        var x1, y1, x2, y2, markerid, markerdef, line, cls;
        // Calculation:
        // (top left pixel coordinate of the desired square) + (offset to reach the square midpoint)
        x1 = parseFloat(fromX) / this.cols * this.width + (this.width / 2.0 / this.cols);
        x2 = parseFloat(toX) / this.cols * this.width + (this.width / 2.0 / this.cols);
        y1 = parseFloat(fromY) / this.rows * this.height + (this.height / 2.0 / this.rows);
        y2 = parseFloat(toY) / this.rows * this.height + (this.height / 2.0 / this.rows);
        markerid = "arrowhead" + team;
        markerdef = '\
<defs>\
    <marker id="' + markerid + '" markerWidth="10" markerHeight="10" refX="3" refY="2" orient="auto">\
      <path d="M0,0 L0,4 L5,2 z" stroke-width="0" />\
    </marker>\
</defs>';
        line = '<line x1="' + x1 + '" y1="' + y1 + '" x2="' + x2 + '" y2="' + y2 + '" marker-end="url(#' + markerid + ')"/>';
        if (playerTeam !== 0) {
            cls = (team === playerTeam ? "overlay-arrow-friendly" : "overlay-arrow-enemy");
        } else {
            cls = (team === 1 ? "overlay-arrow-friendly" : "overlay-arrow-enemy");
        }
        $(this.container).append('<svg class="overlay-arrow ' + cls + '" height="' + this.height + '" width="' + this.width + '">' + markerdef + line + '</svg>');
    };
    Overlay.prototype.renderPiece = function (x, y, unit) {
        var drawX, drawY, width, node;
        drawX = Math.round(parseFloat(x) / this.cols * this.width) + 1;
        drawY = Math.round(parseFloat(y) / this.rows * this.height) + 1;
        width = Math.round(this.width / this.cols * 0.75);
        node = $.parseHTML(unit.getRenderHtml(width));
        $(node).addClass("overlay-piece");
        $(node).css("position", "absolute")
            .css("left", drawX)
            .css("top", drawY)
            .css("opacity", 0.5);
        $(this.container).append(node);
    };
    Overlay.prototype.renderCollision = function (x, y) {
        var drawX, drawY, node;
        node = $.parseHTML('<div class="overlay-collision"><i class="fa fa-times" aria-hidden="true"></i></div>');
        drawX = Math.round(parseFloat(x) / this.cols * this.width) + (this.width / 2.0 / this.cols);
        drawY = Math.round(parseFloat(y) / this.rows * this.height) + (this.height / 2.0 / this.rows);
        $(node).css("position", "absolute")
            .css("left", drawX)
            .css("top", drawY);
        $(this.container).append(node);
    };
    return Overlay;
}());
