var ChangeLogModel = function($rootScope, $interval, changeLogService, userModel) {

	var newLogCount = 0,
		pollingIntervalInSecs = 30,
		changeLogInterval;

	this.newLogCount = function() {
		return newLogCount;
	};

	var createChangeLogInterval = function() {
		killChangeLogInterval();
		changeLogInterval = $interval(function() { getNewLogCount() }, (pollingIntervalInSecs*1000)); // every X minutes
	};

	var killChangeLogInterval = function() {
		if (angular.isDefined(changeLogInterval)) {
			$interval.cancel(changeLogInterval);
			changeLogInterval = undefined;
		}
	};

	var getNewLogCount = function() {
		changeLogService.getNewLogCount()
			.then(function(result) {
				newLogCount = result.data.response.newLogcount;
			});
	};

	$rootScope.$on('authService::login', function() {
		getNewLogCount();
		createChangeLogInterval();
	});

	$rootScope.$on('authService::logout', function() {
		killChangeLogInterval();
	});

	$rootScope.$on('changeLogService::getChangeLogs', function() {
		newLogCount = 0;
	});

	var init = function () {
		if (userModel.loaded) {
			getNewLogCount();
			createChangeLogInterval();
		}
	};
	init();

};

ChangeLogModel.$inject = ['$rootScope', '$interval', 'changeLogService', 'userModel'];
module.exports = ChangeLogModel;