window.Board = (function () {
    "use strict";
    function Square(board, x, y) {
        this.board = board;
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
            var u = new Unit();
            u.name = "sup";
            self.board.main.api.addUnit(self.x, self.y, u)
                .then(function () {
                    console.log("Added unit successfully!");
                })
                .catch(function (err) {
                    console.log(err);
                });
        });
    };
    function Board(htmlTable, main) {
        this.cols = 0;
        this.rows = 0;
        this.htmlTable = htmlTable;
        this.grid = [];
        this.main = main;
    }
    Board.prototype.render = function (board) {
        if (board.board.length !== (board.cols * board.rows)) {
            throw ("Invalid board data");
        }
        if (this.cols !== board.cols || this.rows !== board.rows) {
            this.createGrid(board.cols, board.rows);
        }
        this.renderPieces(board.board);
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
        var i, x, y;
        for (i = 0; i < pieces.length; i += 1) {
            x = Math.floor(i / this.rows);
            y = i % this.rows;
            this.grid[x][y].unit = null;
            if (pieces[i].exists) {
                this.grid[x][y].unit = pieces[i];
                $(this.grid[x][y].dom).html(pieces[i].name);
            }
        }
    };
    return Board;
}());