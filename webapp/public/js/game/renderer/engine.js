(function () {
    "use strict";
    if (!Engine) {
        throw "The gopherjs 'Engine' entity must be imported first";
    }
    Engine.newEngineBoard = function (board) {
        var mappedUnits = [], i, ptr, engineBoard;
        for (i = 0; i < board.Board.length; i += 1) {
            ptr = {
                Team: board.Board[i].Team,
                Stack: board.Board[i].Stack,
                Exists: board.Board[i].Exists
            };
            ptr.$val = ptr;
            mappedUnits.push(ptr);
        }
        engineBoard = Engine.NewBoardFromBoard(board.Cols, board.Rows, mappedUnits);
        return engineBoard;
    };
}());
