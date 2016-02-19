var OffsetFilter = function() {
    return function(input, start) {
        if($.isArray(input)) {
            start = parseInt(start, 10);
            return input.slice(start);
        }
    };
};

OffsetFilter.$inject = ['$log'];
module.exports = OffsetFilter;
