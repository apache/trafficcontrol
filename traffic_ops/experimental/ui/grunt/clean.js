
module.exports = {
    all: [
        '<%= globalConfig.distdir %>/*',
        '<%= globalConfig.tmpdir %>/*'
    ],
    options: {
        force: true
    }
};