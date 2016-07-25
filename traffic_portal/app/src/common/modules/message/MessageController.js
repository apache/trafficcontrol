var MessageController = function($rootScope, $scope, messageModel, chartModel) {

    $scope.messageData = messageModel;

    $scope.chartData = chartModel.chart;

    $scope.dismissMessage = function(message) {
        messageModel.removeMessage(message);
    };

    $scope.showConnectionLostMsg = function() {
        $('#lostConnectionMsg').show();
    };

    $scope.hideConnectionLostMsg = function() {
        $('#lostConnectionMsg').hide();
    };

    $rootScope.$watch('online', function(newStatus) {
        if (newStatus === false && $scope.chartData.autoRefresh === true) {
            $scope.showConnectionLostMsg();
        }
    });

};

MessageController.$inject = ['$rootScope', '$scope', 'messageModel', 'chartModel'];
module.exports = MessageController;
