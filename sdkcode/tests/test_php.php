<?php
require_once '../php/UniTranslate.php';

$testData = json_decode(file_get_contents('test_data.json'), true);
$results = [];

foreach ($testData['testCases'] as $case) {
    $hash = UniTranslate::authEncrypt($testData['key'], $case['params']);
    $results[] = [
        'name' => $case['name'],
        'success' => $hash === $case['expectedHash'],
        'got' => $hash,
        'expected' => $case['expectedHash']
    ];
}

echo json_encode($results, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE);
