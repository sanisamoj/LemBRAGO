### Iniciar um banco de dados Mongodb:

    docker run -d --name mongodb -e MONGO_INITDB_ROOT_USERNAME=jomasinas -e MONGO_INITDB_ROOT_PASSWORD=senha -p 27017:27017 mongodb/mongodb-community-server

### Para construir a imagem docker:

    docker build -t lembrago .

### Para rodar a aplicação

    docker run --name lembragro -p 7888:7888 lembrago:latest

### Caso queira utilizar o Docker-compose

    docker-compose -p lembrago up --build -d

### Formato do arquivo .env

```.env
HOST=0.0.0.0
PORT=7888
SELF_URL=

MONGO_HOST=host.docker.internal
MONGO_PORT=27017
MONGO_USER=jomasinas
MONGO_PASSWORD=senha

REDIS_HOST=host.docker.internal
REDIS_PORT=6379
REDIS_PASSWORD=senha

JWT_SECRET=secret_key
JWT_SECRET_USER_CREATION=secret_key_user_creation
JWT_SECRET_ADMIN=secret_key
JWT_ISSUER=LemBraGO

EMAIL_AUTH_USER=
EMAIL_AUTH_PASS=
EMAIL_HOST=
EMAIL_PORT=465
```