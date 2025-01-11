const UniTranslate = require('../javascript/UniTranslate');
const fs = require('fs');

const testData = JSON.parse(fs.readFileSync('test_data.json', 'utf8'));
const results = [];

testData.testCases.forEach(testCase => {
    const hash = UniTranslate.authEncrypt(testData.key, testCase.params);
    results.push({
        name: testCase.name,
        success: hash === testCase.expectedHash,
        got: hash,
        expected: testCase.expectedHash
    });
});

console.log(JSON.stringify(results, null, 2));
