# go-clean-architecture (Codex Skills)

Este repositório contém **skills do Codex** para criar serviços com:

- Arquitetura em camadas (`domain` / `application` / `infra`)
- DDD com organização por `bounded-context/aggregate/action`
- SOLID + estratégia de testes
- Checks automatizados de dependências (camadas e bounded contexts)

## Onde está a skill

- Go: `go-clean-architecture/` (doc: `go-clean-architecture/SKILL.md`, template: `go-clean-architecture/assets/go-layered-service-template/`)
- Node.js/TypeScript: `nodejs-clean-architecture/` (doc: `nodejs-clean-architecture/SKILL.md`, template: `nodejs-clean-architecture/assets/ts-layered-service-template/`)

## Uso rápido

Gerar um projeto a partir do template:

```bash
cd go-clean-architecture
./scripts/scaffold.sh --all example.com/minha-app /tmp/minha-app
```

```bash
cd nodejs-clean-architecture
./scripts/scaffold.sh --all @acme/minha-app /tmp/minha-app
```

Validar regras de arquitetura em um projeto:

```bash
./go-clean-architecture/scripts/check_arch.sh /caminho/do/projeto
```

```bash
./nodejs-clean-architecture/scripts/check_arch.sh /caminho/do/projeto
```
