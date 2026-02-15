# Tecnologias e Pacotes Go Utilizados no Projeto

## **Principais Recursos da Biblioteca Padrão do Go**

### **1. Logging**

- `log/slog`: Biblioteca moderna para logging estruturado.

## **Principais Pacotes Externos Utilizados**

### **1. Banco de Dados**

- SQL puro usando `database/sql`
- `github.com/golang-migrate/migrate/v4`: Ferramenta para controle de **migrations** do banco de dados.

Obs.: versões recentes de dependências podem exigir versões mais novas do Go; se precisar de compatibilidade com um Go mais antigo, fixe versões compatíveis via `go.mod`.

### **2. Web Framework e API REST**

- `github.com/go-chi/chi/v5`: Roteador minimalista e eficiente para APIs HTTP.

### **3. Upload para AWS S3**

- `github.com/aws/aws-sdk-go-v2`: SDK oficial da AWS para integração com o S3.
- `github.com/localstack/localstack-go`: Para testar o upload para S3 localmente.

### **4. Variáveis de Ambiente**

- `github.com/joho/godotenv`: Para carregar variáveis de ambiente de um arquivo `.env`.

### **5. Identificadores Únicos (UUID)**

- `github.com/google/uuid`: Biblioteca para geração de UUIDs, garantindo identificadores únicos para registros no banco de dados e outras entidades da aplicação.

### **6. Testes automatizados**

- `github.com/stretchr/testify`: Framework de testes para facilitar a escrita de testes unitários e de integração.

### **7. Validação**

- `github.com/asaskevich/govalidator`: Validação baseada em tags/structs (útil em entidades do domínio quando necessário).
