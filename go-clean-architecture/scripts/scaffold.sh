#!/usr/bin/env bash
set -euo pipefail

run_tidy=0
run_tests=0
run_check_arch=0

usage() {
  cat >&2 <<'EOF'
usage: scaffold.sh [--tidy] [--test] [--check-arch] [--all] <module_path> <output_dir>

options:
  --tidy        run 'go mod tidy' in the new project
  --test        run 'go test ./...' in the new project
  --check-arch  run the layer dependency checker
  --all         run --check-arch --tidy --test
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --tidy) run_tidy=1; shift ;;
    --test) run_tests=1; shift ;;
    --check-arch) run_check_arch=1; shift ;;
    --all) run_check_arch=1; run_tidy=1; run_tests=1; shift ;;
    -h|--help) usage; exit 0 ;;
    --) shift; break ;;
    -*) echo "unknown option: $1" >&2; usage; exit 2 ;;
    *) break ;;
  esac
done

if [[ $# -lt 2 ]]; then
  usage
  exit 2
fi

module_path="$1"
output_dir="$2"

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
skill_root="$(cd "$script_dir/.." && pwd)"
template_dir="$skill_root/assets/go-layered-service-template"

if [[ ! -d "$template_dir" ]]; then
  echo "template not found: $template_dir" >&2
  exit 1
fi

mkdir -p "$output_dir"
if [[ -n "$(ls -A "$output_dir" 2>/dev/null || true)" ]]; then
  echo "output_dir must be empty: $output_dir" >&2
  exit 1
fi

cp -R "$template_dir/." "$output_dir"

mkdir -p "$output_dir/scripts"
cp "$skill_root/scripts/check_arch.sh" "$output_dir/scripts/check_arch.sh"
chmod +x "$output_dir/scripts/check_arch.sh"

if ! command -v perl >/dev/null 2>&1; then
  echo "perl is required" >&2
  exit 2
fi

perl -pi -e "s#\\bexample\\.com/service\\b#$module_path#g" "$output_dir/go.mod"

if command -v rg >/dev/null 2>&1; then
  while IFS= read -r file; do
    perl -pi -e "s#\\bexample\\.com/service\\b#$module_path#g" "$file"
  done < <(rg -l --glob '*.go' 'example\.com/service' "$output_dir" || true)
else
  while IFS= read -r -d '' file; do
    if grep -q 'example.com/service' "$file"; then
      perl -pi -e "s#\\bexample\\.com/service\\b#$module_path#g" "$file"
    fi
  done < <(find "$output_dir" -type f -name '*.go' -print0)
fi

echo "scaffolded at: $output_dir"
echo "next:"
echo "  cd \"$output_dir\""

if [[ "$run_check_arch" -eq 1 ]]; then
  (cd "$output_dir" && ./scripts/check_arch.sh .)
else
  echo "  ./scripts/check_arch.sh .    # optional"
fi

go_env=()
if [[ -z "${GOCACHE:-}" ]]; then
  mkdir -p "$output_dir/.cache/go-build"
  go_env+=( "GOCACHE=$output_dir/.cache/go-build" )
fi

if [[ "$run_tidy" -eq 1 ]]; then
  (cd "$output_dir" && env "${go_env[@]}" go mod tidy)
else
  echo "  go mod tidy                  # optional"
fi

if [[ "$run_tests" -eq 1 ]]; then
  (cd "$output_dir" && env "${go_env[@]}" go test ./...)
else
  echo "  go test ./..."
fi
