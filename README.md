# go-clean-architecture (Codex Skill)

Este repositório contém uma **skill do Codex** para criar serviços em **Go** com:

- Arquitetura em camadas (`domain` / `application` / `infra`)
- DDD com organização por `bounded-context/aggregate/action`
- SOLID + estratégia de testes
- Checks automatizados de dependências (camadas e bounded contexts)

## Onde está a skill

- Skill: `go-clean-architecture/`
- Documentação/guia: `go-clean-architecture/SKILL.md`
- Template de projeto (exemplo): `go-clean-architecture/assets/go-layered-service-template/`

## Uso rápido

Gerar um projeto a partir do template:

```bash
cd go-clean-architecture
./scripts/scaffold.sh --all example.com/minha-app /tmp/minha-app
```

Validar regras de arquitetura em um projeto:

```bash
./go-clean-architecture/scripts/check_arch.sh /caminho/do/projeto
```

