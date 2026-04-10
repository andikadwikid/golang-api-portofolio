import http from 'k6/http';
import { check } from 'k6';

export const options = {
    scenarios: {
        high_traffic: {
            executor: 'ramping-arrival-rate',
            startRate: 50, // mulai dari 50 request/sec
            timeUnit: '1s',
            preAllocatedVUs: 100,
            maxVUs: 1000,
            stages: [
                { target: 100, duration: '1m' },  // naik ke 100 rps
                { target: 200, duration: '2m' },  // naik ke 200 rps
                { target: 300, duration: '2m' },  // spike
                { target: 0, duration: '1m' },    // turun
            ],
        },
    },

    thresholds: {
        http_req_failed: ['rate<0.02'],       // toleransi sedikit error
        http_req_duration: ['p(95)<800'],     // lebih longgar karena high load
    },
};

const BASE_URL = 'https://api-portofolio.declarationdigital.tech';

// 🔑 Ambil token SEKALI saja (bukan tiap request)
export function setup() {
    const loginRes = http.post(`${BASE_URL}/users/login`, JSON.stringify({
        email: 'johndoe@example.com',
        password: 'password',
    }), {
        headers: { 'Content-Type': 'application/json' },
    });

    return {
        token: loginRes.json('token'),
    };
}

export default function (data) {
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${data.token}`,
    };

    // ⚡ Fokus ke endpoint penting saja (bukan register)
    const responses = http.batch([
        ['GET', `${BASE_URL}/users/health`],
        ['GET', `${BASE_URL}/portofolio/all`],
        ['GET', `${BASE_URL}/portofolio/`, null, { headers }],
    ]);

    check(responses[0], {
        'health OK': (r) => r.status === 200,
    });

    check(responses[1], {
        'portfolio list OK': (r) => r.status === 200,
    });

    check(responses[2], {
        'my portfolio OK': (r) => r.status === 200,
    });
}