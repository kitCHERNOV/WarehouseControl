# WarehouseControl — Архитектурные решения

## Обзор задачи

Система управления складом должна обеспечивать:
- CRUD-операции для товаров
- Историю всех изменений (кто, когда, что изменил)
- Ролевую модель доступа (admin, manager, viewer)
- JWT-авторизацию
- Веб-интерфейс для взаимодействия

**Важное примечание:** История изменений должна реализовываться через триггеры в PostgreSQL (антипаттерн для учебных целей).

---

## Вариант 1: Классическая монолитная архитектура с триггерами БД

### Описание архитектуры

Этот вариант представляет собой традиционный монолитный подход, где вся бизнес-логика сосредоточена в одном приложении, а история изменений полностью делегируется триггерам базы данных.

### Структура системы

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend Layer                           │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Web UI (HTML/CSS/JavaScript)                             │  │
│  │  - Login page с выбором роли                              │  │
│  │  - Dashboard со списком товаров                           │  │
│  │  - Forms для CRUD операций                                │  │
│  │  - History view для каждого товара                        │  │
│  └───────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTP/REST API
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Backend Layer                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  HTTP Server (Express.js / FastAPI / Gin)                 │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  Auth Middleware                                    │  │  │
│  │  │  - JWT verification                                 │  │  │
│  │  │  - Role extraction                                  │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  Routes / Controllers                              │  │  │
│  │  │  - POST   /auth/login                              │  │  │
│  │  │  - GET    /items                                   │  │  │
│  │  │  - POST   /items                                   │  │  │
│  │  │  - PUT    /items/{id}                              │  │  │
│  │  │  - DELETE /items/{id}                              │  │  │
│  │  │  - GET    /items/{id}/history                      │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  Business Logic Layer                               │  │  │
│  │  │  - Role-based access control (RBAC)                  │  │  │
│  │  │  - Input validation                                  │  │  │
│  │  │  - Error handling                                    │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  Data Access Layer                                   │  │  │
│  │  │  - SQL queries / ORM                                 │  │  │
│  │  │  - Connection pool management                        │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             │ PostgreSQL Protocol
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Database Layer (PostgreSQL)                  │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  Tables:                                               │    │
│  │  - users (id, username, role, password_hash)          │    │
│  │  - items (id, name, quantity, price, description)     │    │
│  │  - item_history (id, item_id, action, old_value,      │    │
│  │                 new_value, changed_by, changed_at)     │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  Triggers (АНТИПАТТЕРН):                                 │    │
│  │  - trigger_items_insert() → AFTER INSERT ON items       │    │
│  │  - trigger_items_update() → AFTER UPDATE ON items       │    │
│  │  - trigger_items_delete() → AFTER DELETE ON items       │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  Functions:                                             │    │
│  │  - log_item_change() → записывает в item_history       │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

### Детальное описание компонентов

#### 1. Frontend Layer

**Технологии:** HTML5, CSS3, Vanilla JavaScript или React/Vue

**Компоненты:**
- **Login Page:** Выпадающий список для выбора роли (admin/manager/viewer), кнопка входа
- **Dashboard:** Таблица со списком товаров, кнопки действий (согласно роли)
- **Item Form:** Модальное окно для создания/редактирования товара
- **History View:** Таблица с историей изменений конкретного товара

**Функциональность:**
- Хранение JWT токена в localStorage
- Автоматическое добавление Authorization header к запросам
- Обработка ошибок авторизации (redirect на login)
- Динамическое отображение кнопок согласно роли пользователя

#### 2. Backend Layer

**Технологии:** Node.js + Express.js / Python + FastAPI / Go + Gin

**Auth Middleware:**
```javascript
// Псевдокод middleware
function authMiddleware(req, res, next) {
    const token = req.headers.authorization?.split(' ')[1];
    const decoded = jwt.verify(token, SECRET_KEY);
    req.user = { id: decoded.userId, role: decoded.role };
    next();
}
```

**Routes с RBAC:**
```javascript
// Пример маршрутизации с проверкой ролей
app.get('/items', authMiddleware, (req, res) => {
    // viewer, manager, admin - все могут смотреть
});

app.post('/items', authMiddleware, requireRole(['admin', 'manager']), (req, res) => {
    // только admin и manager могут создавать
});

app.put('/items/:id', authMiddleware, requireRole(['admin', 'manager']), (req, res) => {
    // только admin и manager могут редактировать
});

app.delete('/items/:id', authMiddleware, requireRole(['admin']), (req, res) => {
    // только admin может удалять
});
```

**Business Logic:**
- Валидация входных данных (name, quantity, price)
- Проверка существования товара перед обновлением/удалением
- Формирование ответов с соответствующими HTTP статусами

#### 3. Database Layer

**Таблицы:**

```sql
-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'manager', 'viewer')),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица товаров
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    price DECIMAL(10, 2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица истории изменений
CREATE TABLE item_history (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    action VARCHAR(20) NOT NULL, -- INSERT, UPDATE, DELETE
    old_data JSONB,
    new_data JSONB,
    changed_by VARCHAR(50) NOT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Триггеры (АНТИПАТТЕРН):**

```sql
-- Функция логирования изменений
CREATE OR REPLACE FUNCTION log_item_change()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO item_history (item_id, action, new_data, changed_by)
        VALUES (NEW.id, 'INSERT', row_to_json(NEW), current_user);
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO item_history (item_id, action, old_data, new_data, changed_by)
        VALUES (NEW.id, 'UPDATE', row_to_json(OLD), row_to_json(NEW), current_user);
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO item_history (item_id, action, old_data, changed_by)
        VALUES (OLD.id, 'DELETE', row_to_json(OLD), current_user);
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Создание триггеров
CREATE TRIGGER trigger_items_insert
    AFTER INSERT ON items
    FOR EACH ROW EXECUTE FUNCTION log_item_change();

CREATE TRIGGER trigger_items_update
    AFTER UPDATE ON items
    FOR EACH ROW EXECUTE FUNCTION log_item_change();

CREATE TRIGGER trigger_items_delete
    AFTER DELETE ON items
    FOR EACH ROW EXECUTE FUNCTION log_item_change();
```

### Потоки данных

**1. Аутентификация:**
```
Frontend → POST /auth/login → Backend → JWT generation → Frontend (stores token)
```

**2. Создание товара:**
```
Frontend → POST /items (with JWT) → Auth Middleware → Controller → 
DB (INSERT items) → Trigger fires → INSERT item_history → Response
```

**3. Просмотр истории:**
```
Frontend → GET /items/{id}/history (with JWT) → Auth Middleware → 
Controller → DB (SELECT * FROM item_history WHERE item_id = ?) → Response
```

### Преимущества и недостатки

**Преимущества:**
- Простота реализации и понимания
- Минимальное количество компонентов
- Быстрый старт разработки
- Легкая отладка

**Недостатки:**
- Бизнес-логика истории скрыта в БД (триггеры)
- Сложность тестирования (нужна настоящая БД)
- Проблемы с миграциями (триггеры могут ломаться)
- Трудно отслеживать изменения в логике истории
- Зависимость от конкретной СУБД
- Триггеры выполняются синхронно, могут замедлять операции
- Невозможность добавить дополнительную логику в процесс истории без модификации БД

---

## Вариант 2: Слоистая архитектура с паттерном Repository

### Описание архитектуры

Этот вариант использует более структурированный подход с четким разделением ответственности между слоями. История изменений всё ещё реализуется через триггеры (по требованию), но приложение имеет более чистую архитектуру с использованием паттерна Repository.

### Структура системы

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend Layer                           │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Single Page Application (React / Vue)                    │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  Components:                                        │  │  │
│  │  │  - LoginForm (role selector)                         │  │  │
│  │  │  - ItemsList (table with filters)                   │  │  │
│  │  │  - ItemForm (create/edit modal)                     │  │  │
│  │  │  - HistoryView (timeline/table)                     │  │  │
│  │  │  - PermissionWrapper (hides/shows UI by role)       │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  State Management (Redux / Pinia)                  │  │  │
│  │  │  - auth state (token, user role)                    │  │  │
│  │  │  - items state (list, loading, error)              │  │  │
│  │  │  - history state (by item_id)                       │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  API Client (Axios with interceptors)               │  │  │
│  │  │  - Automatic token injection                         │  │  │
│  │  │  - Request/response logging                          │  │  │
│  │  │  - Error handling (401 → redirect)                 │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTP/REST API
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Backend Layer                              │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Presentation Layer (Controllers/Handlers)                 │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  AuthController                                    │  │  │
│  │  │  - login(dto: LoginDto) → TokenResponse            │  │  │
│  │  │  - refreshToken(dto: RefreshTokenDto) → Token       │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  ItemsController                                    │  │  │
│  │  │  - findAll() → ItemsResponse                        │  │  │
│  │  │  - findById(id: number) → ItemResponse              │  │  │
│  │  │  - create(dto: CreateItemDto) → ItemResponse       │  │  │
│  │  │  - update(id: number, dto: UpdateItemDto) → Item   │  │  │
│  │  │  - delete(id: number) → void                        │  │  │
│  │  │  - getHistory(id: number) → HistoryResponse         │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────┘  │
│                              │                                  │
│                              ▼                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Application Layer (Services)                             │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  AuthService                                        │  │  │
│  │  │  - validateCredentials(username, role)              │  │  │
│  │  │  - generateToken(user) → JWT                        │  │  │
│  │  │  - verifyToken(token) → UserPayload                │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  ItemsService                                       │  │  │
│  │  │  - getAllItems() → Item[]                           │  │  │
│  │  │  - getItemById(id) → Item                           │  │  │
│  │  │  - createItem(dto, user) → Item                     │  │  │
│  │  │  - updateItem(id, dto, user) → Item                 │  │  │
│  │  │  - deleteItem(id, user) → void                     │  │  │
│  │  │  - getItemHistory(id) → HistoryEntry[]              │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  AuthorizationService                                │  │  │
│  │  │  - canCreate(user) → boolean                        │  │  │
│  │  │  - canUpdate(user) → boolean                        │  │  │
│  │  │  - canDelete(user) → boolean                        │  │  │
│  │  │  - canView(user) → boolean                          │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────┘  │
│                              │                                  │
│                              ▼                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Domain Layer (Entities, DTOs, Validators)               │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  Entities:                                           │  │  │
│  │  │  - User (id, username, role)                        │  │  │
│  │  │  - Item (id, name, quantity, price, description)     │  │  │
│  │  │  - HistoryEntry (id, itemId, action, oldData, ...)   │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  DTOs:                                               │  │  │
│  │  │  - LoginDto, CreateItemDto, UpdateItemDto, ...      │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  Validators:                                         │  │  │
│  │  │  - ItemValidator (name length, quantity >= 0, ...)   │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────┘  │
│                              │                                  │
│                              ▼                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Infrastructure Layer (Repositories, Database)             │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  Repository Interface:                               │  │  │
│  │  │  interface IItemsRepository {                       │  │  │
│  │  │    findAll(): Promise<Item[]>                       │  │  │
│  │  │    findById(id): Promise<Item>                       │  │  │
│  │  │    create(item): Promise<Item>                       │  │  │
│  │  │    update(id, item): Promise<Item>                   │  │  │
│  │  │    delete(id): Promise<void>                         │  │  │
│  │  │    findHistory(itemId): Promise<HistoryEntry[]>      │  │  │
│  │  │  }                                                   │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  PostgreSQL Implementation:                          │  │  │
│  │  │  class PostgresItemsRepository implements IItemsRepository │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │                                                           │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │  Database Connection:                               │  │  │
│  │  │  - Connection pool (pg / Prisma / TypeORM)          │  │  │
│  │  │  - Migration scripts                                 │  │  │
│  │  │  - Seed data                                         │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             │ PostgreSQL Protocol
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Database Layer (PostgreSQL)                  │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  Tables:                                               │    │
│  │  - users                                               │    │
│  │  - items                                               │    │
│  │  - item_history                                        │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  Triggers (АНТИПАТТЕРН - обязательное требование):      │    │
│  │  - trigger_items_insert                                 │    │
│  │  - trigger_items_update                                 │    │
│  │  - trigger_items_delete                                 │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  Indexes:                                               │    │
│  │  - idx_item_history_item_id                             │    │
│  │  - idx_item_history_changed_at                          │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

### Детальное описание компонентов

#### 1. Frontend Layer (SPA)

**Технологии:** React + TypeScript + Redux Toolkit / Vue 3 + TypeScript + Pinia

**Структура директорий:**
```
src/
├── components/
│   ├── auth/
│   │   └── LoginForm.tsx
│   ├── items/
│   │   ├── ItemsList.tsx
│   │   ├── ItemForm.tsx
│   │   └── ItemActions.tsx
│   ├── history/
│   │   └── HistoryView.tsx
│   └── common/
│       ├── PermissionWrapper.tsx
│       └── ProtectedRoute.tsx
├── store/
│   ├── authSlice.ts
│   ├── itemsSlice.ts
│   └── historySlice.ts
├── services/
│   └── api.ts (Axios instance with interceptors)
├── types/
│   └── index.ts
└── App.tsx
```

**Компоненты:**

**LoginForm.tsx:**
```typescript
interface LoginFormProps {
  onLogin: (role: UserRole) => void;
}

// Выпадающий список с ролями: admin, manager, viewer
// При отправке вызывает API для получения JWT
```

**PermissionWrapper.tsx:**
```typescript
interface PermissionWrapperProps {
  allowedRoles: UserRole[];
  children: React.ReactNode;
}

// Скрывает/показывает children в зависимости от роли пользователя
```

**API Service (api.ts):**
```typescript
const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
});

// Request interceptor - добавляет токен
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - обрабатывает 401
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Redirect to login
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

#### 2. Backend Layer

**Технологии:** NestJS / Spring Boot / .NET Core

**Структура директорий (NestJS пример):**
```
src/
├── modules/
│   ├── auth/
│   │   ├── auth.controller.ts
│   │   ├── auth.service.ts
│   │   ├── auth.module.ts
│   │   ├── strategies/
│   │   │   └── jwt.strategy.ts
│   │   └── guards/
│   │       ├── jwt-auth.guard.ts
│   │       └── roles.guard.ts
│   ├── items/
│   │   ├── items.controller.ts
│   │   ├── items.service.ts
│   │   ├── items.module.ts
│   │   ├── dto/
│   │   │   ├── create-item.dto.ts
│   │   │   ├── update-item.dto.ts
│   │   │   └── query-items.dto.ts
│   │   └── entities/
│   │       └── item.entity.ts
│   └── history/
│       ├── history.controller.ts
│       ├── history.service.ts
│       └── history.module.ts
├── common/
│   ├── decorators/
│   │   └── roles.decorator.ts
│   ├── enums/
│   │   └── role.enum.ts
│   └── guards/
│       └── authorization.guard.ts
├── database/
│   ├── migrations/
│   ├── seeds/
│   └── triggers.sql
└── main.ts
```

**Presentation Layer - ItemsController:**
```typescript
@Controller('items')
@UseGuards(JwtAuthGuard)
export class ItemsController {
  constructor(private readonly itemsService: ItemsService) {}

  @Get()
  @Roles('admin', 'manager', 'viewer')
  @UseGuards(RolesGuard)
  findAll(@Query() query: QueryItemsDto) {
    return this.itemsService.findAll(query);
  }

  @Get(':id')
  @Roles('admin', 'manager', 'viewer')
  @UseGuards(RolesGuard)
  findOne(@Param('id') id: number) {
    return this.itemsService.findOne(id);
  }

  @Post()
  @Roles('admin', 'manager')
  @UseGuards(RolesGuard)
  create(
    @Body() createItemDto: CreateItemDto,
    @CurrentUser() user: UserPayload
  ) {
    return this.itemsService.create(createItemDto, user);
  }

  @Put(':id')
  @Roles('admin', 'manager')
  @UseGuards(RolesGuard)
  update(
    @Param('id') id: number,
    @Body() updateItemDto: UpdateItemDto,
    @CurrentUser() user: UserPayload
  ) {
    return this.itemsService.update(id, updateItemDto, user);
  }

  @Delete(':id')
  @Roles('admin')
  @UseGuards(RolesGuard)
  delete(
    @Param('id') id: number,
    @CurrentUser() user: UserPayload
  ) {
    return this.itemsService.delete(id, user);
  }

  @Get(':id/history')
  @Roles('admin', 'manager', 'viewer')
  @UseGuards(RolesGuard)
  getHistory(@Param('id') id: number) {
    return this.itemsService.getHistory(id);
  }
}
```

**Application Layer - ItemsService:**
```typescript
@Injectable()
export class ItemsService {
  constructor(
    private readonly itemsRepository: IItemsRepository,
    private readonly authorizationService: AuthorizationService
  ) {}

  async findAll(query: QueryItemsDto): Promise<PaginatedResponse<Item>> {
    return this.itemsRepository.findAll(query);
  }

  async findOne(id: number): Promise<Item> {
    const item = await this.itemsRepository.findById(id);
    if (!item) {
      throw new NotFoundException(`Item with id ${id} not found`);
    }
    return item;
  }

  async create(dto: CreateItemDto, user: UserPayload): Promise<Item> {
    if (!this.authorizationService.canCreate(user)) {
      throw new ForbiddenException('Insufficient permissions');
    }

    // Валидация бизнес-правил
    this.validateItemData(dto);

    // Создаём товар - триггер автоматически запишет историю
    const item = await this.itemsRepository.create({
      ...dto,
      createdBy: user.username
    });

    return item;
  }

  async update(
    id: number,
    dto: UpdateItemDto,
    user: UserPayload
  ): Promise<Item> {
    if (!this.authorizationService.canUpdate(user)) {
      throw new ForbiddenException('Insufficient permissions');
    }

    const existingItem = await this.findOne(id);
    this.validateItemData(dto);

    // Обновляем товар - триггер автоматически запишет историю
    const updatedItem = await this.itemsRepository.update(id, {
      ...dto,
      updatedBy: user.username
    });

    return updatedItem;
  }

  async delete(id: number, user: UserPayload): Promise<void> {
    if (!this.authorizationService.canDelete(user)) {
      throw new ForbiddenException('Insufficient permissions');
    }

    await this.findOne(id); // Проверка существования

    // Удаляем товар - триггер автоматически запишет историю
    await this.itemsRepository.delete(id);
  }

  async getHistory(id: number): Promise<HistoryEntry[]> {
    await this.findOne(id); // Проверка существования
    return this.itemsRepository.findHistory(id);
  }

  private validateItemData(dto: CreateItemDto | UpdateItemDto): void {
    if (dto.name && dto.name.length < 3) {
      throw new BadRequestException('Name must be at least 3 characters');
    }
    if (dto.quantity !== undefined && dto.quantity < 0) {
      throw new BadRequestException('Quantity cannot be negative');
    }
    if (dto.price !== undefined && dto.price < 0) {
      throw new BadRequestException('Price cannot be negative');
    }
  }
}
```

**AuthorizationService:**
```typescript
@Injectable()
export class AuthorizationService {
  canCreate(user: UserPayload): boolean {
    return ['admin', 'manager'].includes(user.role);
  }

  canUpdate(user: UserPayload): boolean {
    return ['admin', 'manager'].includes(user.role);
  }

  canDelete(user: UserPayload): boolean {
    return user.role === 'admin';
  }

  canView(user: UserPayload): boolean {
    return ['admin', 'manager', 'viewer'].includes(user.role);
  }
}
```

**Infrastructure Layer - Repository:**

**Interface:**
```typescript
export interface IItemsRepository {
  findAll(query: QueryItemsDto): Promise<PaginatedResponse<Item>>;
  findById(id: number): Promise<Item | null>;
  create(dto: CreateItemDto & { createdBy: string }): Promise<Item>;
  update(id: number, dto: UpdateItemDto & { updatedBy: string }): Promise<Item>;
  delete(id: number): Promise<void>;
  findHistory(itemId: number): Promise<HistoryEntry[]>;
}
```

**PostgreSQL Implementation:**
```typescript
@Injectable()
export class PostgresItemsRepository implements IItemsRepository {
  constructor(@InjectDataSource() private dataSource: DataSource) {}

  async findAll(query: QueryItemsDto): Promise<PaginatedResponse<Item>> {
    const { page = 1, limit = 10, search } = query;
    const offset = (page - 1) * limit;

    let sql = 'SELECT * FROM items';
    const params: any[] = [];

    if (search) {
      sql += ' WHERE name ILIKE $1';
      params.push(`%${search}%`);
    }

    sql += ' ORDER BY created_at DESC LIMIT $2 OFFSET $3';
    params.push(limit, offset);

    const result = await this.dataSource.query(sql, params);

    // Получаем общее количество
    let countSql = 'SELECT COUNT(*) FROM items';
    if (search) {
      countSql += ' WHERE name ILIKE $1';
    }
    const countResult = await this.dataSource.query(countSql, search ? [params[0]] : []);
    const total = parseInt(countResult[0].count);

    return {
      data: result,
      meta: {
        page,
        limit,
        total,
        totalPages: Math.ceil(total / limit)
      }
    };
  }

  async findById(id: number): Promise<Item | null> {
    const result = await this.dataSource.query(
      'SELECT * FROM items WHERE id = $1',
      [id]
    );
    return result[0] || null;
  }

  async create(dto: CreateItemDto & { createdBy: string }): Promise<Item> {
    const result = await this.dataSource.query(
      `INSERT INTO items (name, quantity, price, description, created_by)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING *`,
      [dto.name, dto.quantity, dto.price, dto.description, dto.createdBy]
    );
    // Триггер автоматически запишет историю
    return result[0];
  }

  async update(
    id: number,
    dto: UpdateItemDto & { updatedBy: string }
  ): Promise<Item> {
    const result = await this.dataSource.query(
      `UPDATE items
       SET name = COALESCE($1, name),
           quantity = COALESCE($2, quantity),
           price = COALESCE($3, price),
           description = COALESCE($4, description),
           updated_by = $5,
           updated_at = CURRENT_TIMESTAMP
       WHERE id = $6
       RETURNING *`,
      [dto.name, dto.quantity, dto.price, dto.description, dto.updatedBy, id]
    );
    // Триггер автоматически запишет историю
    return result[0];
  }

  async delete(id: number): Promise<void> {
    await this.dataSource.query('DELETE FROM items WHERE id = $1', [id]);
    // Триггер автоматически запишет историю
  }

  async findHistory(itemId: number): Promise<HistoryEntry[]> {
    const result = await this.dataSource.query(
      `SELECT * FROM item_history
       WHERE item_id = $1
       ORDER BY changed_at DESC`,
      [itemId]
    );
    return result;
  }
}
```

#### 3. Database Layer

**Миграции:**

```sql
-- 001_create_tables.up.sql
CREATE TYPE user_role AS ENUM ('admin', 'manager', 'viewer');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    role user_role NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    description TEXT,
    created_by VARCHAR(50),
    updated_by VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE item_history (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    action VARCHAR(20) NOT NULL CHECK (action IN ('INSERT', 'UPDATE', 'DELETE')),
    old_data JSONB,
    new_data JSONB,
    changed_by VARCHAR(50) NOT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_item_history_item_id ON item_history(item_id);
CREATE INDEX idx_item_history_changed_at ON item_history(changed_at DESC);
```

**Триггеры (АНТИПАТТЕРН):**

```sql
-- 002_create_triggers.up.sql
CREATE OR REPLACE FUNCTION log_item_change()
RETURNS TRIGGER AS $$
DECLARE
    user_context VARCHAR(50);
BEGIN
    -- Получаем пользователя из текущего контекста (устанавливается в приложении)
    user_context := current_setting('app.current_user', true);

    IF TG_OP = 'INSERT' THEN
        INSERT INTO item_history (item_id, action, new_data, changed_by)
        VALUES (NEW.id, 'INSERT', row_to_json(NEW), COALESCE(user_context, 'system'));
        RETURN NEW;

    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO item_history (item_id, action, old_data, new_data, changed_by)
        VALUES (NEW.id, 'UPDATE', row_to_json(OLD), row_to_json(NEW), COALESCE(user_context, 'system'));
        RETURN NEW;

    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO item_history (item_id, action, old_data, changed_by)
        VALUES (OLD.id, 'DELETE', row_to_json(OLD), COALESCE(user_context, 'system'));
        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_items_insert
    AFTER INSERT ON items
    FOR EACH ROW EXECUTE FUNCTION log_item_change();

CREATE TRIGGER trigger_items_update
    AFTER UPDATE ON items
    FOR EACH ROW EXECUTE FUNCTION log_item_change();

CREATE TRIGGER trigger_items_delete
    AFTER DELETE ON items
    FOR EACH ROW EXECUTE FUNCTION log_item_change();
```

**Seed Data:**

```sql
-- 003_seed_data.up.sql
INSERT INTO users (username, role, password_hash) VALUES
('admin', 'admin', '$2b$10$...'), -- hashed password
('manager', 'manager', '$2b$10$...'),
('viewer', 'viewer', '$2b$10$...');

INSERT INTO items (name, quantity, price, description, created_by) VALUES
('Laptop Dell XPS 15', 10, 1299.99, 'High-performance laptop', 'admin'),
('Mouse Logitech MX Master', 50, 99.99, 'Wireless mouse', 'manager'),
('Keyboard Keychron K2', 30, 89.99, 'Mechanical keyboard', 'manager');
```

### Потоки данных

**1. Аутентификация:**
```
Frontend (LoginForm)
  → POST /auth/login { username, role }
  → AuthController.login()
  → AuthService.validateCredentials()
  → AuthService.generateToken()
  → Response { accessToken, user }
  → Frontend stores token in localStorage
```

**2. Создание товара:**
```
Frontend (ItemForm)
  → POST /items { name, quantity, price, description }
  → AuthMiddleware (verifies JWT)
  → RolesGuard (checks role: admin/manager)
  → ItemsController.create()
  → ItemsService.create()
    → AuthorizationService.canCreate()
    → ItemValidator.validate()
  → ItemsRepository.create()
  → Database: INSERT INTO items
  → Trigger: log_item_change() fires
  → Database: INSERT INTO item_history
  → Response: created Item
  → Frontend updates Redux store
```

**3. Просмотр истории:**
```
Frontend (HistoryView)
  → GET /items/{id}/history
  → AuthMiddleware (verifies JWT)
  → RolesGuard (checks role: admin/manager/viewer)
  → ItemsController.getHistory()
  → ItemsService.getHistory()
  → ItemsRepository.findHistory()
  → Database: SELECT * FROM item_history WHERE item_id = ?
  → Response: HistoryEntry[]
  → Frontend displays timeline/table
```

### Преимущества и недостатки

**Преимущества:**
- Чёткое разделение ответственности (SRP)
- Легкая тестируемость каждого слоя (unit tests)
- Возможность замены реализации Repository без изменения бизнес-логики
- Dependency Injection упрощает управление зависимостями
- DTO обеспечивают валидацию на входе
- Чистый код с явными интерфейсами
- Легкая масштабируемость (добавление новых фич)
- Возможность использования CQRS для сложных сценариев

**Недостатки:**
- Большее количество файлов и слоёв
- Более сложная начальная настройка
- Избыточность для простых проектов
- Триггеры всё ещё скрывают логику истории в БД
- Сложнее отслеживать полный поток данных
- Требует больше времени на разработку
- Необходимость в DI контейнере

### Сравнение двух вариантов

| Характеристика | Вариант 1 (Монолит) | Вариант 2 (Слоистый) |
|---------------|-------------------|---------------------|
| Сложность | Низкая | Средняя |
| Время разработки | Быстрее | Медленнее |
| Тестируемость | Средняя | Высокая |
| Поддерживаемость | Средняя | Высокая |
| Масштабируемость | Низкая | Высокая |
| Чистота кода | Средняя | Высокая |
| Количество слоёв | 3 | 5+ |
| Использование паттернов | Минимальное | Repository, DI, DTO |
| Подходит для MVP | Да | Нет |
| Подходит для Enterprise | Нет | Да |

---

## Рекомендации

### Когда выбрать Вариант 1 (Монолит):
- Проект является MVP или прототипом
- Ограниченные сроки разработки
- Маленькая команда разработчиков
- Ожидается небольшой объём кода
- Не планируется масштабирование

### Когда выбрать Вариант 2 (Слоистый):
- Проект рассчитан на долгосрочную поддержку
- Команда разработчиков более 2-3 человек
- Ожидается рост функциональности
- Требуется высокая тестируемость
- Планируется интеграция с другими системами

### Важное примечание о триггерах

Оба варианта используют триггеры для логирования истории изменений, что является **антипаттерном**. В реальных проектах рекомендуется:

1. **Application-level logging:** Записывать историю в сервисном слое
2. **Event Sourcing:** Хранить все события изменений
3. **Audit Log Service:** Отдельный сервис для аудита
4. **CDC (Change Data Capture):** Использовать Debezium или аналоги

**Проблемы триггеров:**
- Скрытая бизнес-логика
- Сложность отладки
- Проблемы с миграциями
- Зависимость от конкретной СУБД
- Синхронное выполнение (блокировка)
- Невозможность добавить дополнительную логику

**Альтернативы:**
```typescript
// Вместо триггера - логирование в сервисе
async create(dto: CreateItemDto, user: UserPayload): Promise<Item> {
  const item = await this.itemsRepository.create(dto);
  
  // Явное логирование
  await this.historyService.log({
    itemId: item.id,
    action: 'INSERT',
    newData: item,
    changedBy: user.username
  });
  
  return item;
}
```

---

## Технологический стек (рекомендации)

### Backend:
- **Язык:** TypeScript / Python / Go
- **Фреймворк:** NestJS / FastAPI / Gin
- **ORM:** TypeORM / Prisma / SQLAlchemy / GORM
- **Валидация:** class-validator / Pydantic
- **Аутентификация:** Passport.js / python-jose
- **База данных:** PostgreSQL 14+

### Frontend:
- **Фреймворк:** React / Vue 3
- **State Management:** Redux Toolkit / Pinia
- **HTTP Client:** Axios
- **UI Components:** Material-UI / Ant Design / Vuetify
- **Routing:** React Router / Vue Router

### DevOps:
- **Контейнеризация:** Docker
- **База данных:** PostgreSQL
- **Миграции:** TypeORM Migrations / Alembic
- **Тестирование:** Jest / Pytest

---

## Заключение

Оба архитектурных варианта решают поставленную задачу, но имеют разные подходы к организации кода. Выбор зависит от целей проекта, размера команды и планов на развитие.

Важно помнить, что использование триггеров для логирования истории является учебным антипаттерном и не рекомендуется для production-систем. В реальных проектах следует использовать application-level logging или event sourcing подходы.
