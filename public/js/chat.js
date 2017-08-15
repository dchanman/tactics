window.Chat = (function () {
    "use strict";
    function Chat(main) {
        this.main = main;
    }
    Chat.prototype.receiveMessage = function (msgEvent) {
        console.log(msgEvent);
        var log, sender, p;
        log = document.createElement("div");
        sender = msgEvent.sender + ": ";
        p = document.createElement("p");
        $(p).addClass("chat-message");
        $(p).html(sender.bold());
        p.appendChild(document.createTextNode(msgEvent.message));
        log.appendChild(p);
        $("#chatlog").append(log);
        $("#chatlog").animate({ scrollTop: $('#chatlog').prop("scrollHeight")}, 100);
    };
    Chat.prototype.sendMessage = function () {
        var text = $("#chatmsg").val();
        if (text) {
            this.main.api.sendChat(text + "\n");
        }
        $("#chatmsg").val("");
    };
    Chat.prototype.notification = function (msg) {
        var log, p;
        log = document.createElement("div");
        p = document.createElement("p");
        $(p).html(msg.italics().bold());
        log.appendChild(p);
        $("#chatlog").append(log);
        $("#chatlog").animate({ scrollTop: $('#chatlog').prop("scrollHeight")}, 100);
    };
    return Chat;
}());
