import json
import sys
sys.path.append('../python')
from uni_translate import UniTranslate

with open('test_data.json', 'r', encoding='utf-8') as f:
    test_data = json.load(f)

results = []
for case in test_data['testCases']:
    hash = UniTranslate.auth_encrypt(test_data['key'], case['params'])
    results.append({
        'name': case['name'],
        'success': hash == case['expectedHash'],
        'got': hash,
        'expected': case['expectedHash']
    })

print(json.dumps(results, indent=2, ensure_ascii=False))
