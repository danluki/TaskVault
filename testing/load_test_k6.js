import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '5s', target: 100 },
        { duration: '10s', target: 100 },
        { duration: '5s', target: 0 },
    ],
};

const BASE_URLS = [
    'http://localhost:8080/v1/storage',
    'http://localhost:8081/v1/storage',
];

function tryPost(key, value) {
    for (const baseUrl of BASE_URLS) {
        try {
            const res = http.post(baseUrl, JSON.stringify({ key, value }), {
                headers: { 'Content-Type': 'application/json' },
            });
            if (res.status === 201) {
                check(res, { 'POST status is 201': () => true });
                return { success: true, baseUrl };
            }
        } catch (_) {
        }
    }
    check(null, { 'POST status is 201': () => false });
    return { success: false };
}

function tryGet(baseUrl, key, expectedValue) {
    try {
        const res = http.get(`${baseUrl}/${key}`);
        check(res, {
            'GET status is 200': (r) => r.status === 200,
            'GET returns correct value': (r) => r.json() === expectedValue,
        });
    } catch (_) {
        check(null, { 'GET status is 200': () => false });
    }
}

export default function () {
    const key = `key_${Math.random().toString(36).substring(7)}`;
    const value = `value_${Math.random().toString(36).substring(7)}`;

    const postResult = tryPost(key, value);
    if (postResult.success) {
        tryGet(postResult.baseUrl, key, value);
    }

    sleep(1);
}