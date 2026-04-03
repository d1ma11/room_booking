import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend } from 'k6/metrics';

// Метрика для отслеживания 95-го перцентиля
const slotListDuration = new Trend('slot_list_duration', true);

export const options = {
    // Имитация реального использования
    stages: [
        { duration: '30s', target: 10 },   // Ramp-up: 10 пользователей за 30 сек
        { duration: '1m', target: 10 },    // Стабильная нагрузка: 10 пользователей 1 мин
        { duration: '30s', target: 50 },   // Пик: 50 пользователей за 30 сек
        { duration: '1m', target: 50 },    // Стабильный пик: 50 пользователей 1 мин
        { duration: '30s', target: 0 },    // Ramp-down
    ],

    // Пороговые значения (SLI/SLA)
    thresholds: {
        http_req_duration: [
            'p(95)<200',  // 95% запросов должны укладываться в 200 мс
            'p(99)<500',  // 99% запросов — в 500 мс
        ],
        http_req_failed: ['rate<0.01'],    // Менее 1% ошибок
        slot_list_duration: ['p(95)<200'], // Кастомная метрика для /slots/list
    },

    // Дополнительные настройки
    summaryTrendStats: ['avg', 'min', 'med', 'p(90)', 'p(95)', 'p(99)'],
};

// Глобальные переменные (заполняются в setup)
let authToken;
let roomId;
let testDate;

// Setup: подготовка данных перед тестом
export function setup() {
    const baseUrl = __ENV.BASE_URL || 'http://localhost:8080';

    // 1. Получаем токен админа (для создания комнаты и расписания)
    const loginRes = http.post(`${baseUrl}/dummyLogin`,
        JSON.stringify({ role: 'admin' }),
        { headers: { 'Content-Type': 'application/json' } }
    );
    if (loginRes.status !== 200) {
        throw new Error('Failed to get admin token');
    }
    authToken = JSON.parse(loginRes.body).token;

    // 2. Создаём тестовую переговорку
    const roomRes = http.post(`${baseUrl}/rooms/create`,
        JSON.stringify({ name: 'Load Test Room', capacity: 10 }),
        { headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            }
        }
    );
    if (roomRes.status !== 201) {
        throw new Error('Failed to create room');
    }
    roomId = JSON.parse(roomRes.body).room.id;

    // 3. Создаём расписание (Пн-Пт, 09:00-18:00)
    const scheduleRes = http.post(`${baseUrl}/rooms/${roomId}/schedule/create`,
        JSON.stringify({
            daysOfWeek: [1,2,3,4,5],
            startTime: '09:00',
            endTime: '18:00'
        }),
        { headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            }
        }
    );
    // Игнорируем 409 (расписание уже есть) — это ок для повторных запусков
    if (scheduleRes.status !== 201 && scheduleRes.status !== 409) {
        throw new Error('Failed to create schedule');
    }

    // 4. Определяем тестовую дату (завтра)
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    testDate = tomorrow.toISOString().split('T')[0]; // YYYY-MM-DD

    return { authToken, roomId, testDate, baseUrl };
}

// Основной сценарий теста
export default function (data) {
    const { authToken, roomId, testDate, baseUrl } = data;

    // Запрос к самому нагруженному эндпоинту
    const res = http.get(`${baseUrl}/rooms/${roomId}/slots/list?date=${testDate}`, {
        headers: {
            'Authorization': `Bearer ${authToken}`,
            'Content-Type': 'application/json',
        },
        tags: { name: 'GET /rooms/{id}/slots/list' },
    });

    // Записываем время ответа в кастомную метрику
    slotListDuration.add(res.timings.duration);

    // Проверки (валидация ответа)
    check(res, {
        'status is 200': (r) => r.status === 200,
        'response has slots array': (r) => {
            const body = r.json();
            return body && Array.isArray(body.slots);
        },
        'slots have valid structure': (r) => {
            const body = r.json();
            if (!body.slots || body.slots.length === 0) return true; // пустой список — ок
            const slot = body.slots[0];
            return slot.id && slot.roomId && slot.start && slot.end;
        },
    });

    // Небольшая пауза между запросами (имитация пользователя)
    sleep(0.1 + Math.random() * 0.2); // 100-300 мс
}

// Teardown: очистка после теста (опционально)
export function teardown(data) {
    // Можно добавить удаление тестовых данных, если нужно
    // console.log('Test completed');
}