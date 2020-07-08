exports.config.onPrepare = () => {
  const jasmineReporters = require('jasmine-reporters');
  jasmine.getEnv().addReporter(
    new jasmineReporters.JUnitXmlReporter({
      savePath: '/portaltestresults',
      filePrefix: 'portaltestresults',
      consolidateAll: true
    }))
}
