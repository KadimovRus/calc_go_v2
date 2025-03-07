# Распределённый вычислитель арифметических выражений

Сервис для выполнения сложных арифметических вычислений в распределённой среде с использованием центрального координирующего узла и рабочих агентов.

## Общая схема работы

```mermaid
graph TD
    subgraph Клиент["Клиентская часть"]
        Client["Клиентское приложение"]
    end

    subgraph Оркестратор["Оркестратор (порт 8080)"]
        API["REST API"]
        Parser["Парсер выражений"]
        TaskManager["Менеджер задач"]
        ResultStore["Хранилище результатов"]
    end

    subgraph Агенты["Рабочие агенты"]
        Agent1["Агент 1"]
        Agent2["Агент 2"]
        AgentN["Агент N"]
    end

    Client -->|"POST /api/v1/calculate"| API
    API -->|"Выражение"| Parser
    Parser -->|"Операции"| TaskManager
    TaskManager -->|"Задания"| Agent1
    TaskManager -->|"Задания"| Agent2
    TaskManager -->|"Задания"| AgentN
    Agent1 -->|"Результаты"| TaskManager
    Agent2 -->|"Результаты"| TaskManager
    AgentN -->|"Результаты"| TaskManager
    TaskManager -->|"Объединение"| ResultStore
    ResultStore -->|"GET /api/v1/expressions"| API
    API -->|"Результат"| Client

    classDef client fill:#f9f,stroke:#333,stroke-width:2px,color:#000
    classDef orchestrator fill:#bbf,stroke:#333,stroke-width:2px,color:#000
    classDef agent fill:#bfb,stroke:#333,stroke-width:2px,color:#000
    
    class Client client
    class API,Parser,TaskManager,ResultStore orchestrator
    class Agent1,Agent2,AgentN agent
```

### Компоненты системы

- **Оркестратор** (порт 8080 по умолчанию):
   - Принимает запросы на вычисления через API
   - Разбивает выражения на составные операции
   - Управляет распределением задач
   - Собирает и хранит результаты

- **Агенты**:
   - Получают задания по HTTP
   - Производят вычисления с регулируемой задержкой
   - Передают ответы обратно в систему

## Технические требования

- Go 1.20+
- Поддерживаемые операции: `+`, `-`, `*`, `/`
- Обработка приоритетов и скобок
- Многопоточное выполнение

## Развёртывание

### Запуск оркестратора

```bash
export TIME_ADDITION_MS=200
export TIME_SUBTRACTION_MS=200
export TIME_MULTIPLICATIONS_MS=300
export TIME_DIVISIONS_MS=400

go run ./cmd/coordinator/main.go
```

### Запуск агента

```bash
export COMPUTING_POWER=4
export ORCHESTRATOR_URL=http://localhost:8080

go run ./cmd/agent/main.go
```

## REST API

### Отправка выражения на вычисление

```bash
POST /api/v1/calculate
```

#### Пример запроса:
```json
{
  "expression": "(3+5)*2-4/2"
}
```

#### Пример ответа:
```json
{
    "id": "42"
}
```

### Получение списка выражений
```bash
GET /api/v1/expressions
```

Пример ответа:
```json
{
    "expressions": [
        {
            "id": "42",
            "expression": "(3+5)*2-4/2",
            "status": "completed",
            "result": 14
        }
    ]
}
```

### Получение данных по конкретному вычислению
```bash
GET /api/v1/expressions/{id}
```

Пример ответа:
```json
{
    "expression": {
        "id": "42",
        "status": "completed",
        "result": 14
    }
}
```

## Внутренний API (для рабочих агентов)

### Запрос на выполнение задания
```bash
GET /internal/task
```

Пример ответа:
```json
{
    "task": {
        "id": "7",
        "arg1": 4,
        "arg2": 2,
        "operation": "/",
        "operation_time": 400
    }
}
```

### Отправка результата
```bash
POST /internal/task
```

Пример запроса:
```json
{
  "id": "7",
  "result": 2
}
```

## Конфигурация через переменные окружения

### Оркестратор
- `PORT` — порт сервиса (по умолчанию 8080)
- `TIME_ADDITION_MS` — задержка сложения (мс)
- `TIME_SUBTRACTION_MS` — задержка вычитания (мс)
- `TIME_MULTIPLICATIONS_MS` — задержка умножения (мс)
- `TIME_DIVISIONS_MS` — задержка деления (мс)

### Рабочие узлы
- `ORCHESTRATOR_URL` — адрес координатора
- `COMPUTING_POWER` — количество параллельных потоков

## Тестирование

```bash
go test ./...
```

