# 🛒 E-commerce (LoLoL200)

## 📌 Project Description

This is a simple educational e-commerce project that implements the core functionality of an online store: browsing products, adding items to a cart, placing orders, and user management.

The project demonstrates a classic web application structure with separation between frontend and backend logic.

---

## ⚙️ Main Features

### 👤 User Side:

* Browse products (catalog)
* Search and filter products
* View product details
* Add products to cart
* Manage cart (update quantity, remove items)
* User registration and login
* Checkout process
* User profile page

### 🛠 Admin Panel (if implemented):

* Manage products (create / edit / delete)
* Manage orders
* Manage users

---

## 🧱 Technologies Used

* HTML — page structure
* CSS / Bootstrap — styling and responsiveness
* JavaScript — interactivity
* PHP — backend logic
* MySQL — database

---

## 📂 Project Structure (example)

```
/css        — styles
/js         — scripts
/images     — images
/includes   — reusable components
/admin      — admin panel
/database   — database logic

index.php        — home page
product.php      — product page
cart.php         — cart
checkout.php     — checkout
login.php        — login
signup.php       — registration
profile.php      — user profile
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

* Apache or Nginx
* PHP (version 7 or higher)
* MySQL

### 3. Database setup

* Create a database in MySQL
* Import the `.sql` file (if provided in the project)
* Configure database connection settings in the config file (host, username, password, database name)

### 4. Run the project

Place the project into your web server directory:

```
/var/www/html   (Linux)
htdocs          (XAMPP)
```

Start your server and open in browser:

```
http://localhost/
```

Or run using PHP built-in server:

```bash
php -S localhost:8000
```

---

## 💻 Useful Commands

### Git commands

```bash
git clone https://github.com/LoLoL200/E-commerce.git
git pull
git add .
git commit -m "update"
git push
```

### Run local server

```bash
php -S localhost:8000
```

---

## 📦 Possible Improvements

* Integrate payment systems (Stripe / PayPal)
* Add REST API
* Improve security (password hashing, SQL injection protection)
* Optimize performance

---

## 📎 Summary

This project is a basic implementation of an online store with key features such as product catalog, cart, authentication, and checkout. It is suitable for learning purposes and can be used as a foundation for more advanced e-commerce systems.
