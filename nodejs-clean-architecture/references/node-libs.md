# Libs recomendadas (Node.js + TypeScript)

## HTTP

- **NestJS** (`@nestjs/*`): framework para API, DI e módulos.

## Validação

- **Zod** (`zod`): schemas tipados e validação usada em:
  - entities do `domain` (via `validate()`)
  - controllers do `infrastructure/nest` (validação de request/params)

## ORM / DB

- **Prisma** (`prisma` + `@prisma/client`): ORM moderno e tipado.
  - Fica **somente em `infrastructure`**
  - Template mantém repo `memory` como default para testes não dependerem de DB

## Testes

- **Vitest** (`vitest`): unit/integration/e2e com rapidez e DX bom.

## Utilitários

- UUID: `node:crypto` (`randomUUID`) via adapter em `infrastructure/nest/common/services`
