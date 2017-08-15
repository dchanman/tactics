window.Board = (function () {
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
            console.log("Committing move (" + this.board.selectedSquare.x + "," + this.board.selectedSquare.y + ") to (" + this.x + "," + this.y + ")");
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
            var self = this;
            this.board.main.api.getValidMoves(this.x, this.y)
                .then(function (result) {
                    var squares = result.validMoves,
                        i;
                    if (!squares) {
                        return;
                    }
                    for (i = 0; i < squares.length; i += 1) {
                        self.board.setActiveSquare(squares[i].x, squares[i].y);
                    }
                })
                .catch(function (err) {
                    console.log(err);
                });
        }
    };
    function Overlay(htmlTable, rows, cols) {
        this.rows = rows;
        this.cols = cols;
        this.container = document.createElement("div");
        $(this.container).addClass("overlay-container");
        this.width = $(htmlTable).width();
        this.height = $(htmlTable).height();
        $(this.container).width(this.width);
        $(this.container).height(this.height);
        htmlTable.appendChild(this.container);
    }
    Overlay.prototype.clear = function () {
        $(this.container).html("");
    };
    Overlay.prototype.renderMove = function (fromX, fromY, toX, toY, team) {
        var x1, y1, x2, y2, markerid, markerdef, line, cls;
        console.log("Drawing arrow from (" + fromX + "," + fromY + ") to (" + toX + "," + toY + ")");
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
        cls = (team === 1 ? "overlay-arrow-1" : "overlay-arrow-2");
        $(this.container).append('<svg class="overlay-arrow ' + cls + '" height="' + this.height + '" width="' + this.width + '">' + markerdef + line + '</svg>');
    };
    function Board(htmlTable, main) {
        this.cols = 0;
        this.rows = 0;
        this.playerTeam = 0;
        this.currentBoard = [];
        this.htmlTable = htmlTable;
        this.grid = [];
        this.overlay = null;
        this.selectedSquare = null;
        this.main = main;
    }
    Board.prototype.setPlayerTeam = function (team) {
        if (this.playerTeam !== team) {
            this.playerTeam = team;
            this.renderPieces(this.currentBoard);
        }
    };
    Board.prototype.removeSelectableSquares = function () {
        var i, j;
        for (i = 0; i < this.grid.length; i += 1) {
            for (j = 0; j < this.grid[i].length; j += 1) {
                $(this.grid[i][j].container).removeClass("grid-square-selectable");
                $(this.grid[i][j].container).removeClass("grid-square-commit-src");
                $(this.grid[i][j].container).removeClass("grid-square-commit-dst");
            }
        }
    };
    Board.prototype.render = function (board) {
        if (board.board.length !== (board.cols * board.rows)) {
            throw ("Invalid board data");
        }
        if (this.cols !== board.cols || this.rows !== board.rows) {
            this.createGrid(board.cols, board.rows);
        }
        this.renderPieces(board.board);
    };
    Board.prototype.renderHistory = function (history) {
        this.overlay.clear();
        if (history.length > 0) {
            var lastMove = history[history.length - 1],
                m1 = lastMove[1],
                m2 = lastMove[2];
            this.overlay.renderMove(m1.Src.x, m1.Src.y, m1.Dst.x, m1.Dst.y, 1);
            this.overlay.renderMove(m2.Src.x, m2.Src.y, m2.Dst.x, m2.Dst.y, 2);
        }
    };
    Board.prototype.setActiveSquare = function (x, y) {
        if (x < 0 || y < 0 || x > this.cols || y > this.rows) {
            throw ("Invalid square: (" + x + "," + y + ")");
        }
        $(this.grid[x][y].container).addClass("grid-square-selectable");
    };
    Board.prototype.createGrid = function (cols, rows) {
        var x, y, tr, td, div;
        this.cols = cols;
        this.rows = rows;
        // Create matrix
        this.grid = [];
        for (x = 0; x < cols; x += 1) {
            this.grid.push([]);
            for (y = 0; y < rows; y += 1) {
                this.grid[x].push(new Square(this, x, y));
            }
        }
        // Clean up existing DOM
        while (this.htmlTable.firstChild) {
            this.htmlTable.moveChild(this.htmlTable.firstChild);
        }
        // Create new DOM grid
        for (y = 0; y < rows; y += 1) {
            tr = document.createElement("tr");
            $(tr).addClass("grid");
            for (x = 0; x < cols; x += 1) {
                td = document.createElement("td");
                $(td).addClass("grid-square");
                div = document.createElement("div");
                $(div).addClass("grid-square-container");
                td.appendChild(div);
                tr.appendChild(td);
                this.grid[x][y].setDOM(td, div);
            }
            this.htmlTable.appendChild(tr);
        }
        // Create overlay
        this.overlay = new Overlay(this.htmlTable, rows, cols);
        console.log("Created grid");
    };
    Board.prototype.renderPieces = function (pieces) {
        var i, x, y, cls, num;
        this.currentBoard = pieces;
        for (i = 0; i < pieces.length; i += 1) {
            x = Math.floor(i / this.rows);
            y = i % this.rows;
            this.grid[x][y].unit = null;
            $(this.grid[x][y].dom).html("");
            $(this.grid[x][y].container).removeClass("piece-friendly");
            $(this.grid[x][y].container).removeClass("piece-enemy");
            if (pieces[i].exists) {
                this.grid[x][y].unit = pieces[i];
                cls = "piece";
                cls += " " + (pieces[i].team === 1 ? "piece-1" : "piece-2");
                num = (pieces[i].stack > 1 ? pieces[i].stack : "");
                $(this.grid[x][y].dom).html('<svg class="' + cls + '"><circle></circle><text x="50%" y="50%">' + num + '</text></svg>');
            }
        }
    };
    return Board;
}());