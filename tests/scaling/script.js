import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  discardResponseBodies: true,
  scenarios: {
    contacts: {
      executor: 'ramping-vus',
      startVUs: 8,
      stages: [
        { duration: '30s', target: 32 },
        { duration: '30s', target: 16 },
        { duration: '30s', target: 24 },
        { duration: '30s', target: 50 },
      ],
      gracefulRampDown: '0s',
    },
  },
};

export default function() {
  let res = http.get('http://localhost:8080/api/screenshot?url=https://google.com');
  check(res, { "status is 200": (res) => res.status === 200 });
  sleep(1);
}
