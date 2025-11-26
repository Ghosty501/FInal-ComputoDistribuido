import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

export let options = {
    vus: 50,
    duration: '3m',

    thresholds: {
        http_req_duration: ['p(95)<3000'],
        http_req_failed: ['rate<0.05'], // menos del 5% falla permitido
    },
};

const BASE = 'http://192.168.49.2:30536';

// Cargamos una lista de nombres para simular usuarios reales
const names = new SharedArray('names', () => [
  "Camila",
  "Armando",
  "Sofia",
  "Daniel",
  "Miguel",
  "Fernanda",
  "Karina",
  "Oscar",
  "Luis",
  "Beatriz",
]);

export default function () {

    // 1) Crear un contacto aleatorio
    let name = names[Math.floor(Math.random() * names.length)];
    let email = `${name.toLowerCase()}_${Math.random().toString(36).substring(7)}@example.com`;
    let phone = "555-" + Math.floor(100 + Math.random() * 900);

    let createContactRes = http.post(`${BASE}/contacts/create`, JSON.stringify({
        name,
        email,
        phone
    }), {
        headers: { 'Content-Type': 'application/json' }
    });

    check(createContactRes, {
        'createContact status is 200': (r) => r.status === 200
    });

    let contactId = createContactRes.json('id'); // obtener ID de contacto creado

    // 2) Crear un deal asociado a ese contacto
    let createDealRes = http.post(`${BASE}/deals/create`, JSON.stringify({
        contact_id: contactId,
        title: "Test Deal " + Math.random(),
        amount: Math.floor(100 + Math.random() * 900),
        currency: "USD",
        status: "open",
    }), {
        headers: { 'Content-Type': 'application/json' }
    });

    check(createDealRes, {
        'createDeal status is 200': (r) => r.status === 200
    });

    // 3) Leer contactos
    let getContacts = http.get(`${BASE}/contacts/list`);
    check(getContacts, {
        'listContacts status is 200': (r) => r.status === 200
    });

    // 4) Leer deals
    let getDeals = http.get(`${BASE}/deals/list`);
    check(getDeals, {
        'listDeals status is 200': (r) => r.status === 200
    });

    sleep(0.5); // pequeña pausa para evitar patrón robotizado
}
