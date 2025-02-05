import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 100 },
        { duration: '1m', target: 100 },
        { duration: '30s', target: 0 },
    ],
};

const BASE_URL = 'http://localhost:8080/v1/storage';

export default function () {
    let key = `key_${Math.random().toString(36).substring(7)}`;
    let value = `value_${Math.random().toString(36).substring(7)}`;

    let postRes = http.post(BASE_URL, JSON.stringify({ key: key, value: value }), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(postRes, {
        'POST status is 201': (r) => r.status === 201,
    });

    let getRes = http.get(`${BASE_URL}/${key}`);
    check(getRes, {
        'GET status is 200': (r) => r.status === 200,
        'GET returns correct value': (r) => r.json() === value,
    });

    sleep(1);
}
