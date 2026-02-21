#!/usr/bin/env bash
set -euo pipefail

project_root="${1:-.}"

if [[ ! -d "$project_root" ]]; then
  echo "project_root not found: $project_root" >&2
  exit 2
fi

if ! command -v rg >/dev/null 2>&1; then
  echo "ripgrep (rg) is required" >&2
  exit 2
fi

fail=0

read_allowed_contexts() {
  local file="$1"

  if [[ ! -f "$file" ]]; then
    return 0
  fi

  while IFS= read -r line; do
    line="${line%%#*}"
    line="${line//$'\r'/}"
    line="$(echo "$line" | xargs || true)"
    if [[ -z "$line" ]]; then
      continue
    fi
    echo "$line"
  done <"$file"
}

contains() {
  local needle="$1"
  shift
  local item
  for item in "$@"; do
    if [[ "$item" == "$needle" ]]; then
      return 0
    fi
  done
  return 1
}

check() {
  local glob="$1"
  local pattern="$2"
  local message="$3"

  if (cd "$project_root" && rg -n --glob "$glob" "$pattern" . >/dev/null); then
    echo "FAIL: $message" >&2
    (cd "$project_root" && rg -n --glob "$glob" "$pattern" . >&2) || true
    fail=1
  fi
}

# Layering rules (imports via #/* alias).
check "src/domain/**/*.ts" "['\"]#/(application|infrastructure)/" "domain importing application/infrastructure"
check "src/application/**/*.ts" "['\"]#/infrastructure/" "application importing infrastructure"

# HTTP/Zod rule: controllers must not declare Zod schemas inline.
if (cd "$project_root" && rg -n --glob "src/infrastructure/nest/**/controllers/**/*.controller.ts" --glob "!**/*.test.ts" "from ['\"]zod['\"]" . >/dev/null); then
  echo "FAIL: controllers importing zod (schemas must live in adjacent *.schema.ts files)" >&2
  (cd "$project_root" && rg -n --glob "src/infrastructure/nest/**/controllers/**/*.controller.ts" --glob "!**/*.test.ts" "from ['\"]zod['\"]" . >&2) || true
  fail=1
fi

# Tech guardrails: keep NestJS/Prisma out of domain and application.
domain_tech="$(
  (cd "$project_root" && rg -n --glob "src/domain/**/*.ts" --glob "!**/*.test.ts" "from ['\"]@nestjs/|from ['\"]@prisma/client|from ['\"]prisma" .) \
    || true
)"
if [[ -n "$domain_tech" ]]; then
  echo "FAIL: domain production code importing framework/orm packages" >&2
  echo "$domain_tech" >&2
  fail=1
fi

application_tech="$(
  (cd "$project_root" && rg -n --glob "src/application/**/*.ts" --glob "!**/*.test.ts" "from ['\"]@nestjs/|from ['\"]@prisma/client|from ['\"]prisma" .) \
    || true
)"
if [[ -n "$application_tech" ]]; then
  echo "FAIL: application production code importing framework/orm packages" >&2
  echo "$application_tech" >&2
  fail=1
fi

# DDD-aware check: prevent cross-bounded-context imports by default.
entity_root="$project_root/src/domain/entity"
if [[ -d "$entity_root" ]]; then
  contexts=()
  while IFS= read -r dir; do
    [[ -z "$dir" ]] && continue
    ctx="$(basename "$dir")"
    if [[ "$ctx" == "shared" ]]; then
      continue
    fi
    contexts+=("$dir")
  done < <(find "$entity_root" -mindepth 1 -maxdepth 1 -type d -print 2>/dev/null | sort)

  if [[ "${#contexts[@]}" -gt 1 ]]; then
    for ctx_dir in "${contexts[@]}"; do
      ctx="$(basename "$ctx_dir")"
      allowed_file="$entity_root/$ctx/.allowed_contexts"

      allowed_contexts=()
      while IFS= read -r allowed; do
        allowed_contexts+=("$allowed")
      done < <(read_allowed_contexts "$allowed_file")

      for other_dir in "${contexts[@]}"; do
        other="$(basename "$other_dir")"
        if [[ "$other" == "$ctx" ]]; then
          continue
        fi
        if [[ "${#allowed_contexts[@]}" -gt 0 ]] && contains "$other" "${allowed_contexts[@]}"; then
          continue
        fi

        entity_pattern="['\"]#/domain/entity/${other}/"
        repo_pattern="['\"]#/domain/repository/${other}/"

        if (cd "$project_root" && rg -n \
          --glob "src/domain/entity/${ctx}/**/*.ts" \
          --glob "src/domain/repository/${ctx}/**/*.ts" \
          --glob "src/application/usecase/${ctx}/**/*.ts" \
          --glob "src/application/service/${ctx}/**/*.ts" \
          --glob "src/infrastructure/nest/${ctx}/**/*.ts" \
          --glob "src/infrastructure/repository/${ctx}/**/*.ts" \
          -e "$entity_pattern" -e "$repo_pattern" \
          . >/dev/null); then
          echo "FAIL: bounded context '$ctx' importing '$other' (configure allowlist at $allowed_file)" >&2
          (cd "$project_root" && rg -n \
            --glob "src/domain/entity/${ctx}/**/*.ts" \
            --glob "src/domain/repository/${ctx}/**/*.ts" \
            --glob "src/application/usecase/${ctx}/**/*.ts" \
            --glob "src/application/service/${ctx}/**/*.ts" \
            --glob "src/infrastructure/nest/${ctx}/**/*.ts" \
            --glob "src/infrastructure/repository/${ctx}/**/*.ts" \
            -e "$entity_pattern" -e "$repo_pattern" \
            . >&2) || true
          fail=1
        fi
      done
    done
  fi
fi

if [[ "$fail" -ne 0 ]]; then
  exit 1
fi

echo "OK: layer dependency rule satisfied"
