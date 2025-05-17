# LemBRAGO

LemBRAGO é uma aplicação backend projetada para oferecer uma maneira segura e privada de armazenar, gerenciar e compartilhar suas senhas e informações sensíveis. Construído sobre uma arquitetura Zero-Knowledge, garantimos que nem mesmo nós (os desenvolvedores do serviço) podemos acessar suas senhas descriptografadas. Seus dados são criptografados e descriptografados localmente no seu dispositivo usando sua Senha Mestra, que nunca sai do seu computador.

## Funcionalidades

* **Gerenciamento de Organizações:** Crie e gerencie organizações.
* **Gerenciamento de Usuários:**
    * Registro de usuários dentro de uma organização.
    * Autenticação via e-mail e códigos de 6 dígitos.
    * Login de usuário e gerenciamento de sessão usando JWT.
    * Convidar usuários para uma organização.
    * Controle de acesso baseado em função (Admin, Member).
* **Gerenciamento de Cofres (Vaults):**
    * Crie, atualize e exclua cofres (compartilhados e pessoais).
    * Gerencie membros do cofre e suas permissões (Admin, Write, Read).
    * Armazene e gerencie senhas criptografadas dentro dos cofres.
* **Gerenciamento de Mídia:**
    * Faça upload e sirva arquivos de mídia de forma segura, associados a uma organização.
    * Exclusão de arquivos de mídia.
* **Versionamento da Aplicação:**
    * Registre e gerencie versões da aplicação, particularmente para um cliente desktop.
    * Endpoints para buscar a versão mais recente e todas as versões da aplicação.
    * Upload e download de builds da aplicação desktop para diferentes plataformas/arquiteturas.
* **Segurança:**
    * Autenticação baseada em JWT para endpoints protegidos.
    * Limitação de taxa (rate limiting) para prevenir abuso.
    * Configuração de CORS.
    * Middleware de tratamento de erros.
    * Uso de Argon2id para verificação de senha.
    * Armazenamento criptografado de dados sensíveis (metadados de cofres, senhas, chaves de usuário).
* **Notificações por E-mail:**
    * E-mails de boas-vindas para novos usuários.
    * E-mails de convite.
    * E-mails com código de autenticação.
* **Ambiente Dockerizado:** Fácil configuração e deployment usando Docker e Docker Compose.

## Stack de Tecnologias

* **Linguagem:** Go
* **Framework:** Gin (Framework Web)
* **Banco de Dados:** MongoDB
* **Cache:** Redis
* **Containerização:** Docker, Docker Compose

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