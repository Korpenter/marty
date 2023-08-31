# Marty
Marty is an HTTP API for a loyalty points application that provides various functionalities to users. It offers features such as user registration, authentication, authorization, and more.

## Usage
Application is configured using flags or environmental variables:
- Service Address (-a or RUN_ADDRESS): This specifies the address and port on which the Marty service will be hosted. The default value is localhost:8081.

- Postgres URI (-d or DATABASE_URI): This defines the URI for connecting to the PostgreSQL database. If not provided, the application will use the default value.

- Accrual Address (-r or ACCRUAL_SYSTEM_ADDRESS): This sets the address of the accrual system used for calculating and managing points accrual. The default value is localhost:8080.

- Secret Key (SECRET_KEY): This is the secret key used for various cryptographic operations, such as hashing and encryption.

### Docker
Build image:

```shell
sudo docker buildx build -t marty .  
```

Run it:

```shell
docker run marty
```
