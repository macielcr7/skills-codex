# Estrutura de Camadas, DDD e Responsabilidades

## Visão Geral

O ConversorGo segue uma arquitetura em camadas para organizar seu código, garantindo separação de responsabilidades, testabilidade e manutenibilidade. Este documento descreve a estrutura arquitetural do projeto, as responsabilidades de cada camada e as diretrizes para implementação.

## Princípios Fundamentais

1. **Separação de Responsabilidades**: Cada camada tem um propósito específico e bem definido
2. **Testabilidade**: Todas as camadas podem ser testadas de forma isolada
3. **Independência de Detalhes Técnicos**: O núcleo da aplicação não depende de detalhes de implementação
4. **Inversão de Dependência**: Dependências apontam para dentro, não para fora
5. **Substituibilidade**: Componentes podem ser substituídos sem afetar o restante do sistema

## Estrutura de Camadas

O ConversorGo é organizado nas seguintes camadas:

```
internal/
├── domain/             # Regras de negócio e entidades
├── application/        # Casos de uso e orquestração de serviços
└── infra/              # Implementações técnicas e interfaces com o mundo externo
```

### Extensão: organização por DDD (Bounded Context / Aggregate / Action)

Quando usar DDD, aplicar a mesma segmentação em todas as camadas para evitar “pacotes gigantes”:

```
internal/domain/
  entity/{bounded-context}/{aggregate}/...
  repository/{bounded-context}/{aggregate}/...

internal/application/
  usecase/{bounded-context}/{aggregate}/{action}.go
  service/{bounded-context}/...              # opcional
  service/shared/...                         # cross-cutting

internal/infra/
  api/handler/{bounded-context}/{aggregate}/...
  repository/{bounded-context}/{aggregate}/{tech}/...
```

### Regra adicional: dependências entre Bounded Contexts

Por padrão, um Bounded Context não deve importar diretamente entidades/repositórios de outro contexto.
Quando for inevitável, explicitar a exceção com um allowlist por contexto:

- `internal/domain/entity/{bounded-context}/.allowed_contexts` (1 contexto permitido por linha)

### 1. Domain (Domínio)

A camada de domínio contém as entidades e regras de negócio centrais da aplicação, independentes de qualquer detalhe de implementação.

#### Responsabilidades:
- Definir entidades de negócio
- Definir interfaces de repositórios
- Implementar regras de negócio puras

#### Estrutura:
```
domain/
├── entity/                             # Entidades e agregados
│   └── {bounded-context}/{aggregate}/  # DDD
└── repository/                         # Interfaces de repositórios
    └── {bounded-context}/{aggregate}/  # DDD
```

#### Diretrizes:
- Não deve depender de nenhuma outra camada
- Não deve importar pacotes externos exceto os da biblioteca padrão Go (código de produção), **salvo exceções explícitas** (ex.: validação com `github.com/asaskevich/govalidator`)
- Deve conter apenas regras de negócio puras
- Deve definir interfaces que serão implementadas por camadas externas

#### Exemplo:
```go
// domain/entity/media/video/video.go
package video

type Video struct {
    ID           string
    Title        string
    FilePath     string
    Status       string
    // ...
}

func (v *Video) CanBeProcessed() bool {
    return v.Status == "pending" || v.Status == "failed"
}

// domain/repository/media/video/video_repository.go
package video

type Repository interface {
    FindByID(ctx context.Context, id string) (Video, error)
    Create(ctx context.Context, v Video) error
}
```

### 2. Application (Aplicação)

A camada de aplicação contém os casos de uso da aplicação, orquestrando o fluxo entre entidades e serviços de infraestrutura.

#### Responsabilidades:
- Implementar casos de uso
- Orquestrar entre domínio e infraestrutura
- Gerenciar transações e fluxo de dados
- Implementar serviços de aplicação

#### Estrutura:
```
application/
├── service/                             # Serviços de aplicação (opcional)
│   ├── {bounded-context}/...
│   └── shared/...
└── usecase/                             # Casos de uso
    └── {bounded-context}/{aggregate}/   # DDD
```

#### Diretrizes:
- Pode depender apenas da camada de domínio
- Não deve conter regras de negócio, apenas orquestração
- Pode receber implementações de infraestrutura via injeção de dependência

#### Exemplo:
```go
// application/usecase/media/video/upload_video.go
package video

import (
    "github.com/devfullcycle/golangtechweek/internal/domain/repository/media/video"
)

type UploadVideo struct {
    repo video.Repository
}
```

### 3. Infrastructure (infra)

A camada de infraestrutura contém implementações concretas de interfaces definidas no domínio e na aplicação, bem como todos os componentes que interagem com o mundo externo.

#### Responsabilidades:
- Implementar repositórios
- Integrar com serviços externos
- Fornecer adaptadores para frameworks
- Implementar detalhes técnicos
- Fornecer interfaces com o mundo externo (API, CLI)

#### Estrutura:
```
infra/
├── database/       # Implementações de banco de dados
├── repository/     # Implementações de repositórios
├── s3/        # Serviços de armazenamento
├── api/            # API HTTP
│   ├── handler/    # Handlers HTTP
│   └── router.go   # Configuração de rotas
```

Com DDD, preferir:

```
infra/
  api/handler/{bounded-context}/{aggregate}/...
  repository/{bounded-context}/{aggregate}/{tech}/...
```

#### Diretrizes:
- Pode depender das camadas de domínio e aplicação
- Deve implementar interfaces definidas no domínio e na aplicação
- Deve encapsular detalhes técnicos
- Deve ser substituível sem afetar as camadas internas



## Fluxo de Dependências

O fluxo de dependências segue a regra da dependência: as camadas internas não conhecem as camadas externas.

```
Infrastructure → Application → Domain
```

## Injeção de Dependências

A injeção de dependências é utilizada para fornecer implementações concretas para interfaces:

```go
// cmd/app/main.go
func main() {
    // Infraestrutura
    db := postgres.NewConnection()
    videoRepo := postgres.NewVideoRepository(db)
    ffmpegWrapper := ffmpeg.NewFFmpegWrapper()
    
    // Aplicação
    videoConverter := service.NewVideoConverter(ffmpegWrapper, videoRepo)
    uploadUseCase := usecase.NewUploadVideo(videoRepo, videoConverter)
    
    // API (parte da infraestrutura)
    videoHandler := handler.NewVideoHandler(uploadUseCase)
    router := api.NewRouter(videoHandler)
    
    // Iniciar servidor
    http.ListenAndServe(":8080", router)
}
```
