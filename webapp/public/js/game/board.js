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
        this.showLastUnit = false;
        this.showLastMove = false;
        this.showCollisions = false;
        //
        this.engineBoard = null;
        this.engineBoardHistory = [];
        this.engineBoardCols = 0;
        this.engineBoardRows = 0;
    }
    function createInitialEngineBoard(update) {
        var mappedUnits = [], i, ptr, engineBoard;
        for (i = 0; i < update.board.Board.length; i += 1) {
            ptr = {
                Team: update.board.Board[i].Team,
                Stack: update.board.Board[i].Stack,
                Exists: update.board.Board[i].Exists
            };
            ptr.$val = ptr;
            mappedUnits.push(ptr);
        }
        engineBoard = Engine.NewBoardFromBoard(update.board.Cols, update.board.Rows, mappedUnits);
        return engineBoard;
    }
    Board.prototype.runEngine = function (update) {
        if (this.engineBoardHistory.length > update.history.length) {
            // The game was reset, force redraw of board
            this.engineBoard = null;
        }
        this.engineBoardHistory = update.history;
        if (this.engineBoard === null) {
            this.engineBoard = createInitialEngineBoard(update);
            console.log("Created initial board");
            this.render(this.engineBoard.GetBoard());
        } else {
            var up, m1, m2, move1, move2;
            if (update.history.length < 1) {
                return;
            }
            up = update.history[update.history.length - 1];
            m1 = up.moves[1];
            m2 = up.moves[2];
            move1 = Engine.NewMove(m1.Src.X, m1.Src.Y, m1.Dst.X, m1.Dst.Y);
            move2 = Engine.NewMove(m2.Src.X, m2.Src.Y, m2.Dst.X, m2.Dst.Y);
            console.log({"resolution": this.engineBoard.ResolveMove(move1, move2)});
            this.render(this.engineBoard.GetBoard());
        }
    };
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
        console.log("Rendering!");
        console.log(board);
        if (board.Board.length !== (board.Cols * board.Rows)) {
            throw ("Invalid board data");
        }
        if (this.cols !== board.Cols || this.rows !== board.Rows) {
            this.createGrid(board.Cols, board.Rows);
        }
        this.renderPieces(board.Board);
    };
    Board.prototype.renderHistory = function (history) {
        if (this.overlay === null) {
            return;
        }
        this.currentHistory = history;
        this.overlay.clear();
        if (history.length > 0) {
            var i,
                lastMove = history[history.length - 1],
                m1 = lastMove.moves[1],
                m2 = lastMove.moves[2];
            if (this.showLastUnit) {
                this.overlay.renderPiece(m1.Src.X, m1.Src.Y, Unit.fromJSON(lastMove.oldUnits[1]));
                this.overlay.renderPiece(m2.Src.X, m2.Src.Y, Unit.fromJSON(lastMove.oldUnits[2]));
            }
            if (this.showLastMove) {
                this.overlay.renderMove(m1.Src.X, m1.Src.Y, m1.Dst.X, m1.Dst.Y, 1, this.playerTeam);
                this.overlay.renderMove(m2.Src.X, m2.Src.Y, m2.Dst.X, m2.Dst.Y, 2, this.playerTeam);
            }
            if (this.showCollisions) {
                for (i = 0; i < lastMove.collisions.length; i += 1) {
                    this.overlay.renderCollision(lastMove.collisions[i].X, lastMove.collisions[i].Y);
                }
            }
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
    };
    Board.prototype.renderPieces = function (pieces) {
        var i, x, y, width;
        this.currentBoard = pieces;
        width = $(this.grid[0][0].container).width();
        for (i = 0; i < pieces.length; i += 1) {
            x = Math.floor(i / this.rows);
            y = i % this.rows;
            this.grid[x][y].unit = null;
            $(this.grid[x][y].dom).html("");
            $(this.grid[x][y].container).removeClass("piece-friendly");
            $(this.grid[x][y].container).removeClass("piece-enemy");
            if (pieces[i].Exists) {
                this.grid[x][y].unit = Unit.fromJSON(pieces[i]);
                $(this.grid[x][y].dom).html(this.grid[x][y].unit.getRenderHtml(width));
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
        this.renderPieces(this.currentBoard);
        this.overlay.resize();
        this.renderHistory(this.currentHistory);
    };
    Board.prototype.setOverlaySettings = function (showLastMove, showLastUnit, showCollisions) {
        this.showLastUnit = showLastUnit;
        this.showLastMove = showLastMove;
        this.showCollisions = showCollisions;
        this.renderHistory(this.currentHistory);
    };
    return Board;
}());