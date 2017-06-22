window.Board = (function () {
    "use strict";
    function Square() {
        this.unit = null;
        this.dom = null;
    }
    function Board(htmlTable) {
        this.cols = 0;
        this.rows = 0;
        this.htmlTable = htmlTable;
        this.grid = [];
    }
    Board.prototype.render = function (json) {
        var board = JSON.parse(json);
        console.log(board);
        if (board.Board.length !== (board.Cols * board.Rows)) {
            throw ("Invalid board data");
        }
        if (this.cols !== board.Cols || this.rows !== board.Rows) {
            this.createGrid(board.Cols, board.Rows);
        }
        this.renderPieces(board.Board);
    };
    Board.prototype.createGrid = function (cols, rows) {
        var x, y, tr, td;
        this.cols = cols;
        this.rows = rows;
        // Create matrix
        this.grid = [];
        for (x = 0; x < cols; x += 1) {
            this.grid.push([]);
            for (y = 0; y < rows; y += 1) {
                this.grid[x].push(new Square());
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
                tr.appendChild(td);
                this.grid[x][y].dom = td;
            }
            this.htmlTable.appendChild(tr);
        }
        console.log("Created grid");
    };
    Board.prototype.renderPieces = function (pieces) {
        console.log(pieces);
    };
    return Board;
}());