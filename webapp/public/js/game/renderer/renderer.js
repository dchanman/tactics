window.Renderer = (function () {
    "use strict";
    var colLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    function xCoord(x) {
        return colLetters[x];
    }
    function yCoord(y, nRows) {
        return (nRows - y).toString();
    }
    function moveToString(move, nRows) {
        var src = xCoord(move.Src.X) + yCoord(move.Src.Y, nRows),
            dst = xCoord(move.Dst.X) + yCoord(move.Dst.Y, nRows);
        return src + "→" + dst;
    }
    function Board(board, turn, resolution) {
        this.board = board;
        this.turn = turn;
        this.resolution = resolution;
    }
    function Renderer(main, htmlTable, historyTable) {
        this.main = main;
        // HTML DOM
        this.htmlTable = htmlTable;
        this.historyTable = historyTable;
        this.historyTableRowDOMs = [];
        this.turnNumber = -1;
        this.grid = [];
        this.overlay = null;
        this.selectedSquare = null;
        // Game state
        this.cols = 0;
        this.rows = 0;
        this.playerTeam = 0;
        // Configuration
        this.showLastMove = false;
        this.showCollisions = false;
        // Internal engine
        this.engineBoard = null;
        this.historyBoards = [];
    }
    Renderer.prototype.engineResolveMove = function (turn) {
        var m1, m2, move1, move2;
        m1 = turn[1];
        m2 = turn[2];
        move1 = Engine.NewMove(m1.Src.X, m1.Src.Y, m1.Dst.X, m1.Dst.Y);
        move2 = Engine.NewMove(m2.Src.X, m2.Src.Y, m2.Dst.X, m2.Dst.Y);
        return this.engineBoard.ResolveMove(move1, move2);
    };
    Renderer.prototype.handleGameTurn = function (history) {
        var turn, resolution;
        if (history.length < 1) {
            return;
        }
        turn = history[history.length - 1];
        resolution = this.engineResolveMove(turn);
        this.historyBoards.push(new Board(this.engineBoard.GetBoard(), turn, resolution));
        this.renderHistory(history, this.rows);
        this.selectTurn(history.length);
    };
    Renderer.prototype.engineInit = function (gameInformation) {
        var i, resolution, clickEnabled;
        this.engineBoard = Engine.newEngineBoard(gameInformation.board);
        clickEnabled = (gameInformation.history.length === 0);
        this.historyBoards = [new Board(this.engineBoard.GetBoard(), null, null, clickEnabled)];
        for (i = 0; i < gameInformation.history.length; i += 1) {
            resolution = this.engineResolveMove(gameInformation.history[i]);
            clickEnabled = (i === gameInformation.history.length - 1);
            this.historyBoards.push(new Board(this.engineBoard.GetBoard(), gameInformation.history[i], resolution));
        }
        this.renderHistory(gameInformation.history, gameInformation.board.Rows);
        this.selectTurn(gameInformation.history.length);
    };
    Renderer.prototype.setPlayerTeam = function (team) {
        if (this.playerTeam !== team) {
            this.playerTeam = team;
            this.renderEndzones();
            this.selectTurn(this.turnNumber);
        }
    };
    Renderer.prototype.removeSelectableSquares = function () {
        var i, j;
        for (i = 0; i < this.grid.length; i += 1) {
            for (j = 0; j < this.grid[i].length; j += 1) {
                $(this.grid[i][j].container).removeClass("grid-square-selectable");
                $(this.grid[i][j].container).removeClass("grid-square-commit-src");
                $(this.grid[i][j].container).removeClass("grid-square-commit-dst");
                $(this.grid[i][j].container).droppable({disabled: true});
            }
        }
    };
    Renderer.prototype.render = function (board) {
        if (board.Board.length !== (board.Cols * board.Rows)) {
            throw ("Invalid board data");
        }
        if (this.cols !== board.Cols || this.rows !== board.Rows) {
            this.createGrid(board.Cols, board.Rows);
        }
        this.renderPieces(board.Board);
    };
    Renderer.prototype.renderLastMove = function (turn, resolution) {
        if (this.overlay === null) {
            return;
        }
        this.overlay.clear();
        var m1 = turn[1],
            m2 = turn[2],
            i;
        if (this.showLastMove) {
            this.overlay.renderMove(m1.Src.X, m1.Src.Y, m1.Dst.X, m1.Dst.Y, 1, this.playerTeam);
            this.overlay.renderMove(m2.Src.X, m2.Src.Y, m2.Dst.X, m2.Dst.Y, 2, this.playerTeam);
        }
        if (this.showCollisions) {
            for (i = 0; i < resolution.Collisions.length; i += 1) {
                this.overlay.renderCollision(resolution.Collisions[i].X, resolution.Collisions[i].Y);
            }
        }
    };
    function historyRowOnClickHandler(renderer, turnNumber) {
        return function () {
            renderer.selectTurn(turnNumber);
        };
    }
    Renderer.prototype.selectTurn = function (turnNumber) {
        var i, j, board, clickEnabled;
        if (turnNumber < 0) {
            return;
        }
        this.turnNumber = turnNumber;
        for (i = 0; i < this.historyTableRowDOMs.length; i += 1) {
            this.historyTableRowDOMs[i].removeClass("table-active");
        }
        this.historyTableRowDOMs[turnNumber].addClass("table-active");
        board = this.historyBoards[turnNumber];
        this.render(board.board);
        this.overlay.clear();
        if (board.turn !== null && board.resolution !== null) {
            this.renderLastMove(board.turn, board.resolution);
        }
        clickEnabled = (turnNumber === this.historyTableRowDOMs.length - 1);
        for (i = 0; i < this.grid.length; i += 1) {
            for (j = 0; j < this.grid[i].length; j += 1) {
                this.grid[i][j].setClickEnabled(clickEnabled);
            }
        }
    };
    Renderer.prototype.renderHistory = function (history, boardRows) {
        this.currentRendereredHistory = history;
        this.historyTableRowDOMs = [];
        $(this.historyTable).html("");
        var i, tr;
        tr = $("<tr></tr>")
            .append('<th scope="row">0</th>')
            .append('<td>-</td>')
            .append('<td>-</td>');
        tr.click(historyRowOnClickHandler(this, 0));
        $(this.historyTable).append(tr);
        this.historyTableRowDOMs.push(tr);
        for (i = 0; i < history.length; i += 1) {
            tr = $("<tr></tr>")
                .append('<th scope="row">' + (i + 1) + '</th>')
                .append('<td>' + moveToString(history[i][1], boardRows) + '</td>')
                .append('<td>' + moveToString(history[i][2], boardRows) + '</td>');
            tr.click(historyRowOnClickHandler(this, i + 1));
            $(this.historyTable).append(tr);
            this.historyTableRowDOMs.push(tr);
            $("#history-container").animate({ scrollTop: $("#history-container").prop("scrollHeight")}, 100);
        }
    };
    Renderer.prototype.setActiveSquare = function (x, y) {
        if (x < 0 || y < 0 || x > this.cols || y > this.rows) {
            throw ("Invalid square: (" + x + "," + y + ")");
        }
        var sq = this.grid[x][y];
        $(sq.container).addClass("grid-square-selectable");
        $(sq.container).droppable({
            disabled: false,
            hoverClass: "grid-square-selectable-hovered",
            drop: function () {
                sq.click();
            }
        });
    };
    Renderer.prototype.renderEndzones = function () {
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
    Renderer.prototype.createGrid = function (cols, rows) {
        var x, y, tr, td, div, leftMargin, botMargin;
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
            leftMargin = $('<td class="grid-coord">' + yCoord(y, this.rows) + "</td>");
            tr = $("<tr></tr>")
                .addClass("grid")
                .append(leftMargin);
            for (x = 0; x < cols; x += 1) {
                td = $("<td></td>")
                    .addClass("grid-square")
                    .appendTo(tr);
                div = $("<div></div>")
                    .addClass("grid-square-container")
                    .appendTo(td);
                this.grid[x][y].setDOM(td[0], div[0]);
            }
            $(this.htmlTable).append(tr);
        }
        tr = $("<tr></tr>")
            .addClass("grid");
        for (x = 0; x <= cols; x += 1) {
            botMargin = $('<td class="grid-coord">' + (x > 0 ? xCoord(x - 1) : "") + "</td>");
            tr.append(botMargin);
        }
        $(this.htmlTable).append(tr);
        this.renderEndzones();
        // Create overlay
        this.overlay = new Overlay(this.htmlTable, rows, cols,
            leftMargin.outerWidth(), botMargin.outerHeight());
    };
    Renderer.prototype.renderPieces = function (pieces) {
        var i, x, y, width;
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
    Renderer.prototype.handleGameOver = function (team) {
        if (team === 0) {
            this.overlay.displayDrawScreen();
        } else if (this.playerTeam !== 0 && team === this.playerTeam) {
            this.overlay.displayWinScreen();
        } else if (this.playerTeam !== 0 && team !== this.playerTeam) {
            this.overlay.displayLoseScreen();
        }
    };
    Renderer.prototype.onWindowResize = function () {
        if (this.overlay !== null) {
            this.overlay.resize();
        }
        this.selectTurn(this.turnNumber);
    };
    Renderer.prototype.setOverlaySettings = function (showLastMove, showCollisions, moveConfirmation) {
        this.showLastMove = showLastMove;
        this.showCollisions = showCollisions;
        Square.moveConfirmation = moveConfirmation;
        this.selectTurn(this.turnNumber);
    };
    return Renderer;
}());