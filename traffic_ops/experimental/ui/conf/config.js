module.exports = {
    timeout: '60s',
    port: 8080,
    proxyPort: 8009,
    api: {
        base_url: 'http://ipcdn-cache-12.cdnlab.comcast.net:8080',
        key: ''
    },
    files: {
        static: './app/dist/public/'
    },
    log: {
        stream: './server/log/access.log'
    }
};
