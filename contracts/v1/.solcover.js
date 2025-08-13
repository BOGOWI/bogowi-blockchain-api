module.exports = {
    skipFiles: [
        'mocks/',
        'test/',
        'interfaces/'
    ],
    mocha: {
        grep: "@skip-on-coverage",
        invert: true
    }
};