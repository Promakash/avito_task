# **Авито shop**

## **Магазин мерча**

В Авито существует внутренний магазин мерча, где сотрудники могут приобретать товары за монеты (coin). Каждому новому сотруднику выделяется 1000 монет, которые можно использовать для покупки товаров. Кроме того, монеты можно передавать другим сотрудникам в знак благодарности или как подарок.

## **Описание задачи**

Необходимо реализовать сервис, который позволит сотрудникам обмениваться монетками и приобретать на них мерч. Каждый сотрудник должен иметь возможность видеть:

- Список купленных им мерчовых товаров  
- Сгруппированную информацию о перемещении монеток в его кошельке, включая:  
  - Кто ему передавал монетки и в каком количестве  
  - Кому сотрудник передавал монетки и в каком количестве  

Количество монеток не может быть отрицательным, запрещено уходить в минус при операциях с монетками.

## **Общие вводные**

Мерч — это продукт, который можно купить за монетки. Всего в магазине доступно 10 видов мерча. Каждый товар имеет уникальное название и цену. Ниже приведён список наименований и их цены.

| Название     | Цена |
|--------------|------|
| t-shirt      | 80   |
| cup          | 20   |
| book         | 50   |
| pen          | 10   |
| powerbank    | 200  |
| hoody        | 300  |
| umbrella     | 200  |
| socks        | 10   |
| wallet       | 50   |
| pink-hoody   | 500  |

Предполагается, что в магазине бесконечный запас каждого вида мерча.

## **Условия**

- Используйте этот [API](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/schema.json).
- Сотрудников может быть до 100k, RPS — 1k, SLI времени ответа — 50 мс, SLI успешности ответа — 99.99%.  
- Для авторизации доступов должен использоваться JWT. Пользовательский токен доступа к API выдается после авторизации/регистрации пользователя. При первой авторизации пользователь должен создаваться автоматически.
- Реализуйте покрытие бизнес-сценариев юнит-тестами. Общее тестовое покрытие проекта должно превышать 40%.
- Реализуйте интеграционный или E2E-тест на сценарий покупки мерча.
- Реализуйте интеграционный или E2E-тест на сценарий передачи монеток другим сотрудникам.

## **Дополнительные задания**

Эти задания не являются обязательными, но выполнение всех или части из них даст вам преимущество перед другими кандидатами:

- Провести нагрузочное тестирование полученного решения и приложить результаты тестирования
- Реализовать интеграционное или E2E-тестирование для остальных сценариев
- Описать конфигурацию линтера (.golangci.yaml в корне проекта для Go, phpstan.neon для PHP или аналогичный файл, если вы используете другой ЯП)

---

## **Запуск проекта**

### **1️⃣ Клонирование репозитория**

```bash
git clone git@github.com:Promakash/avito_task.git
cd avito_task
```

### **2️⃣ Сборка и запуск проекта**

```bash
# Сборка Docker-изображений
make build
```

```bash
# Запуск всех сервисов (detached)
make run
```

```bash
# Запуск интеграционных тестов 
make run_tests
```

```bash
# Запуск юнит тестов 
make run_unit_tests
```

```bash
# Остановка всех сервисов
make stop
```

---

## **Конфигурация**

Конфигурация сервиса задается через YAML-файл или переменные окружения в `docker-compose.yaml`.

### **📌 HTTP-сервер**

| Параметр       | Значение | Описание               |
|----------------|----------|------------------------|
| address        | ":8080"  | Адрес сервера          |
| read_timeout   | 5s       | Таймаут чтения запроса |
| write_timeout  | 5s       | Таймаут записи ответа  |
| idle_timeout   | 30s      | Таймаут простоя        |

### **📌 JWT-секрет**

Задается исключительно через env `AUTH_SECRET`.

### **📌 PostgreSQL**

| Параметр | Значение  | Описание         |
|----------|-----------|------------------|
| host     | db        | Хост базы данных |
| port     | 5432      | Порт PostgreSQL  |
| user     | postgres  | Имя пользователя |
| password | password  | Пароль           |
| db_name  | shop      | Название БД      |

### **📌 Redis**

| Параметр       | Значение    | Описание                |
|----------------|-------------|-------------------------|
| host           | user-cache  | Хост Redis              |
| port           | 6379        | Порт Redis              |
| password       | redis       | Пароль Redis            |
| TTL            | 30m         | Время жизни кэша        |
| write_timeout  | 3s          | Таймаут записи в Redis  |
| read_timeout   | 2s          | Таймаут чтения из Redis |

### **📌 Логирование**

| Параметр  | Значение    | Описание                                                    |
|-----------|------------|-------------------------------------------------------------|
| level     | debug       | Уровень логов (debug, info, error) (default = info)        |
| format    | json        | Формат логов (json, text) (default = json)                 |
| directory | /app/logs   | Директория для логов в контейнере (default = "/app/logs")  |

---

## **Возникшие вопросы и уточнения**

Ниже приведён список вопросов по реализации, а также мой подход к их решению.

---

### **Вопрос 1**
В описании задачи и документации API сказано, что у пользователя в истории должна отображаться информация, кто ему передавал монетки и кому передавал он. Нужно ли в список транзакций добавлять покупки в шопе?

#### **Решение**
Логично отображать покупки пользователя в магазине, так как это также транзакции, демонстрирующие расход монеток

Пример `InfoResponse` для пользователя, который купил товар:

```json
{
  "coins": 800,
  "inventory": [
    {
      "type": "umbrella",
      "quantity": 1
    }
  ],
  "coinHistory": {
    "received": [],
    "sent": [
      {
        "toUser": "shop",
        "amount": 200
      }
    ]
  }
}
```

---

### **Вопрос 2**
Надо ли хранить переводы при удалении аккаунта пользователя?

#### **Решение**
Хоть сейчас и недоступна операция удаления аккаунта пользователя, но, как мне кажется, это необходимая вещь, так как хоть и удаляется отправитель/получатель, но факт транзакции имеется, поэтому при проектировании схемы БД, я решил оставлять переводы, даже при удалении аккаунта работника.

Схема бд при таком подходе:
```sql
    id         SERIAL PRIMARY KEY,
    sender     INT NULL,
    recipient  INT NULL,
    amount     INT NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    FOREIGN KEY (sender) REFERENCES employees (id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (recipient) REFERENCES employees (id) ON DELETE SET NULL ON UPDATE CASCADE
```

Соответственно при запросе на получение UserInfo получатель продолжит видеть всю историю транзакций, даже если запись второго участника транзакций уже удалена. Выглядеть это будет примерно вот так:
```json
{
  "coins": 800,
  "inventory": [],
  "coinHistory": {
    "received": [],
    "sent": [
      {
        "toUser": "deleted",
        "amount": 200
      }
    ]
  }
}
```
Из минусов стоило бы отметить, что при удалении и аккаунта получателя, и аккаунта отправителя, транзакция так же остается в базе данных. Если в будущем добавится возможность удаления аккаунтов, можно реализовать триггер AFTER DELETE на employees, который будет автоматически очищать транзакции, где receiver и sender стали NULL.

---

### **Вопрос 3**
В каком порядке должны отображаться транзакции?

#### **Решение**
В задании нет уточнений, как нужно возвращать транзакции, но мне показалось бы правильным, если бы они возвращались в отсортированном виде: сначала новые затем старые. Поэтому при запросе на получение транзакций происходит сортировка по времени транзакции. В начале массивов транзакций находятся более новые.

---

### **Вопрос 4**
Будет ли в реальном времени изменяться цена товаров или сам ассортимент?

#### **Решение**
В задании сказано, что всего будет 10 категорий товаров, поэтому я сделал допущение, что товары и их характеристики остаются статичными после запуска сервиса. Это позволило мне имплементировать кэширование через sync.Map (поднимать целый redis для 10 записей мне кажется это чересчур). Что позволило значительно уменьшить нагрузку на базу данных при покупке товара.

---

### **Вопрос 5**
Есть ли смысл использовать Redis?

#### **Решение**
В задании оговорена высокая нагрузка (1k RPS) и до 100k пользователей. Разумно закэшировать в Redis наиболее частые запросы — например, проверку имени пользователя при переводах. Так можно разгрузить базу данных, если мы постоянно обращаемся к ней для получения `id` по имени. При этом, чтобы не захламлять Redis, можно использовать TTL (например, 30 минут).
