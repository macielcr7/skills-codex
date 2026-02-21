# Convenções (TypeScript + ESM)

## Módulos

- Projeto em **ESM** (`"type": "module"`).
- `tsconfig` com `module/moduleResolution: NodeNext`.

## Imports

- Preferir alias interno `#/*` (mapeado para `src/*` no `tsconfig`).
- Evitar imports relativos longos (melhora legibilidade e facilita `check_arch.sh`).

## Infra Nest (organização)

- Em `src/infrastructure/nest/<modulo>/`, usar subpastas para escalar:
  - `controllers/`, `services/`, `guards/`, `decorators/`, `gateways/` (externo), `__tests__/` (integração)
  - manter na raiz do módulo apenas: `<modulo>.module.ts`, `<modulo>.tokens.ts`, `<modulo>.types.ts` (quando necessário)
- Itens cross-cutting (ex.: exception filters globais) ficam em `src/infrastructure/nest/common/**`.

## HTTP (parsing / Zod) — obrigatório

- Não declarar schemas Zod inline em `*.controller.ts`.
- Criar arquivo adjacente `*.schema.ts` no mesmo diretório do controller e importar de lá.
- Preferir reutilização via:
  - `src/infrastructure/nest/common/http/**` (schemas/helpers compartilhados)
  - `src/infrastructure/nest/<modulo>/shared/**` quando fizer sentido

## NestJS + Vitest (DI)

- Em testes rodando TS via Vitest, nem sempre há `emitDecoratorMetadata`; quando necessário, prefira injeção explícita com `@Inject(Token)` nos construtores (especialmente em controllers).

## Naming e arquivos

- Pastas e arquivos de **usecase/controller/action** em `snake_case`:
  - `upload_video.ts`, `get_video.ts`
- Use cases seguem o padrão `bounded-context/entity/action/*`, com múltiplos arquivos por action quando necessário.
- Se vários use cases do mesmo `{entity}` compartilharem tipos/helpers, crie `src/application/usecase/{bounded-context}/{entity}/shared/**`.
- Types de use case devem ser nomeados por action (evitar `Input`/`Output` genéricos):
  - `GetVideoInput`, `GetVideoOutput`
  - `UploadVideoInput`, `UploadVideoOutput`
- Sufixos:
  - Interfaces: `*.interface.ts`
  - Types: `*.type.ts`
  - Repository interfaces: `*.repository.interface.ts`
  - Repository implementations: `*.repository.ts`
  - Controllers: `*.controller.ts`
- Classes/Tipos em `PascalCase`.
- 1 use case por arquivo; 1 controller por action.

## Erros

- Domínio: erros do domínio (ex.: `VideoNotFoundError`, `EntityValidationError`).
- Infrastructure: mapeia erros para HTTP (400/404/500).

## Testes

- Unit tests próximos ao código (`*.test.ts`) em `domain` e `application`.
- E2E em `test/e2e/**/*.test.ts`, subindo o Nest app e chamando endpoints.

Regra do projeto:
- Todo arquivo `.ts` em `src/` que contenha lógica (classes/funções/branches/parsing/validações) deve ter um teste correspondente.
- Exceções permitidas:
  - `*.module.ts`, `*.tokens.ts`, `*.types.ts`
  - `src/infrastructure/database/schema/**`
  - `src/infrastructure/database/database.tokens.ts`
  - `src/main.ts`, `src/app.module.ts`
  - arquivos “type-only” (apenas `type/interface`, sem runtime)
