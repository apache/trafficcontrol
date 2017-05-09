var PercentFilter = function() {
	return function(input) {
		input = parseFloat(input);
		input *= 100;
		if(input % 1 === 0) {
			input = input.toFixed(0);
		}
		else {
			input = input.toFixed(2);
		}
		return input + '%';
	};
};

PercentFilter.$inject = [];
module.exports = PercentFilter;
