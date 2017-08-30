window.Square = (function () {
    "use strict";
    function Square(renderer, x, y) {
        this.renderer = renderer;
        this.x = x;
        this.y = y;
        this.unit = null;
        this.dom = null;
        this.container = null;
    }
    Square.prototype.setDOM = function (td, div) {
        this.container = td;
        this.dom = div;
        var self = this;
        $(td).click(function () {
            self.onClick();
        });
    };
    Square.prototype.onClick = function () {
        if (this === this.renderer.selectedSquare) {
            this.renderer.removeSelectableSquares();
            this.renderer.selectedSquare = null;
        } else if ($(this.container).hasClass("grid-square-selectable")) {
            this.renderer.removeSelectableSquares();
            $(this.renderer.selectedSquare.container).addClass("grid-square-commit-src");
            $(this.container).addClass("grid-square-commit-dst");
        } else if ($(this.container).hasClass("grid-square-commit-dst")) {
            this.renderer.removeSelectableSquares();
            this.renderer.main.api.commitMove(
                this.renderer.selectedSquare.x,
                this.renderer.selectedSquare.y,
                this.x,
                this.y
            );
            this.renderer.selectedSquare = null;
        } else {
            this.renderer.removeSelectableSquares();
            this.renderer.selectedSquare = this;
            var self = this,
                x,
                y;
            // Set valid moves
            if (this.unit && this.renderer.playerTeam === this.unit.team) {
                for (x = 0; x < this.renderer.cols; x += 1) {
                    if (x !== this.x) {
                        self.renderer.setActiveSquare(x, this.y);
                    }
                }
                for (y = 0; y < this.renderer.rows; y += 1) {
                    if (y !== this.y) {
                        self.renderer.setActiveSquare(this.x, y);
                    }
                }
            }
        }
    };
    return Square;
}());