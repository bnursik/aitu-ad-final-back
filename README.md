# AITU AD Final Backend

A RESTful API backend built with Go, Gin framework, and MongoDB.

## Getting Started

### Prerequisites
- Go 1.21+
- MongoDB

### Environment Variables

Create a `.env` file in the root directory:

```env
MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/
DB_NAME=your_database_name
JWT_SECRET=your_jwt_secret
PORT=8080
```

### Running the Application

```bash
go run cmd/api/main.go
```

### Testing MongoDB Connection

```bash
go run cmd/test/main.go
```

---

## API Documentation

Base URL: `/api/v1`

### Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

---

## Endpoints

### Health Check

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/health` | Check API health | No |

**Response:**
```json
{"status": "ok"}
```

---

### Auth

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/auth/register` | Register new user | No |
| POST | `/auth/login` | Login user | No |
| POST | `/admin/auth/admin/register` | Register new admin | Admin |

#### Register User
**POST** `/auth/register`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "...",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

#### Login
**POST** `/auth/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "...",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

#### Register Admin
**POST** `/admin/auth/admin/register`

**Request Body:**
```json
{
  "name": "Admin User",
  "email": "admin@example.com",
  "password": "adminpass123"
}
```

---

### Categories

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/categories` | List all categories | No |
| GET | `/categories/:id` | Get category by ID | No |
| POST | `/admin/categories` | Create category | Admin |
| PUT | `/admin/categories/:id` | Update category | Admin |
| DELETE | `/admin/categories/:id` | Delete category | Admin |

#### List Categories
**GET** `/categories`

**Response (200):**
```json
[
  {
    "id": "...",
    "name": "Electronics",
    "description": "Electronic devices",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

#### Get Category
**GET** `/categories/:id`

**Response (200):**
```json
{
  "id": "...",
  "name": "Electronics",
  "description": "Electronic devices",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Create Category (Admin)
**POST** `/admin/categories`

**Request Body:**
```json
{
  "name": "Electronics",
  "description": "Electronic devices"
}
```

#### Update Category (Admin)
**PUT** `/admin/categories/:id`

**Request Body:**
```json
{
  "name": "Updated Name",
  "description": "Updated description"
}
```

#### Delete Category (Admin)
**DELETE** `/admin/categories/:id`

**Response:** `204 No Content`

---

### Products

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/products` | List all products | No |
| GET | `/products/:id` | Get product by ID | No |
| POST | `/admin/products` | Create product | Admin |
| PUT | `/admin/products/:id` | Update product | Admin |
| DELETE | `/admin/products/:id` | Delete product | Admin |

#### List Products
**GET** `/products`

**Query Parameters:**
- `categoryId` (optional): Filter by category ID

**Response (200):**
```json
[
  {
    "id": "...",
    "categoryId": "...",
    "name": "iPhone 15",
    "description": "Latest iPhone",
    "price": 999.99,
    "stock": 100,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

#### Get Product
**GET** `/products/:id`

**Response (200):**
```json
{
  "id": "...",
  "categoryId": "...",
  "name": "iPhone 15",
  "description": "Latest iPhone",
  "price": 999.99,
  "stock": 100,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
  "reviews": [
    {
      "id": "...",
      "userId": "...",
      "rating": 5,
      "comment": "Great product!",
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### Create Product (Admin)
**POST** `/admin/products`

**Request Body:**
```json
{
  "categoryId": "...",
  "name": "iPhone 15",
  "description": "Latest iPhone",
  "price": 999.99,
  "stock": 100
}
```

#### Update Product (Admin)
**PUT** `/admin/products/:id`

**Request Body:**
```json
{
  "categoryId": "...",
  "name": "iPhone 15 Pro",
  "description": "Updated description",
  "price": 1199.99,
  "stock": 50
}
```

#### Delete Product (Admin)
**DELETE** `/admin/products/:id`

**Response:** `204 No Content`

---

### Product Reviews

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/products/:id/reviews` | Add review | User |
| DELETE | `/products/:id/reviews/:reviewId` | Delete review | User |

#### Add Review
**POST** `/products/:id/reviews`

**Request Body:**
```json
{
  "rating": 5,
  "comment": "Great product!"
}
```

**Response (201):**
```json
{
  "id": "...",
  "userId": "...",
  "rating": 5,
  "comment": "Great product!",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

#### Delete Review
**DELETE** `/products/:id/reviews/:reviewId`

**Response:** `204 No Content`

---

### Orders

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/orders` | List orders (user: own, admin: all) | User/Admin |
| GET | `/orders/:id` | Get order by ID | User/Admin |
| POST | `/orders` | Create order | User |
| PUT | `/admin/orders/:id/status` | Update order status | Admin |
| POST | `/admin/orders/find` | Find order by ID | Admin |

#### List Orders
**GET** `/orders`

**Response (200):**
```json
[
  {
    "id": "...",
    "items": [
      {
        "productId": "...",
        "quantity": 2
      }
    ],
    "status": "pending",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "userId": "..." // Only for admin
  }
]
```

#### Get Order
**GET** `/orders/:id`

**Response (200):**
```json
{
  "id": "...",
  "items": [
    {
      "productId": "...",
      "quantity": 2
    }
  ],
  "status": "pending",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Create Order
**POST** `/orders`

**Request Body:**
```json
{
  "items": [
    {
      "productId": "...",
      "quantity": 2
    }
  ]
}
```

**Response (201):**
```json
{
  "id": "...",
  "items": [
    {
      "productId": "...",
      "quantity": 2
    }
  ],
  "status": "pending",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Update Order Status (Admin)
**PUT** `/admin/orders/:id/status`

**Request Body:**
```json
{
  "status": "shipped"
}
```

**Valid statuses:** `pending`, `shipped`, `delivered`, `cancelled`

#### Find Order by ID (Admin)
**POST** `/admin/orders/find`

**Request Body:**
```json
{
  "order_id": "..."
}
```

**Response (200):**
```json
{
  "id": "...",
  "userId": "...",
  "items": [...],
  "status": "pending",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

---

### Wishlist

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/wishlist` | Get user's wishlist | User |
| POST | `/wishlist` | Add product to wishlist | User |
| DELETE | `/wishlist/:id` | Remove from wishlist | User |

#### Get Wishlist
**GET** `/wishlist`

**Response (200):**
```json
[
  {
    "id": "...",
    "productId": "...",
    "createdAt": "2024-01-01T00:00:00Z"
  }
]
```

#### Add to Wishlist
**POST** `/wishlist`

**Request Body:**
```json
{
  "product_id": "..."
}
```

**Response (201):**
```json
{
  "id": "...",
  "productId": "...",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

#### Remove from Wishlist
**DELETE** `/wishlist/:id`

**Response:** `204 No Content`

---

### Statistics (Admin Only)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/admin/statistics/sales` | Get all sales statistics | Admin |
| POST | `/admin/statistics/sales/date-range` | Get sales stats by date range | Admin |
| POST | `/admin/statistics/sales/year` | Get sales stats by year | Admin |
| GET | `/admin/statistics/products` | Get all products statistics | Admin |
| POST | `/admin/statistics/products/date-range` | Get products stats by date range | Admin |
| POST | `/admin/statistics/products/year` | Get products stats by year | Admin |

#### Get All Sales Statistics
**GET** `/admin/statistics/sales`

**Response (200):**
```json
{
  "total_orders": 150,
  "total_revenue": 25000.50,
  "average_order": 166.67,
  "pending_orders": 10,
  "shipped_orders": 50,
  "delivered_orders": 85,
  "cancelled_orders": 5
}
```

#### Get Sales Statistics by Date Range
**POST** `/admin/statistics/sales/date-range`

**Request Body:**
```json
{
  "start_date": "2024-01-01",
  "end_date": "2024-12-31"
}
```

#### Get Sales Statistics by Year
**POST** `/admin/statistics/sales/year`

**Request Body:**
```json
{
  "year": 2024
}
```

#### Get All Products Statistics
**GET** `/admin/statistics/products`

**Response (200):**
```json
{
  "total_products": 500,
  "total_stock": 10000,
  "out_of_stock": 15,
  "total_reviews": 2500,
  "average_rating": 4.2,
  "total_categories": 25
}
```

#### Get Products Statistics by Date Range
**POST** `/admin/statistics/products/date-range`

**Request Body:**
```json
{
  "start_date": "2024-01-01",
  "end_date": "2024-12-31"
}
```

#### Get Products Statistics by Year
**POST** `/admin/statistics/products/year`

**Request Body:**
```json
{
  "year": 2024
}
```

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "error message"
}
```

### Common HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 204 | No Content |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 500 | Internal Server Error |

---

## Swagger Documentation

Swagger UI is available at: `/swagger/index.html`
