module.exports = angular.module('trafficOps.private.monitor.dashboard.view', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficOps.private.monitor.dashboard.view', {
				url: '',
				views: {
					cacheGroupsContent: {
						templateUrl: 'common/modules/widget/cacheGroups/widget.cacheGroups.tpl.html',
						controller: 'WidgetCacheGroupsController',
						resolve: {
							cacheGroupHealth: function() {
								// this is already defined in a parent template that shares the $scope
								return null;
							}
						}
					},
					capacityContent: {
						templateUrl: 'common/modules/widget/capacity/widget.capacity.tpl.html',
						controller: 'WidgetCapacityController'
					},
					cdnChartContent: {
						templateUrl: 'common/modules/widget/cdnChart/widget.cdnChart.tpl.html',
						controller: 'WidgetCDNChartController',
						resolve: {
							cdn: function() {
								// the controller will take care of fetching the cdn
								return null;
							}
						}
					},
					changeLogsContent: {
						templateUrl: 'common/modules/widget/changeLogs/widget.changeLogs.tpl.html',
						controller: 'WidgetChangeLogsController',
						resolve: {
							changeLogs: function(changeLogService) {
								return changeLogService.getChangeLogs({ limit: 5 });
							}
						}
					},
					routingContent: {
						templateUrl: 'common/modules/widget/routing/widget.routing.tpl.html',
						controller: 'WidgetRoutingController',
						resolve: {
							routing: function() {
								return [];
							}
						}
					},
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
