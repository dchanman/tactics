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
    function Board(htmlTable, main) {
        this.cols = 0;
        this.rows = 0;
        this.playerTeam = 0;
        this.currentBoard = [];
        this.htmlTable = htmlTable;
        this.grid = [];
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
        console.log("Created grid");
    };
    Board.prototype.renderPieces = function (pieces) {
        var i, x, y, name;
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
                name = (pieces[i].team === 1 ? "X" : "O") + pieces[i].stack;
                if (this.playerTeam !== 0) {
                    if (pieces[i].team === this.playerTeam) {
                        $(this.grid[x][y].container).addClass("piece-friendly");
                    } else {
                        $(this.grid[x][y].container).addClass("piece-enemy");
                    }
                }
                $(this.grid[x][y].dom).html(name);
            }
        }
    };
    return Board;
}());