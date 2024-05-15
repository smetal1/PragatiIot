
# Pragati IOT

PragatiIoT is a robust IoT and Industrial IoT framework written in Go, designed to provide seamless device management, data streaming, and user interaction capabilities. This framework facilitates the creation of scalable IoT solutions, enabling real-time data processing and control.


## Features

- Device Management: Add, update, and manage IoT devices.
- User Roles: Supports role-based access control with Admin and View roles.
- Data Streaming: Utilizes MQTT for real-time data communication.
- Secure Communication: TLS support for secure MQTT communication.
- Database Integration: PostgreSQL for data storage and management.
- Scalable Architecture: Docker and Kubernetes for deployment.




## Installation

Prerequisites
- Go 1.21 or later
- Docker
- Docker Compose
- PostgreSQL
- RabbitMQ
- Mosquitto MQTT broker (or any other MQTT broker)

## Clone the Repository
```bash
  git clone https://github.com/yourusername/PragatiIoT.git
  cd PragatiIoT
```
## Setup Environment Variables

Create a .env file in the platform directory with the following content:

```bash
  POSTGRES_USER=your_postgres_user
  POSTGRES_PASSWORD=your_postgres_password
  POSTGRES_DB=your_database_name
```

## Generate Certificates (Ignore it for now)
Run the script to generate TLS certificates for MQTT:

```bash
  cd platform
  ./generate_cert.sh
```
## Database Initialization
Initialize the PostgreSQL database:

```bash
  docker-compose up db
  docker-compose exec db psql -U your_postgres_user -d your_database_name -f db_init.sql
```

## Start PostGres, RabbitMQ, MQTT Broker Services

Use Docker Compose to start the services:

```bash
  docker-compose up -d
```

## Start the Platform

```bash
  cd pragatiiot/platform
  go run main.go
```

## Usage
### API Endpoints
The API documentation is available via Swagger. Access it at http://localhost:8080/swagger/index.html.

### MQTT Configuration
Configure your IoT devices to connect to the MQTT broker at mqtt://localhost:1883 using the generated certificates.

# Contributing
Contributions are welcome! Please fork the repository and submit a pull request for review.

# License
This project is licensed under the Creative Commons Attribution-NonCommercial 4.0 International License. See the LICENSE file for details.

# Acknowledgements
Special thanks to the open-source community for their invaluable tools and resources.
- Gin Web Framework
- Mosquitto
For any questions or support, please contact 
[saurav7055@gmail.com].
## License

[CC BY-NC 4.0](https://creativecommons.org/licenses/by-nc/4.0/)

