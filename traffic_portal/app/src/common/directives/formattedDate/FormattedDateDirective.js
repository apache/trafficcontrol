var FormattedDateDirective = function(dateUtils) {
    return {
        link: function (scope, element, attrs, ctrl) {
            ctrl.$formatters.unshift(function (modelValue) {
                return dateUtils.dateFormat(modelValue, "mmm d yyyy h:MM tt (Z)")
            });

            ctrl.$parsers.unshift(function (viewValue) {
                return dateUtils.dateFormat(modelValue, "mmm d yyyy h:MM tt (Z)")
            });
        },
        restrict: 'A',
        require: 'ngModel'
    }
};

module.exports = FormattedDateDirective;
