import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    // 1. Smoke Test (Uncomment to run)
    /*
    vus: 1,
    duration: '1m',
    */

    // 2. Load Test (Ramping up to 20 users)
    stages: [
        { duration: '30s', target: 20 }, 
        { duration: '1m', target: 20 },  
        { duration: '30s', target: 0 },  
    ],

    thresholds: {
        http_req_failed: ['rate<0.01'], 
        http_req_duration: ['p(95)<500'], 
    },
};

const BASE_URL = 'https://api-portofolio.declarationdigital.tech';

export default function () {
    // --- SCENARIO 1: Public Health Check ---
    let healthRes = http.get(`${BASE_URL}/users/health`);
    check(healthRes, {
        'health check status is 200': (r) => r.status === 200,
    });

    // --- SCENARIO 2: Public Portfolio List ---
    let portfolioRes = http.get(`${BASE_URL}/portofolio/all`);
    check(portfolioRes, {
        'get all portfolios status is 200': (r) => r.status === 200,
    });

    // --- SCENARIO 3: Auth Flow (Register & Login) ---
    const username = `user_${randomString(8)}`;
    const email = `${username}@example.com`;
    const password = 'password123';

    // Register
    let regPayload = JSON.stringify({
        name: 'Load Test User',
        username: username,
        email: email,
        password: password,
    });
    let regParams = { headers: { 'Content-Type': 'application/json' } };
    
    let regRes = http.post(`${BASE_URL}/users/register`, regPayload, regParams);
    check(regRes, {
        'register status is 201': (r) => r.status === 201,
    });

    // Login
    let loginPayload = JSON.stringify({
        email: email,
        password: password,
    });
    
    let loginRes = http.post(`${BASE_URL}/users/login`, loginPayload, regParams);
    const loginOk = check(loginRes, {
        'login status is 200': (r) => r.status === 200,
        'has token': (r) => r.json().token !== undefined,
    });

    if (loginOk) {
        const token = loginRes.json().token;
        const authParams = {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        };

        // --- SCENARIO 4: Protected Endpoint (Get My Portfolio) ---
        let myPortfolioRes = http.get(`${BASE_URL}/portofolio/`, authParams);
        check(myPortfolioRes, {
            'get my portfolio status is 200': (r) => r.status === 200,
        });
    }

    sleep(1);
}
