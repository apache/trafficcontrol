module.exports = {
    timeout: '60s',
    port: 8080,
    proxyPort: 8009,
    api: {
        url: 'http://localhost:3000/api/2.0/',
        key: ''
    },
    files: {
        static: './app/dist/public/'
    },
    log: {
        stream: './server/log/access.log'
    }
};
