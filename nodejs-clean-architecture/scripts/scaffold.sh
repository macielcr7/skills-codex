#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  scaffold.sh [--pm npm|pnpm|yarn] [--install] [--test] [--check-arch] [--all] <package_name> <output_dir>

Examples:
  ./scripts/scaffold.sh @acme/media-service ./tmp/media-service
  ./scripts/scaffold.sh --all @acme/media-service ./tmp/media-service
  ./scripts/scaffold.sh --pm pnpm --install @acme/media-service ./tmp/media-service
EOF
}

pm="npm"
run_install=0
run_test=0
run_check_arch=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      usage
      exit 0
      ;;
    --pm)
      pm="${2:-}"
      shift 2
      ;;
    --install)
      run_install=1
      shift 1
      ;;
    --test)
      run_test=1
      shift 1
      ;;
    --check-arch)
      run_check_arch=1
      shift 1
      ;;
    --all)
      run_install=1
      run_test=1
      run_check_arch=1
      shift 1
      ;;
    --*)
      echo "Unknown flag: $1" >&2
      usage >&2
      exit 2
      ;;
    *)
      break
      ;;
  esac
done

package_name="${1:-}"
output_dir="${2:-}"

if [[ -z "$package_name" || -z "$output_dir" ]]; then
  usage >&2
  exit 2
fi

if [[ "$pm" != "npm" && "$pm" != "pnpm" && "$pm" != "yarn" ]]; then
  echo "Unsupported --pm: $pm (use npm|pnpm|yarn)" >&2
  exit 2
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
skill_root="$(cd "$script_dir/.." && pwd)"
template_dir="$skill_root/assets/ts-layered-service-template"

if [[ ! -d "$template_dir" ]]; then
  echo "template not found: $template_dir" >&2
  exit 2
fi

mkdir -p "$output_dir"
cp -R "$template_dir/." "$output_dir/"

# Update package.json name
if [[ -f "$output_dir/package.json" ]]; then
  perl -pi -e "s/\"name\"\\s*:\\s*\"[^\"]*\"/\"name\": \"$package_name\"/g" "$output_dir/package.json"
fi

# Ensure architecture check script exists (overwrite with skill version)
mkdir -p "$output_dir/scripts"
cp "$skill_root/scripts/check_arch.sh" "$output_dir/scripts/check_arch.sh"
chmod +x "$output_dir/scripts/check_arch.sh"

echo "OK: scaffolded at $output_dir"

if [[ "$run_install" -eq 1 ]]; then
  echo "Running install ($pm)..."
  (cd "$output_dir" && "$pm" install)
fi

if [[ "$run_test" -eq 1 ]]; then
  echo "Running tests ($pm)..."
  (cd "$output_dir" && "$pm" test)
fi

if [[ "$run_check_arch" -eq 1 ]]; then
  echo "Running architecture checks..."
  (cd "$output_dir" && bash scripts/check_arch.sh .)
fi

