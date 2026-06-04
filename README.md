# 🛒 E-commerce (LoLoL200)

## 📌 Project Description

The E-commerce project is a REST API application that implements the business logic.

---

## ⚙️ Main Features
==============================================

### 👤 User Side:

*  REGISTER
*  LOGIN
*  UPDATE TOKEN
*  LOGOUT
==============================================
### 🍏,🖥️,🎨,📱 Product Side:

*  List Product
*  Search Product
*  Products by Category
*  Dateils product
==============================================
### 🛒 Cart Side:

*  GET cart user
*  Clear cart
*  Add new product in cart
*  Update quantity
==============================================
### 🛍️ Order Side:

*  Create new order from cart
*  List user orders
*  Datail one order
*  Cancel order
==============================================

---

## 🧱 Technologies Used

* Golang - base all project
* Docker - For a virtual server and testing
* Swagger - For a Frontend Developer
* PostgreSQL — database

---

## 📂 Project Structure (example)

```
.
├── cmd
│   └── api
│       └── main.go
├── docker-compose.yml
├── Dockerfile.dev
├── go.mod
├── go.sum
├── internal
│   ├── db
│   │   └── postgres.go
│   ├── domain
│   │   ├── cart.go
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   ├── handler
│   │   ├── http
│   │   │   ├── auth_handler.go
│   │   │   ├── cart_handler.go
│   │   │   ├── middleware.go
│   │   │   ├── order_handler.go
│   │   │   ├── product_handler.go
│   │   │   ├── router.go
│   │   │   ├── user_handler.go
│   │   │   └── utils.go
│   │   └── middleware
│   │       ├── product.go
│   │       └── users.go
│   ├── repository
│   │   ├── mocks
│   │   └── postgres
│   │       ├── cart_repository.go
│   │       ├── mocks
│   │       │   ├── cart_mock.go
│   │       │   └── order_mock.go
│   │       ├── orders_repository.go
│   │       ├── product_repository.go
│   │       └── user_repository.go
│   └── service
│       ├── auth
│       │   ├── auth_service.go
│       │   ├── dto_user.go
│       │   ├── mocks
│       │   │   ├── auth_mocks.go
│       │   │   └── mocks.go
│       │   └── user_service.go
│       ├── cart
│       │   └── cart_service.go
│       ├── order
│       │   ├── mocks
│       │   │   └── order_service_mock.go
│       │   └── order_service.go
│       └── product
│           ├── dto_product.go
│           └── product_service.go
├── migrations
│   ├── 000001_create_users_table.down.sql
│   ├── 000001_create_users_table.up.sql
│   ├── 000002_create_categories_table.down.sql
│   ├── 000002_create_categories_table.up.sql
│   ├── 000003_create_products_table.down.sql
│   ├── 000003_create_products_table.up.sql
│   ├── 000004_create_cart_items_table.down.sql
│   ├── 000004_create_cart_items_table.up.sql
│   ├── 000005_create_orders_tables.down.sql
│   └── 000005_create_orders_tables.up.sql
├── pkg
│   └── utils
│       └── errors.go
├── README.md
├── script
│   └── seed.sql
├── swager
│   ├── swager_auth.yaml
│   ├── swager_cart.yaml
│   ├── swager_order.yaml
│   └── swager_product.yaml
├── test
│   └── user_test.go
└── tmp
    ├── build-errors.log
    └── main
```

---

## 🚀 Installation & Setup

### 1. Clone the repository

```bash
git clone https://github.com/LoLoL200/E-commerce.git
cd E-commerce
```

### 2. Server setup

Install the following:

* Golang
* Swager
* Docker
* PostgreSQL

### 3. Database setup

* Create a database in MySQL
* Import the `.sql` file (if provided in the project)
* Configure database connection settings in the config file (host, username, password, database name)



## 📦 Possible Improvements
* ADMIN PANEL
* Integrate payment systems (Stripe / PayPal)
* Improve security (password hashing, SQL injection protection)
* Optimize performance

---

## 📎 Summary

This project is a basic implementation of an online store with key features such as product catalog, cart, authentication, and checkout. It is suitable for learning purposes and can be used as a foundation for more advanced e-commerce systems.
