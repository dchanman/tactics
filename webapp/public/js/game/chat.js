window.Chat = (function () {
    "use strict";
    function Chat(main) {
        this.main = main;
    }
    Chat.prototype.receiveMessage = function (msgEvent) {
        var sender;
        sender = msgEvent.sender + ": ";
        $("<div>")
            .append($("<p>", {"class": "chat-message"})
                .append($("<b>")
                    .append(document.createTextNode(sender)))
                .append(document.createTextNode(msgEvent.message)))
            .appendTo("#chatlog");
        $("#chatlog").animate({ scrollTop: $('#chatlog').prop("scrollHeight")}, 100);
    };
    Chat.prototype.receiveJoin = function (msgEvent) {
        this.notification(msgEvent + " has joined the game.");
    };
    Chat.prototype.receiveNameChange = function (msgEvent) {
        this.notification(msgEvent.from + " is now " + msgEvent.to + ".");
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
        $("<div>")
            .append(
                $("<div>", {"class": "chat-notification"}).append(
                    document.createTextNode(msg)
                )
            )
            .appendTo("#chatlog");
        $("#chatlog").animate({ scrollTop: $('#chatlog').prop("scrollHeight")}, 100);
    };
    return Chat;
}());