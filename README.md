# Happy — сервис персональных поздравлений (Nuxt 3 + Go/Gin + MongoDB)

**Идея:** пользователь создаёт поздравление по тематике (Новый год / ДР / … / Кастом), добавляет фото-коллаж и подарок.
Публикация доступна в двух форматах:

1) **Поздравление по коду** — открывается по ссылке `/c/CODE` (дешевле).  
2) **Поздравление по поддомену** — `имя.<BASE_DOMAIN>` (премиум, DNS через GoDaddy API).

Публикация активна **7 дней после оплаты**, затем страница, DNS-запись и файлы удаляются worker-ом.

---

## Стек

- **Frontend**: Nuxt 3 (SSR) + TailwindCSS + Pinia
- **Backend**: Go + Gin + JWT (cookie)
- **DB**: MongoDB
- **Infra**: docker-compose + внутренний Nginx (порт 8088 наружу)

---

## Быстрый старт (локально / сервер)

```bash
cp .env.example .env
# отредактируй .env (JWT_SECRET, BASE_DOMAIN, ADMIN_PASSWORD, ...)
docker compose up -d --build
```

После запуска:
- сайт: `http://SERVER:8088`
- API health: `http://SERVER:8088/api/health` (через nginx)
- админка: `/admin`  
  Логин/пароль берутся из `.env` (`ADMIN_USERNAME`, `ADMIN_PASSWORD`)

> **Важно:** Mongo и API порты наружу НЕ открываются (внутренняя сеть docker).

---

## Внешний Nginx на машине (проксирование)

Контейнерный nginx слушает **8088** на хосте. Снаружи вы можете проксировать туда из своего nginx:

```nginx
server {
  listen 443 ssl;
  server_name happy.example.com *.happy.example.com;

  location / {
    proxy_pass http://127.0.0.1:8088;
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-Proto https;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  }
}
```

---

## Как работает поддомен

### 1) В Nuxt включён “subdomain middleware”
Он определяет поддомен из `Host` и превращает запрос в `/s/<sub>`.

То есть запрос на `love.happy.example.com` будет обслужен страницей `pages/s/[sub].vue`,
которая загрузит поздравление через `GET /api/public/subdomain/love`.

### 2) DNS создаётся через GoDaddy API
При публикации поздравления типа `subdomain` backend вызовет GoDaddy:
`PUT /v1/domains/{GODADDY_DOMAIN}/records/{A|CNAME}/{sub}`

> Если `GODADDY_ENABLED=false`, публикация всё равно возможна — DNS операции просто пропускаются.

---

## Оплата (как подключить платёжный шлюз)

Сейчас по умолчанию включён режим **stub** (для разработки):
- при “Оплатить и опубликовать” создаётся Order со статусом `pending`
- пользователь редиректится на `/dashboard/orders/:id`
- в админке можно вручную нажать **“Пометить оплачено”**
- backend активирует поздравление на 7 дней

### Где интегрировать реальный провайдер
Код провайдера вынесен в `api/internal/payments`.

Интерфейс:
- `CreateCheckout(order)` → возвращает `checkoutURL` и `providerRef`
- `HandleWebhook(...)` → проверяет подпись/статус, возвращает `providerRef` и `paid=true`

Что нужно сделать для YooKassa / CloudPayments / Stripe:
1) Создать новый тип, который реализует `payments.Provider`
2) В `payments.FromConfig()` выбрать его по `PAYMENT_PROVIDER`
3) В `handlers.PaymentWebhook`:
   - получить `providerRef`
   - найти Order по `providerRef`
   - если `paid=true`: пометить Order как paid и активировать Greeting на 7 дней

### Схема связки (упрощённо)
1) Front: `POST /api/greetings/:id/publish`
2) Back:
   - создаёт Order
   - вызывает Provider.CreateCheckout → получает URL оплаты
   - возвращает URL фронту
3) Провайдер вызывает webhook `/api/webhooks/payment`
4) Back подтверждает оплату и активирует поздравление

---

## Модерация / предупреждение

На фронте показывается баннер:
“Запрещённый контент удаляется. Возврат средств не производится. Работаем по законодательству РФ.”

---

## Структура репозитория

- `api/` — Go API + worker
- `web/` — Nuxt 3 SSR
- `nginx/` — конфиг внутреннего nginx
- `docker-compose.yml` — весь стек
- `docker-compose.prod.yml` — override для образов из registry (GitLab CI)

---

## GitLab CI/CD

Файл `.gitlab-ci.yml`:
- build: собирает и пушит образы `api` и `web` в GitLab Container Registry
- deploy: rsync compose-конфигов на сервер и `docker compose pull && up -d`

Нужные CI variables:
- `SSH_PRIVATE_KEY`
- `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_PATH`
- и стандартные `CI_REGISTRY_USER/CI_REGISTRY_PASSWORD`

---

## Полезные эндпоинты

- Auth:
  - `POST /api/auth/register`
  - `POST /api/auth/login`
  - `GET  /api/auth/me`
  - `POST /api/auth/logout`
  - `POST /api/auth/recover/reset`
  - `POST /api/auth/password/change`

- Greetings:
  - `POST /api/greetings`
  - `GET  /api/greetings`
  - `GET  /api/greetings/:id`
  - `PUT  /api/greetings/:id`
  - `POST /api/greetings/:id/photos`
  - `POST /api/greetings/:id/publish`

- Public:
  - `GET /api/public/code/:code`
  - `GET /api/public/subdomain/:sub`

- Admin:
  - `GET /api/admin/users`
  - `DELETE /api/admin/users/:id`
  - `GET /api/admin/orders`
  - `POST /api/admin/orders/:id/mark-paid`
  - `GET /api/admin/greetings`
  - `DELETE /api/admin/greetings/:id`
  - `GET /api/admin/stats`
