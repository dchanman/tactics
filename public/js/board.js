window.Board = (function () {
    "use strict";
    function Square(x, y) {
        this.x = x;
        this.y = y;
        this.unit = null;
        this.dom = null;
    }
    Square.prototype.setDOM = function (td, div) {
        var self = this;
        this.dom = div;
        $(td).click(function () {
            console.log("Clicked " + self.x + "," + self.y);
            console.log(self.unit);
        });
    };
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
        var x, y, tr, td, div;
        this.cols = cols;
        this.rows = rows;
        // Create matrix
        this.grid = [];
        for (x = 0; x < cols; x += 1) {
            this.grid.push([]);
            for (y = 0; y < rows; y += 1) {
                this.grid[x].push(new Square(x, y));
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
        var i, x, y;
        for (i = 0; i < pieces.length; i += 1) {
            x = Math.floor(i / this.cols);
            y = i % this.cols;
            this.grid[x][y].unit = pieces[i];
            $(this.grid[x][y].dom).html(pieces[i].Name);
        }
    };
    return Board;
}());