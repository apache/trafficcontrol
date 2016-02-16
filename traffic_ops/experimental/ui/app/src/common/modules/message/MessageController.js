var MessageController = function($scope, messageModel) {

    $scope.messageData = messageModel;

    $scope.dismissMessage = function(message) {
        messageModel.removeMessage(message);
    };

};

MessageController.$inject = ['$scope', 'messageModel'];
module.exports = MessageController;
