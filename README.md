# Flashy: E-tailing Website for Limited-Time Sales
Flashy is an e-tailing platform that offers a curated selection of products at discounted prices for a limited time. The website utilizes consistent hashing for horizontal scaling, reducing login verification time by 33%. Data persistence, portability, and security are improved by replacing sessions with cookies. The leaky bucket algorithm is employed for rate limiting on the front end. The project follows the Model-View-Controller design pattern
<br><br>
Built with Go, Iris, RabbitMQ, JavaScript, HTML, CSS, MySQL, and Docker.

## Table of Contents
[Features](#features) <br>
[Requirement](#requirement)<br>
[Installation](#requirement)<br>
[Usage](#usage)<br>
[Contribution](#contribution)<br>


## Features
* __Limited-Time Sales:__ Offers a select number of products at discounted prices for a short period.
* __Consistent Hashing:__ Achieves horizontal scaling, reducing login verification time by 33%.
* __Data Persistence:__ Replaces sessions with cookies for improved data persistence, portability, and security.
* __Rate Limiting:__ Implements the leaky bucket algorithm for rate limiting on the front end.
* __MVC Design Pattern:__ Organizes code using the Model-View-Controller design pattern.

## Requirement
* Go 1.16 or higher
* Docker
* MySQL
* RabbitMQ

## Installation

### clone the repo
```bash
git clone https://github.com/jliu0105/JMOOC.git
```

### Install the required dependencies:
```bash
go get
```

### Build the project:
```bash
go build
```
### Build the docker image
```bash
docker build -t flashy .
```

### Run the docker container:
```bash
docker run -p 8080:8080 --name flashy flashy
```
Access the application at http://localhost:8080.

## Usage
* Browse the limited-time sales offerings.
* Create an account or log in to access exclusive deals.
* Add desired products to your cart.
* Complete the checkout process to secure your discounted purchases.

## Contribution
We welcome contributions from the community! If you're interested in contributing to this project, please follow these steps:

1. Fork the repository
2. Create a new branch for your changes
3. Commit your changes and push them to your fork
4. Create a pull request for your changes
