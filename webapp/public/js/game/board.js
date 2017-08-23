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
        if (playerTeam !== 0) {
            cls = (team === playerTeam ? "overlay-arrow-friendly" : "overlay-arrow-enemy");
        } else {
            cls = (team === 1 ? "overlay-arrow-friendly" : "overlay-arrow-enemy");
        }
        $(this.container).append('<svg class="overlay-arrow ' + cls + '" height="' + this.height + '" width="' + this.width + '">' + markerdef + line + '</svg>');
    };
    function Board(htmlTable, main) {
        this.cols = 0;
        this.rows = 0;
        this.playerTeam = 0;
        this.currentBoard = [];
        this.currentHistory = [];
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
            this.renderEndzones();
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
        this.currentHistory = history;
        this.overlay.clear();
        if (history.length > 0) {
            var lastMove = history[history.length - 1],
                m1 = lastMove[1],
                m2 = lastMove[2];
            this.overlay.renderMove(m1.Src.x, m1.Src.y, m1.Dst.x, m1.Dst.y, 1, this.playerTeam);
            this.overlay.renderMove(m2.Src.x, m2.Src.y, m2.Dst.x, m2.Dst.y, 2, this.playerTeam);
        }
    };
    Board.prototype.setActiveSquare = function (x, y) {
        if (x < 0 || y < 0 || x > this.cols || y > this.rows) {
            throw ("Invalid square: (" + x + "," + y + ")");
        }
        $(this.grid[x][y].container).addClass("grid-square-selectable");
    };
    Board.prototype.renderEndzones = function () {
        var x;
        for (x = 0; x < this.cols; x += 1) {
            $(this.grid[x][0].container).addClass("grid-square-endzone");
            $(this.grid[x][this.rows - 1].container).addClass("grid-square-endzone");
        }
        if (this.playerTeam === 1) {
            for (x = 0; x < this.cols; x += 1) {
                $(this.grid[x][0].container).addClass("grid-square-endzone-friendly");
                $(this.grid[x][this.rows - 1].container).addClass("grid-square-endzone-enemy");
            }
        }
        if (this.playerTeam === 2) {
            for (x = 0; x < this.cols; x += 1) {
                $(this.grid[x][0].container).addClass("grid-square-endzone-enemy");
                $(this.grid[x][this.rows - 1].container).addClass("grid-square-endzone-friendly");
            }
        }
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
        this.renderEndzones();
        // Create overlay
        this.overlay = new Overlay(this.htmlTable, rows, cols);
        console.log("Created grid");
    };
    Board.prototype.renderPieces = function (pieces) {
        var i, x, y, cls, num, width;
        this.currentBoard = pieces;
        width = $(this.grid[0][0].container).width();
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
                $(this.grid[x][y].dom).html('<svg class="' + cls + '"><circle cx="50%" cy="50%" r="' + (width * 0.4) + '"></circle><text x="50%" y="50%">' + num + '</text></svg>');
            }
        }
    };
    Board.prototype.handleGameOver = function (team) {
        if (team === 0) {
            this.overlay.displayDrawScreen();
        } else if (this.playerTeam !== 0 && team === this.playerTeam) {
            this.overlay.displayWinScreen();
        } else if (this.playerTeam !== 0 && team !== this.playerTeam) {
            this.overlay.displayLoseScreen();
        }
    };
    Board.prototype.onWindowResize = function () {
        console.log("Rerendering!");
        this.renderPieces(this.currentBoard);
        this.overlay.resize();
        this.renderHistory(this.currentHistory);
    };
    return Board;
}());