# Current Implementation
Techstack:
 - go
 - gin
 - mongodb
 - paseto (for login)


The 'Withdraw' functionality allows a user to withdraw a certain amount from their account. 

The 'Transfer' functionality allows a user to transfer a certain amount from their account to another user's account. Both functionalities check if the user has sufficient balance before proceeding with the transaction.  

Current application uses the PASETO (Platform-Agnostic Security Tokens) protocol for handling authentication. The user's authentication token is retrieved from the header and decoded to get the user's information.

This implementation also using OpenAPI for API documentation. The API documentation can be accessed at http://localhost:8080/openapi.json
## How to run
Before running, this app need to generate public and private key for PASETO. The public key will be used to verify the token and the private key will be used to sign the token. The public and private key can be generated using the following command:
```bash 
go run initial/main.go
```

After running, set Public and Private key in .env file. The .env file should be in the root directory of the project. The .env file should contain the following:
```bash
PASETO_SECRET=
PASETO_PUBLIC=
MONGODB_URI=
```

```bash
docker compose up -d
```

