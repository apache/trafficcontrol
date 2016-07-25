var DateRangeFilter = function() {
    return function(items, dateRange) {
        if (_.isUndefined(dateRange)) return items;
        var filteredItems = _.filter(items, function(item) {
            return (
                ( moment(item.date).isAfter(dateRange.start, 'day') || moment(item.date).isSame(dateRange.start, 'day') ) &&
                ( moment(item.date).isBefore(dateRange.end, 'day') || moment(item.date).isSame(dateRange.end, 'day') )
            );
        });
        return filteredItems;
    };
};

DateRangeFilter.$inject = [];
module.exports = DateRangeFilter;
