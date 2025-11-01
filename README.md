##E-Commerce Platform

A simple full-stack e-commerce web application built with React (frontend) and Go / GORM / SQLite (backend). Supports browsing products, adding to cart, and placing orders.

##Features

Browse products with name, description, and price.

Add products to cart and update quantity.

View cart summary with total price.

Checkout modal confirms order and clears cart.

Backend APIs built in Go with SQLite for storage.

##Tech Stack

Frontend: React, Vite

Backend: Go, Chi router, GORM, SQLite

Other: CORS handling, modal for checkout

##Getting Started
Prerequisites

Node.js >= 18

Go >= 1.20

##Setup Frontend
cd frontend
npm install
npm run dev

Frontend runs at http://localhost:5173.

##Setup Backend
cd backend
go run main.go

Backend runs at http://localhost:8080.

##API Endpoints

GET /products – list all products

GET /cart/{sessionId} – get cart for session

POST /cart/items – add item to cart

DELETE /cart/items/{sessionId} – clear cart for session
