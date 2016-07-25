var rcVerifySetDirective = {
    'rcVerifySet': function () {
        return {
            restrict: 'A',
            require: ['^rcSubmit', 'ngModel'],
            link: function (scope, element, attributes, controllers) {
                var submitController = controllers[0];
                var modelController = controllers[1];

                submitController.onAttempt(function() {
                    modelController.$setViewValue(element.val());
                });
            }
        };
    }
};