var UnitsFilter = function(numberUtils) {
    return function(kilobits, prettify) {
        var friendlyUnit = numberUtils.shrink(kilobits);
        return (prettify) ? numberUtils.addCommas(friendlyUnit[0]) + ' ' + friendlyUnit[1] : friendlyUnit[0];
    };
};

UnitsFilter.$inject = ['numberUtils'];
module.exports = UnitsFilter;
