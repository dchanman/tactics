window.Main = (function () {
    "use strict";
    console.log("Hello World!");
    var Main = function () {
        this.board = new Board(document.getElementById("game"));
    };
    Main.prototype.refresh = function () {
        var xmlHttp = new XMLHttpRequest(),
            board = this.board;
        xmlHttp.onreadystatechange = function () {
            if (xmlHttp.readyState === 4 && xmlHttp.status === 200) {
                board.render(xmlHttp.responseText);
            }
        };
        xmlHttp.open("GET", "/game", true);
        xmlHttp.send(null);
    };
    Main.prototype.start = function () {
        this.refresh();
    };
    return Main;
}());
var main = new Main();
main.start();
