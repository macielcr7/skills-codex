---
name: nodejs-clean-architecture
description: Projetar, implementar e refatorar serviços em Node.js (TypeScript) com NestJS seguindo Clean Architecture em camadas (domain/application/infrastructure + wiring em src/main.ts e src/app.module.ts), aplicando SOLID, DDD (bounded-context/entity/action), validação com Zod (entities com validate()) e testes com Vitest. Inclui template com Prisma (infrastructure) e scripts de scaffold + check de arquitetura.
---

# Clean Architecture (Node.js + TypeScript + NestJS)

## Fonte de verdade (carregar antes de decidir)

- `references/layered-architecture.md`: camadas, responsabilidades e regra de dependência.
- `references/node-libs.md`: libs recomendadas por categoria (NestJS, Zod, Prisma, Vitest).
- `references/conventions-typescript.md`: convenções (ESM, imports, naming, testes).

## Estrutura padrão (DDD + camadas)

```
src/
  main.ts                       # bootstrap (composition root)
  app.module.ts                 # wiring/DI (composition root)

  domain/
    entity/
      {bounded-context}/
        {entity}/
          ...
    repository/
      {bounded-context}/
        {entity}/
          ...
    shared/
      ...

  application/
    service/
      {bounded-context}/
      shared/
        ...
    usecase/
      {bounded-context}/
        {entity}/
          shared/
            ...
          {action}/
            ...

  infrastructure/
    database/
      database.tokens.ts
      schema/
        ...
    nest/
      common/
        common.module.ts
        controllers/
        services/
      {bounded-context}/
        {bounded-context}.module.ts
        {bounded-context}.tokens.ts
        controllers/
          ...
          {entity}/
            {action}.controller.ts
            {action}.schema.ts
    repository/
      {bounded-context}/
        {entity}/
          {tech}/
            ...
```

## Convenções de arquivos (obrigatório)

- Interfaces: `*.interface.ts`
- Types: `*.type.ts`
- Types de use case devem ser nomeados por action (evitar `Input`/`Output` genéricos):
  - `GetVideoInput`, `GetVideoOutput`
- Repository interfaces: `*.repository.interface.ts`
- Repository implementations: `*.repository.ts`
- Controllers: `*.controller.ts`

## DDD (como trabalhar com Bounded Contexts)

- Defina Bounded Contexts como “módulos de domínio” (ex.: `media`, `billing`, `catalog`).
- Replique o contexto em **todas** as camadas (regra do projeto):
  - Entidades: `src/domain/entity/{bounded-context}/{entity}/...`
  - Repositórios (interfaces): `src/domain/repository/{bounded-context}/{entity}/...`
  - Use cases: `src/application/usecase/{bounded-context}/{entity}/{action}/...`
  - Controllers HTTP (Nest): `src/infrastructure/nest/{bounded-context}/controllers/{entity}/{action}.controller.ts`
  - Schemas HTTP (Zod): `src/infrastructure/nest/{bounded-context}/controllers/{entity}/{action}.schema.ts`
  - Repositórios (infrastructure): `src/infrastructure/repository/{bounded-context}/{entity}/{tech}/...`
- Se um contexto crescer, prefira `bounded-context/{entity}/{action}` (evita “pasta gigante por contexto”).
- Se vários use cases do mesmo `{entity}` compartilharem tipos/helpers, crie `src/application/usecase/{bounded-context}/{entity}/shared/**`.
- Cross-cutting: use `shared/` (em cada camada) em vez de acoplamento entre contextos.

## Regras de dependência (obrigatório)

- `src/domain`:
  - não importar `src/application` nem `src/infrastructure`
  - manter regras do domínio (entidades com `validate()` usando **Zod**)
- `src/application`:
  - importar `src/domain`
  - não importar `src/infrastructure`
  - declarar interfaces para integrações técnicas (ex.: `IdGenerator`, `Uploader`) e receber implementações via DI
- `src/infrastructure`:
  - implementar interfaces de `domain` e `application`
  - conter NestJS (controllers), Prisma, integrações, etc.
- `src/main.ts` e `src/app.module.ts`:
  - composição (wiring/DI), nunca regra de negócio

Fluxo: `infrastructure → application → domain`.

## Regras DDD (enforçadas pelo script)

- Por padrão, um Bounded Context **não importa** entidades/repositórios de outro contexto diretamente.
- Configure exceções por contexto criando: `src/domain/entity/{bounded-context}/.allowed_contexts` (1 contexto por linha).

## Libs recomendadas (padrão do template)

- HTTP: NestJS
- Validação: Zod (entities com `validate()` + testes)
- ORM: Prisma (somente em `infrastructure`)
- Testes: Vitest

## Estratégia de testes

- `domain`: unit tests sem dependências de infrastructure.
- `application`: unit tests com fakes (não importar `infrastructure`).
- `infrastructure`: unit/integration tests para controllers/repos + e2e HTTP.

Regra do projeto:
- Todo arquivo `.ts` em `src/` que contenha lógica deve ter teste correspondente.
- Exceções: `*.module.ts`, `*.tokens.ts`, `*.types.ts`, `src/main.ts`, `src/app.module.ts`, arquivos type-only.

## Template e scripts desta skill

- Gerar serviço novo: `./scripts/scaffold.sh <package_name> <output_dir>` (use `--all` para `install + test + check-arch`)
- Validar camadas: `./scripts/check_arch.sh <project_root>` (default `.`)

O template fica em `assets/ts-layered-service-template/` e inclui exemplo DDD com Bounded Context `media`:
- entidade: `src/domain/entity/media/video/`
- repositório (interface): `src/domain/repository/media/video/video.repository.interface.ts`
- erro do domínio: `src/domain/repository/media/video/video_not_found.error.ts`
- use cases: `src/application/usecase/media/video/`
- módulo Nest (media): `src/infrastructure/nest/media/media.module.ts`
- controllers HTTP (Nest): `src/infrastructure/nest/media/controllers/video/`
- repo (memory): `src/infrastructure/repository/media/video/memory/`
- repo (prisma): `src/infrastructure/repository/media/video/prisma/` (opcional)
