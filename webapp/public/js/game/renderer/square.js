window.Square = (function () {
    "use strict";
    function Square(board, x, y) {
        this.board = board;
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
        if (this === this.board.selectedSquare) {
            this.board.removeSelectableSquares();
            this.board.selectedSquare = null;
        } else if ($(this.container).hasClass("grid-square-selectable")) {
            this.board.removeSelectableSquares();
            $(this.board.selectedSquare.container).addClass("grid-square-commit-src");
            $(this.container).addClass("grid-square-commit-dst");
        } else if ($(this.container).hasClass("grid-square-commit-dst")) {
            this.board.removeSelectableSquares();
            this.board.main.api.commitMove(
                this.board.selectedSquare.x,
                this.board.selectedSquare.y,
                this.x,
                this.y
            );
            this.board.selectedSquare = null;
        } else {
            this.board.removeSelectableSquares();
            this.board.selectedSquare = this;
            var self = this,
                x,
                y;
            // Set valid moves
            if (this.unit && this.board.playerTeam === this.unit.team) {
                for (x = 0; x < this.board.cols; x += 1) {
                    if (x !== this.x) {
                        self.board.setActiveSquare(x, this.y);
                    }
                }
                for (y = 0; y < this.board.rows; y += 1) {
                    if (y !== this.y) {
                        self.board.setActiveSquare(this.x, y);
                    }
                }
            }
        }
    };
    return Square;
}());