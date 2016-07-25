var rcSubmitDirective = {
    'rcSubmit': ['$parse', function ($parse) {
        return {
            restrict: 'A',
            require: ['rcSubmit', '?form'],
            controller: ['$scope', function ($scope) {

                var formController = null;
                var attemptHandlers = [];

                this.attempted = false;

                this.onAttempt = function(handler) {
                    attemptHandlers.push(handler);
                };

                this.setAttempted = function() {
                    this.attempted = true;

                    angular.forEach(attemptHandlers, function (handler) {
                        handler();
                    });
                };

                this.setFormController = function(controller) {
                    formController = controller;
                };

                this.needsAttention = function (fieldModelController) {
                    if (!formController) return false;

                    if (fieldModelController) {
                        return fieldModelController.$invalid &&
                            (fieldModelController.$dirty || this.attempted);
                    } else {
                        return formController && formController.$invalid &&
                            (formController.$dirty || this.attempted);
                    }
                };
            }],
            compile: function(cElement, cAttributes, transclude) {
                return {
                    pre: function(scope, formElement, attributes, controllers) {

                        var submitController = controllers[0];
                        var formController = (controllers.length > 1) ? controllers[1] : null;

                        submitController.setFormController(formController);

                        scope.rc = scope.rc || {};
                        scope.rc[attributes.name] = submitController;
                    },
                    post: function(scope, formElement, attributes, controllers) {

                        var submitController = controllers[0];
                        var formController = (controllers.length > 1) ? controllers[1] : null;
                        var fn = $parse(attributes.rcSubmit);

                        formElement.bind('submit', function () {
                            submitController.setAttempted();
                            if (!scope.$$phase) scope.$apply();

                            if (!formController.$valid) return false;

                            scope.$apply(function() {
                                fn(scope, {$event:event});
                            });
                        });
                    }
                };
            }
        };
    }]
};