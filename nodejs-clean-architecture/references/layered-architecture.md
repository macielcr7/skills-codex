# Arquitetura em camadas (Node.js + TypeScript)

## Objetivo

Separar regras de negócio de detalhes técnicos para manter:
- baixo acoplamento
- alta testabilidade
- evolução segura (trocar DB/HTTP sem reescrever domínio)

## Camadas e responsabilidades

### `domain/`

- Regras do domínio (entidades, invariantes, erros do domínio)
- Interfaces de repositório (contratos)
- Nada de framework (NestJS), ORM (Prisma) ou IO

### `application/`

- Casos de uso (orquestração): *o que o sistema faz*
- Depende de contratos do `domain`
- Define portas (interfaces) para dependências técnicas (ex.: `IdGenerator`, `Clock`, `Uploader`)

### `infrastructure/`

- Detalhes técnicos: NestJS (HTTP), Prisma (DB), integrações externas
- Implementa interfaces de `domain` e `application`
- Converte erros/DTOs para o mundo externo (ex.: HTTP 400/404)

### `src/main.ts` e `src/app.module.ts` (composition root)

- Wiring/DI: escolhe implementações concretas e injeta no app
- Bootstrap do processo (ex.: start do servidor)

## Regra de dependência

Fluxo permitido:

```
infrastructure -> application -> domain
main/app.module -> (infrastructure, application, domain)
```

Proibido:
- `domain` importar `application/infrastructure`
- `application` importar `infrastructure`

## DDD por Bounded Context / Entity / Action

O projeto replica o bounded context em todas as camadas:

- `domain/entity/<context>/<entity>/...`
- `domain/repository/<context>/<entity>/...`
- `application/usecase/<context>/<entity>/<action>/...`
- `infrastructure/nest/<context>/controllers/<entity>/<action>.controller.ts`
- `infrastructure/nest/<context>/controllers/<entity>/<action>.schema.ts`
- `infrastructure/repository/<context>/<entity>/<tech>/...`

Quando precisar compartilhar utilitários, use `shared/` dentro da camada (evite “importar domínio de outro contexto”).

## Regras enforçadas

O script `scripts/check_arch.sh` valida:
- dependências entre camadas
- limites de bounded context (com allowlist via `.allowed_contexts`)
