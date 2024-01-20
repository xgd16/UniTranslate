<?php

class AuthEncrypt {
    private string $key;
    private array $params;

    public function __construct(string $key, array $params)
    {
        $this->key = $key;
        $this->params = $params;
    }

    public function encrypt(): string
    {
        return md5($this->key . $this->sortMapToStr($this->params));
    }

    private function isAssociativeArray(array $arr): bool {
        return array_keys($arr) !== range(0, count($arr) - 1);
    }

    private function sortMapToStr(array $params): string
    {
        $mapArr = [];

        foreach ($params as $key => $value) {
            if (is_array($value)) {
                if (!$this->isAssociativeArray($value)) {
                    $mapArr[] = "{$key}:" . implode(',', $value);
                } else {
                    $mapArr[] = "{$key}:|{$this->sortMapToStr($value)}|";
                }
                continue;
            }
            $mapArr[] = "{$key}:" . $value;
        }

        asort($mapArr);
        return implode('&', $mapArr);
    }
}

$a = new AuthEncrypt('123456', [
    'c' => [
        'cc' => 1,
        'cb' => 2,
        'ca' => 3,
        'cd' => 4,
    ],
    'a' => 1,
    'b' => [4, 1, 2],
]);

var_dump($a->encrypt());