window.Chat = (function () {
    "use strict";
    function Chat(main) {
        this.main = main;
    }
    Chat.prototype.receiveMessage = function (msgEvent) {
        var log, sender, sendernode, p;
        log = document.createElement("div");
        sender = msgEvent.sender + ": ";
        p = document.createElement("p");
        $(p).addClass("chat-message");
        sendernode = document.createTextNode(sender);
        $(p).append(sendernode);
        p.appendChild(document.createTextNode(msgEvent.message));
        log.appendChild(p);
        $("#chatlog").append(log);
        $("#chatlog").animate({ scrollTop: $('#chatlog').prop("scrollHeight")}, 100);
    };
    Chat.prototype.getName = function () {
        this.main.api.getChatName()
            .then(function (result) {
                $("#chatname").val(result.name);
            })
            .catch(function (err) {
                console.log(err);
            });
    };
    Chat.prototype.setName = function () {
        var self = this,
            name = $("#chatname").val();
        console.log("Setting name: " + name);
        if (name) {
            this.main.api.setChatName(name)
                .catch(function (err) {
                    console.log(err);
                    self.getName();
                });
        }
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