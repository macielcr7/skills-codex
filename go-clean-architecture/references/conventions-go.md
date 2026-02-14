## Convenções para Go

### 1. Estrutura de Pastas e Arquivos
- **snake_case**: Para arquivos Go, use nomes em snake_case (letras minúsculas separadas por underscores).
- **CamelCase**: Para nomes de pacotes exportados, siga a convenção CamelCase.
- **lowercase**: Para pacotes (diretórios), use nomes em minúsculas sem underscores.

### 2. Estrutura do Projeto (baseada em golang-standards/project-layout)
```
/
├── cmd/                    # Aplicações principais do projeto
│   ├── app/                # Nome da aplicação
│   │   └── main.go         # Ponto de entrada da aplicação
│   └── cli/                # Ferramentas de linha de comando
├── internal/               # Código privado da aplicação
├── pkg/                    # Bibliotecas que podem ser usadas por aplicações externas
├── api/                    # Especificações de API, protobuf, etc.
├── web/                    # Componentes web (templates, assets, etc.)
├── configs/                # Arquivos de configuração
├── deployments/            # Configurações de implantação (docker, k8s, etc.)
└── docs/                   # Documentação do projeto
```

### 3. Nomes de Tipos e Interfaces
- **PascalCase** para tipos exportados (UserService, OrderRepository).
- **camelCase** para variáveis e funções não exportadas (calculatePrice, getUserByID).
- Evite redundância: Não use UserServiceStruct, apenas UserService.
- Interfaces com um único método geralmente são nomeadas com o sufixo "er" (Reader, Writer).

### 4. Nomes de Variáveis e Funções
- **camelCase** para variáveis e funções não exportadas (calculatePrice, getUserByID).
- **PascalCase** para variáveis e funções exportadas (ProcessPayment, SendEmail).
- Use nomes claros e descritivos, evitando abreviações excessivas (GetUserData e não gud).
- Acrônimos em nomes devem ser tratados como uma palavra (HttpServer → HTTPServer, Api → API).
- Nunca use Impl como sufixo para dizer que está implementando uma interface.

### 5. Convenções para Testes
- Arquivos de teste têm o sufixo _test.go (user_service_test.go).
- Funções de teste começam com Test seguido do nome da função testada (TestUserService_Process).
- Benchmarks começam com Benchmark (BenchmarkUserService_Process).
- Exemplos começam com Example (ExampleUserService_Process).
- Os testes ficam no mesmo pacote que o código testado.

### 6. Importações
- Organize as importações em grupos:
  1. Pacotes da biblioteca padrão
  2. Pacotes de terceiros
  3. Pacotes internos do projeto
- Use o caminho completo de importação baseado no módulo Go.

Exemplo:
```go
import (
    "context"
    "fmt"
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    
    "github.com/seu-usuario/seu-projeto/internal/domain"
)
```

### 7. Tratamento de Erros
- Retorne erros explicitamente em vez de usar panics.
- Use pacotes como "errors" ou "github.com/pkg/errors" para criar e enriquecer erros.
- Verifique erros imediatamente após a chamada que pode gerá-los.
- Evite usar _ para ignorar erros, a menos que seja absolutamente necessário.

### 8. Documentação
- Todos os pacotes e funções/tipos exportados devem ter comentários de documentação.
- Comentários de documentação começam com o nome do elemento que estão documentando em inglês.
- Use frases completas com ponto final.

Exemplo:
```go
// UserService providers methods...
type UserService struct {
    // ...
}

// Process process the payment
// Returns a response
func (s *UserService) Process(ctx context.Context, req Request) (Response, error) {
    // ...
}
```