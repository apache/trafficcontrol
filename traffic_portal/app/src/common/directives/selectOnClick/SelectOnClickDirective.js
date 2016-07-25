var SelectOnClickDirective = function() {
    return {
        restrict: 'A',
        link: function (scope, element, attrs) {
            element.on('click', function () {
                this.select();
            });
        }
    };
};

module.exports = SelectOnClickDirective;
