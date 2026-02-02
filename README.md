# AITU AD Final Project

Backend for an peripherals store, paired with a Vite/React/Tailwind frontend. The API is built with Go (Gin) and MongoDB, secured with JWT, and deployed to Railway; the frontend is deployed to Vercel.

## Project Overview
- Domain: catalog, product reviews, orders, wishlists, admin stats, and user profiles.
- Auth: email/password with JWT, roles `user` and `admin` (middleware-enforced).
- Deploy targets: Railway (backend), Vercel (frontend), MongoDB on Railway.
- Docs: OpenAPI available at `/swagger/index.html` once the server is running (sources in `docs/swagger.yaml|json`).

## System Architecture
- **Frontend:** React (TypeScript) + Vite + Tailwind CSS; served from Vercel.
- **Backend:** Go 1.21+, Gin router, layered domain → repository → service → handler.
- **Database:** MongoDB (Atlas friendly). Collections: `users`, `products`, `categories`, `orders`, `wishlist`.
- **Auth:** JWT with Bearer tokens; role-based guards for admin routes.
- **Hosting/CI:** Railway for the API, Vercel for the SPA.

## Database Schema (MongoDB)
- `users`:
  - `_id` ObjectId
  - `name`, `email`, `password_hash`, `role` ("user"|"admin")
  - `address`, `phone`, `bio`, `created_at`
- `categories`:
  - `_id`, `name`, `description`, `createdAt`, `updatedAt`
- `products`:
  - `_id`, `categoryId` (ObjectId), `name`, `description`, `price` (float), `stock` (int)
  - `reviews` (embedded array): `_id`, `userId`, `rating`, `comment`, `createdAt`
  - `createdAt`, `updatedAt`
- `orders`:
  - `_id`, `userId` (string), `items` [{`productId` ObjectId, `quantity` int}]
  - `status` ("pending"|"shipped"|"delivered"|"cancelled")
  - `createdAt`, `updatedAt`
- `wishlist`:
  - `_id`, `userId` (string), `productId` (ObjectId), `createdAt`

## Representative MongoDB Queries
- List products with paging and optional category filter:
  ```js
  db.products.find(
    { ...(categoryId && { categoryId: ObjectId(categoryId) }) }
  ).sort({ createdAt: -1 }).skip(offset).limit(limit)
  ```
- Add review to a product (embedded push):
  ```js
  db.products.updateOne(
    { _id: ObjectId(productId) },
    { $push: { reviews: { _id: ObjectId(), userId, rating, comment, createdAt: new Date() } } }
  )
  ```
- Order status update:
  ```js
  db.orders.findOneAndUpdate(
    { _id: ObjectId(orderId) },
    { $set: { status, updatedAt: new Date() } },
    { returnDocument: "after" }
  )
  ```
- Sales statistics (aggregate excerpt):
  ```js
  db.orders.aggregate([
    { $match: dateFilter },
    { $facet: {
        statusCounts: [{ $group: { _id: "$status", count: { $sum: 1 } } }],
        totals: [
          { $unwind: "$items" },
          { $lookup: { from: "products", localField: "items.productId", foreignField: "_id", as: "product" } },
          { $unwind: { path: "$product", preserveNullAndEmptyArrays: true } },
          { $group: { _id: null,
            totalOrders: { $addToSet: "$_id" },
            totalRevenue: { $sum: { $multiply: [ "$items.quantity", { $ifNull: [ "$product.price", 0 ] } ] } }
          }},
          { $project: { totalOrders: { $size: "$totalOrders" }, totalRevenue: 1 } }
        ]
    }}
  ])
  ```
- Wishlist uniqueness check (compound index):
  ```js
  db.wishlist.createIndex({ userId: 1, productId: 1 }, { unique: true })
  ```

## Indexing & Optimization Strategy
- Unique index on `users.email` (`uniq_email`) to enforce unique accounts.
- Compound unique index on `wishlist.userId + productId` to prevent duplicates.
- Implicit `_id` indexes on all collections.
- Queries sort by `createdAt` and use `skip/limit`; keep `createdAt` indexed if large datasets grow.
- Aggregations reuse `$match` early to reduce pipeline volume; `$facet` used for combined stats in a single round trip.
- Suggested future tuning: add `orders.userId` index for user-specific lists; add `products.categoryId` index to speed catalog filtering.

## API Surface (v1)
Base path: `/api/v1` (Swagger: `/swagger/index.html`)

- **Health**
  - `GET /health` — public

- **Auth**
  - `POST /auth/register` — public user signup
  - `POST /auth/login` — public login
  - `POST /admin/auth/register` — admin creates admin

- **Categories (public + admin)**
  - `GET /categories`
  - `GET /categories/:id`
  - `POST /admin/categories` — admin
  - `PUT /admin/categories/:id` — admin
  - `DELETE /admin/categories/:id` — admin

- **Products**
  - `GET /products`
  - `GET /products/:id`
  - `POST /products/:id/reviews` — auth user
  - `DELETE /products/:id/reviews/:reviewId` — auth user
  - `POST /admin/products` — admin
  - `PUT /admin/products/:id` — admin
  - `DELETE /admin/products/:id` — admin

- **Orders**
  - `POST /orders` — auth user
  - `GET /orders` — auth user/admin (user gets own, admin sees all)
  - `GET /orders/:id` — auth user/admin (own or any for admin)
  - `PUT /admin/orders/:id/status` — admin
  - `GET /admin/orders/:id` — admin
  - `POST /admin/orders/find` — admin (lookup by id)

- **Wishlist** (auth user)
  - `POST /wishlist`
  - `GET /wishlist`
  - `DELETE /wishlist/:id`

- **Profile** (auth user)
  - `GET /profile`
  - `PUT /profile`

- **Admin stats**
  - `GET /admin/stats/sales` — admin (query: `year` or `start`+`end`)
  - `GET /admin/stats/products` — admin (same query pattern)

- **Admin users**
  - `GET /admin/users` — admin (list all users)

Full contract and schemas: `/swagger/index.html` or `docs/swagger.yaml`.

## Deployment Notes
- Railway (prod backend): `https://aitu-ad-final-back-production.up.railway.app/api/v1`
- Vercel (prod frontend/admin): `https://mangustad.vercel.app/admin/dashboard`
- Railway service exposes the Gin server on `PORT`.
- Env vars for Railway/Vercel must mirror `.env` keys; never commit secrets.
- Frontend hits the backend base URL configured per environment; update the SPA env to match the current Railway URL.

## Contributions
- **Frontend**
  - Nursultan: Admin pages with API handling — Orders v1, Products v1, Login, Register, Categories, OrderDetails, Stats v1.
  - Birlik: User pages with API + styling — Catalog, ProductDetail, OrderList, OrderDetail, Wishlist, Profile, Landing page, Stats (Admin) v2, Orders (Admin) v2, AddAdmin.
- **Backend**
  - Nursultan: Authentication, Registration, Orders CRUD, Products CRUD, Categories CRUD.
  - Birlik: AddAdmin, Wishlist CRUD, Profile, Product reviews, Get all users.