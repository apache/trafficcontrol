var MessageModel = function($rootScope) {

    var model = this;
    var queue = [];

    var messages = {};
    messages.loaded = false;
    messages.content = [];

    this.messages = messages;

    this.setMessages = function(messagesArray, delay) {
        // delay should be true if a redirect follows...
        if (!angular.isArray(messagesArray)) return;
        var messages = {};
        messages.loaded = true;
        messages.content = messagesArray;
        if (delay) {
            queue[0] = messages; // queue up messages to display after a location change
        } else {
            model.messages = messages; // show the messages asap
            queue = []; // clear the queue as messages will be shown immediately
        }
    };

    this.resetMessages = function() {
        messages = {};
        messages.loaded = false;
        messages.content = [];

        this.messages = messages;
    };

    this.removeMessage = function(message) {
        model.messages.content = _.without(model.messages.content, message);
    };

    $rootScope.$on("$locationChangeStart", function() {
        model.resetMessages();
    });

    $rootScope.$on("$locationChangeSuccess", function() {
        model.messages = queue.shift() || {};
    });

};

MessageModel.$inject = ['$rootScope'];
module.exports = MessageModel;
