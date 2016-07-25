module.exports = {
    timeout: '120s',
    useSSL: false,
    port: 8080,
    sslPort: 8443,
    proxyPort: 8009,
    ssl: {
        key:    './ssl/tls/private/ssl.key',
        cert:   './ssl/tls/certs/ssl.crt',
        ca:     [
            './ssl/tls/certs/ssl-bundle.crt'
        ]
    },
    api: {
        base_url: 'http://localhost:3000/api/',
        key: ''
    },
    files: {
        static: './app/dist/public/'
    },
    log: {
        stream: './server/log/access.log'
    },
    reject_unauthorized: 0 // 0 if using self-signed certs, 1 if trusted certs
};
