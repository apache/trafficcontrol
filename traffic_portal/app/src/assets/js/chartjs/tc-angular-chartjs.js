/**
 * tc-angular-chartjs - http://carlcraig.github.io/tc-angular-chartjs/
 * Copyright (c) 2017 Carl Craig
 * Dual licensed with the Apache-2.0 or MIT license.
 */
;(function(root, factory) {
  if (typeof define === 'function' && define.amd) {
    define(['angular', 'chart.js'], factory);
  } else if (typeof exports === 'object') {
    module.exports = factory(require('angular'), require('chart.js'));
  } else {
    root.tcAngularChartjs = factory(root.angular, root.Chart);
  }
}(this, function(angular, Chart) {

TcChartjs.$inject = ["TcChartjsFactory"];
TcChartjsLine.$inject = ["TcChartjsFactory"];
TcChartjsBar.$inject = ["TcChartjsFactory"];
TcChartjsHorizontalBar.$inject = ["TcChartjsFactory"];
TcChartjsRadar.$inject = ["TcChartjsFactory"];
TcChartjsPolararea.$inject = ["TcChartjsFactory"];
TcChartjsPie.$inject = ["TcChartjsFactory"];
TcChartjsDoughnut.$inject = ["TcChartjsFactory"];
TcChartjsBubble.$inject = ["TcChartjsFactory"];angular
  .module('tc.chartjs', [])
  .directive('tcChartjs', TcChartjs)
  .directive('tcChartjsLine', TcChartjsLine)
  .directive('tcChartjsBar', TcChartjsBar)
  .directive('tcChartjsHorizontalbar', TcChartjsHorizontalBar)
  .directive('tcChartjsRadar', TcChartjsRadar)
  .directive('tcChartjsPolararea', TcChartjsPolararea)
  .directive('tcChartjsPie', TcChartjsPie)
  .directive('tcChartjsDoughnut', TcChartjsDoughnut)
  .directive('tcChartjsBubble', TcChartjsBubble)
  .directive('tcChartjsLegend', TcChartjsLegend)
  .factory('TcChartjsFactory', TcChartjsFactory);

function TcChartjs(TcChartjsFactory) {
  return new TcChartjsFactory();
}

function TcChartjsLine(TcChartjsFactory) {
  return new TcChartjsFactory('line');
}

function TcChartjsBar(TcChartjsFactory) {
  return new TcChartjsFactory('bar');
}

function TcChartjsHorizontalBar(TcChartjsFactory) {
  return new TcChartjsFactory('horizontalbar');
}

function TcChartjsRadar(TcChartjsFactory) {
  return new TcChartjsFactory('radar');
}

function TcChartjsPolararea(TcChartjsFactory) {
  return new TcChartjsFactory('polararea');
}

function TcChartjsPie(TcChartjsFactory) {
  return new TcChartjsFactory('pie');
}

function TcChartjsDoughnut(TcChartjsFactory) {
  return new TcChartjsFactory('doughnut');
}

function TcChartjsBubble(TcChartjsFactory) {
  return new TcChartjsFactory('bubble');
}

function TcChartjsFactory() {

  return function (chartType) {

    return {
      restrict: 'A',
      scope: {
        data: '=chartData',
        options: '=chartOptions',
        plugins: '=chartPlugins',
        type: '@chartType',
        legend: '=?chartLegend',
        chart: '=?chart',
        click: '&chartClick'
      },
      link: link
    };

    function link($scope, $elem, $attrs) {
      var ctx = $elem[0].getContext('2d');
      var chartObj;
      var showLegend = false;
      var autoLegend = false;
      var exposeChart = false;
      var legendElem = null;

      for (var attr in $attrs) {
        if (attr === 'chartLegend') {
          showLegend = true;
        } else if (attr === 'chart') {
          exposeChart = true;
        } else if (attr === 'autoLegend') {
          autoLegend = true;
        }
      }

      $scope.$on('$destroy', function() {
        if (chartObj && typeof chartObj.destroy === 'function') {
          chartObj.destroy();
        }
      });

      if ($scope.click) {
        $elem[0].onclick = function(evt) {
          var out = {
            chartEvent: evt,
            element: chartObj.getElementAtEvent(evt),
            elements: chartObj.getElementsAtEvent(evt),
            dataset: chartObj.getDatasetAtEvent(evt)
          };

          $scope.click({event: out});
        };
      }

      $scope.$watch('[data, options, plugins]', function (value) {
        if (value && $scope.data) {
          if (chartObj && typeof chartObj.destroy === 'function') {
            chartObj.destroy();
          }

          var type = chartType || $scope.type;
          if (!type) {
            throw 'Error creating chart: Chart type required.';
          }
          type = cleanChartName(type);

          chartObj = new Chart(ctx, {
            type: type,
            data: angular.copy($scope.data),
            options: $scope.options,
            plugins: $scope.plugins
          });
          
          if (showLegend) {
            $scope.legend = chartObj.generateLegend();
          }

          if (autoLegend) {
            if (legendElem) {
              legendElem.remove();
            }
            angular.element($elem[0]).after(chartObj.generateLegend());
            legendElem = angular.element($elem[0] ).next();
          }

          if (exposeChart) {
            $scope.chart = chartObj;
          }
          chartObj.resize();
        }
      }, true);
    }

    function cleanChartName(type) {
      var typeLowerCase = type.toLowerCase();
      switch (typeLowerCase) {
        case 'polararea':
          return 'polarArea';
        case 'horizontalbar':
          return 'horizontalBar';
        default:
          return type;
      }
    }

  };
}

function TcChartjsLegend() {
  return {
    restrict: 'A',
    scope: {
      legend: '=?chartLegend'
    },
    link: link
  };

  function link($scope, $elem) {
    $scope.$watch('legend', function (value) {
      if (value) {
        $elem.html(value);
      }
    }, true);
  }
}

return TcChartjsFactory;
}));
