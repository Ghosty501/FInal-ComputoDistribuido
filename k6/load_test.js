import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: 50,           // empieza con 50 VUs, puedes aumentar
  duration: '3m',    // 3 minutos
  thresholds: {
    http_req_duration: ['p(95)<3000'], // 95% solicitudes < 3s
  },
};

const GATEWAY_HOST = __ENV.GATEWAY_HOST || 'http://<GATEWAY_EXTERNAL_IP>'; // reemplaza o exporta env
const CREATE_DEAL_URL = `${GATEWAY_HOST}/deals/create`;

export default function () {
  const payload = JSON.stringify({
    contact_id: "1",
    title: "Prueba k6",
    amount: 1000.0,
    currency: "USD",
    status: "open"
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
    timeout: '60s',
  };

  const res = http.post(CREATE_DEAL_URL, payload, params);
  check(res, {
    'status 200': (r) => r.status === 200,
  });

  sleep(0.5);
}
