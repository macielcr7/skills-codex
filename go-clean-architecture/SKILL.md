---
name: go-clean-architecture
description: Projetar, implementar e refatorar serviços em Go seguindo a arquitetura em camadas (domain/application/infra + wiring em cmd/), aplicando SOLID e testes (unitário/integração/e2e) e adotando as libs recomendadas (chi, slog, database/sql, migrate, aws-sdk-go-v2, godotenv, uuid, testify). Use quando precisar estruturar um serviço Go, definir responsabilidades por camada, criar casos de uso, repositórios, handlers HTTP e integrações (DB/S3), ou montar pipeline (build/test/container).
---

# Arquitetura em camadas para Go

## Fonte de verdade (carregar antes de decidir)

- `references/layered-architecture.md`: camadas, responsabilidades e regra de dependência.
- `references/go-libs.md`: libs recomendadas por categoria.
- `references/conventions-go.md`: convenções (naming, imports, testes, docs).

## Estrutura padrão

```
cmd/
  app/
    main.go                 # wiring/DI + bootstrap
internal/
  domain/
    entity/
      {bounded-context}/
        {aggregate}/
    repository/
      {bounded-context}/
        {aggregate}/
  application/
    service/
      {bounded-context}/
      shared/
    usecase/
      {bounded-context}/
        {aggregate}/
          {action}.go
  infra/
    api/
      handler/
        {bounded-context}/
          {aggregate}/
    repository/
      {bounded-context}/
        {aggregate}/
          {tech}/
    database/
      shared/
    s3/
      shared/
    id/
      shared/
```

## DDD (como trabalhar com Bounded Contexts)

- Definir Bounded Contexts como “módulos de domínio” com linguagem ubíqua própria (ex.: `media`, `billing`, `catalog`).
- Manter **entidades e regras** dentro do contexto:
  - `internal/domain/entity/{bounded-context}/{entity}/...`
  - Ex.: `internal/domain/entity/media/video/video.go` (`video` é o agregado, `package video`)
- Definir repositórios por contexto:
  - `internal/domain/repository/{bounded-context}/...`
  - Ex.: `internal/domain/repository/media/video/video_repository.go`

- Replicar o contexto nas demais camadas (regra do projeto):
  - Use cases por contexto/agregado/ação: `internal/application/usecase/{bounded-context}/{aggregate}/{action}.go`
  - Serviços de aplicação por contexto (quando fizer sentido): `internal/application/service/{bounded-context}/...`
  - Adapters HTTP por contexto/agregado: `internal/infra/api/handler/{bounded-context}/{aggregate}/...`
  - Repositórios/clients por contexto/agregado/tech: `internal/infra/repository/{bounded-context}/{aggregate}/{tech}/...`
  - Cross-cutting: usar `shared/` (ex.: `internal/application/service/shared/`), evitando acoplamento entre contextos.

Quando o Bounded Context crescer:
- Preferir `bounded-context/{aggregate}/{action}` em vez de “um pacote gigante por contexto”.
- Se necessário, agrupar por tipo:
  - `internal/application/usecase/{bounded-context}/{aggregate}/commands/...`
  - `internal/application/usecase/{bounded-context}/{aggregate}/queries/...`
  - (aplicar o mesmo raciocínio em `infra/api/handler/...` e `infra/repository/...`)

Heurísticas:
- Um agregado (aggregate root) tende a virar 1 “entity package” (ex.: `video`) com invariantes e métodos.
- Value Objects podem viver no mesmo package do agregado ou em subpackages do contexto (quando compartilhados).
- Evitar vazar tipos do contexto para outros contextos: comunicar via IDs/DTOs no `application` quando necessário.

## Regras de dependência (obrigatório)

- `internal/domain`:
  - não importar `internal/application` nem `internal/infra`
  - manter regras puras (preferir só stdlib)
- `internal/application`:
  - importar `internal/domain`
  - não importar `internal/infra`
  - declarar interfaces para integrações técnicas (ex.: `FFmpegWrapper`, `Uploader`, `IDGenerator`) e receber implementações via DI
- `internal/infra`:
  - implementar interfaces do `domain` e do `application`
  - conter HTTP/DB/S3/frameworks/SDKs
- `cmd/`:
  - fazer composição (wiring), nunca regra de negócio

Fluxo: `infra → application → domain`.

## Convenções de código (para manter SOLID)

- Seguir `references/conventions-go.md` (snake_case, imports por grupo, testes no mesmo pacote, comentários em exports).
- SRP: 1 use case por arquivo; 1 handler por recurso; 1 repo por agregado.
- ISP: interfaces pequenas (ex.: separar `VideoFinder` de `VideoSaver` quando fizer sentido).
- DIP: use cases dependem de interfaces; infra depende das camadas internas.

Padrões:
- Use case: `Execute(ctx, input) (output, error)`.
- Domínio: invariantes no construtor/métodos; erros sentinela no domínio.
- HTTP: validação/DTO/mapeamento de erro ficam no handler (infra).

## Regras DDD (enforçadas pelo script)

- Por padrão, um Bounded Context **não importa** entidades/repositórios de outro contexto diretamente.
- Configurar exceções por contexto criando: `internal/domain/entity/{bounded-context}/.allowed_contexts` (1 contexto por linha).

## Libs recomendadas (padrão do projeto)

- Logging: `log/slog`
- REST: `github.com/go-chi/chi/v5`
- DB: `database/sql` + migrations `github.com/golang-migrate/migrate/v4`
- S3: `github.com/aws/aws-sdk-go-v2` + LocalStack `github.com/localstack/localstack-go`
- Env local: `github.com/joho/godotenv`
- UUID: `github.com/google/uuid`
- Validação: `github.com/asaskevich/govalidator`
- Testes: `github.com/stretchr/testify`

## Estratégia de testes

- `domain`: unit tests sem mocks.
- `application`: unit tests com fakes/mocks de repositórios/serviços.
- `infra`:
  - HTTP: `net/http/httptest`
  - DB/S3: testes de integração (LocalStack para S3)

## Template e scripts desta skill

- Gerar serviço novo: `./scripts/scaffold.sh <module_path> <output_dir>` (ou `--all` para `check-arch + tidy + test`)
- Validar camadas: `./scripts/check_arch.sh <project_root>` (default `.`)

O template fica em `assets/go-layered-service-template/` e já contém:
- estrutura `cmd/` + `internal/` em camadas
- exemplo DDD com Bounded Context `media`:
  - entidade: `internal/domain/entity/media/video/`
  - repositório: `internal/domain/repository/media/video/`
  - use cases: `internal/application/usecase/media/video/`
  - handlers HTTP: `internal/infra/api/handler/media/video/`
  - repo infra (memory): `internal/infra/repository/media/video/memory/`
- stubs de DB (`database/sql` + `cmd/migrate` com `golang-migrate`) e S3 (AWS SDK v2)
- testes unitários e de handler (com `testify` e `httptest`)

## Checklist mínimo (para PR “pronto pra prod”)

- `go test ./...` passando.
- build de binário e imagem Docker (multi-stage).
- health endpoint (`GET /health`).
- logs estruturados (`slog`) com request-id (se aplicável).
- config via env (com `.env` somente para dev).
