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

check "internal/domain/**/*.go" "internal/application" "domain importing application"
check "internal/domain/**/*.go" "internal/infra" "domain importing infra"
check "internal/application/**/*.go" "internal/infra" "application importing infra"

domain_external="$(
  (cd "$project_root" && rg -n --glob "internal/domain/**/*.go" --glob "!**/*_test.go" "\"github.com/" .) \
    | rg -v "github.com/asaskevich/govalidator" \
    || true
)"
if [[ -n "$domain_external" ]]; then
  echo "FAIL: domain production code importing external packages (allowed: github.com/asaskevich/govalidator)" >&2
  echo "$domain_external" >&2
  fail=1
fi

if (cd "$project_root" && rg -n --glob "internal/application/**/*.go" --glob "!**/*_test.go" "\"database/sql\"" . >/dev/null); then
  echo "FAIL: application production code importing database/sql (should be infra)" >&2
  (cd "$project_root" && rg -n --glob "internal/application/**/*.go" --glob "!**/*_test.go" "\"database/sql\"" . >&2) || true
  fail=1
fi

if (cd "$project_root" && rg -n --glob "internal/application/**/*.go" --glob "!**/*_test.go" "\"github.com/go-chi/chi" . >/dev/null); then
  echo "FAIL: application production code importing chi (should be infra/api)" >&2
  (cd "$project_root" && rg -n --glob "internal/application/**/*.go" --glob "!**/*_test.go" "\"github.com/go-chi/chi" . >&2) || true
  fail=1
fi

if (cd "$project_root" && rg -n --glob "internal/application/**/*.go" --glob "!**/*_test.go" "\"github.com/aws/aws-sdk-go-v2" . >/dev/null); then
  echo "FAIL: application production code importing aws sdk (should be infra)" >&2
  (cd "$project_root" && rg -n --glob "internal/application/**/*.go" --glob "!**/*_test.go" "\"github.com/aws/aws-sdk-go-v2" . >&2) || true
  fail=1
fi

# DDD-aware check: prevent cross-bounded-context imports by default.
entity_root="$project_root/internal/domain/entity"
if [[ -d "$entity_root" ]]; then
  contexts=()
  while IFS= read -r dir; do
    [[ -z "$dir" ]] && continue
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

        entity_pattern="\"[^\"]*/internal/domain/entity/${other}/"
        repo_pattern="\"[^\"]*/internal/domain/repository/${other}/"

        if (cd "$project_root" && rg -n \
          --glob "internal/domain/entity/${ctx}/**/*.go" \
          --glob "internal/domain/repository/${ctx}/**/*.go" \
          --glob "internal/application/usecase/${ctx}/**/*.go" \
          --glob "internal/application/service/${ctx}/**/*.go" \
          --glob "internal/infra/api/handler/${ctx}/**/*.go" \
          --glob "internal/infra/repository/${ctx}/**/*.go" \
          -e "$entity_pattern" -e "$repo_pattern" \
          . >/dev/null); then
          echo "FAIL: bounded context '$ctx' importing '$other' (configure allowlist at $allowed_file)" >&2
          (cd "$project_root" && rg -n \
            --glob "internal/domain/entity/${ctx}/**/*.go" \
            --glob "internal/domain/repository/${ctx}/**/*.go" \
            --glob "internal/application/usecase/${ctx}/**/*.go" \
            --glob "internal/application/service/${ctx}/**/*.go" \
            --glob "internal/infra/api/handler/${ctx}/**/*.go" \
            --glob "internal/infra/repository/${ctx}/**/*.go" \
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
