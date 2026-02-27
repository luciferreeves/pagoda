#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
GARDEN_DIR="$PROJECT_ROOT/garden"
DIST_DIR="$GARDEN_DIR/dist"
ENV_FILE="$PROJECT_ROOT/.env"
API_BASE="https://nekoweb.org/api"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()   { echo -e "${GREEN}[deploy]${NC} $1"; }
warn()  { echo -e "${YELLOW}[deploy]${NC} $1"; }
error() { echo -e "${RED}[deploy]${NC} $1" >&2; exit 1; }

if [ ! -f "$ENV_FILE" ]; then
    error ".env file not found at $ENV_FILE"
fi

source "$ENV_FILE"

if [ -z "${NEKOWEB_API_KEY:-}" ]; then
    error "NEKOWEB_API_KEY is not set in .env"
fi

if [ -z "${NEKOWEB_SITE:-}" ]; then
    error "NEKOWEB_SITE is not set in .env"
fi

if [ ! -d "$GARDEN_DIR" ]; then
    error "garden/ directory not found at $GARDEN_DIR"
fi

SITE_FOLDER="/${NEKOWEB_SITE}.nekoweb.org"

log "Building garden..."
cd "$GARDEN_DIR" && npm run build
cd "$PROJECT_ROOT"

if [ ! -d "$DIST_DIR" ]; then
    error "Build failed: dist/ not found"
fi

echo ""

api_post() {
    local endpoint="$1"
    shift
    curl -s -w "\n%{http_code}" \
        -H "Authorization: $NEKOWEB_API_KEY" \
        "$API_BASE$endpoint" \
        "$@"
}

api_get() {
    local endpoint="$1"
    curl -s \
        -H "Authorization: $NEKOWEB_API_KEY" \
        "$API_BASE$endpoint"
}

check_response() {
    local response="$1"
    local description="$2"
    local http_code
    http_code=$(echo "$response" | tail -n1)
    local body
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        log "$description - OK"
    else
        warn "$description - HTTP $http_code: $body"
        return 1
    fi
}

delete_remote_path() {
    local path="$1"
    local response
    response=$(api_post "/files/delete" -d "pathname=$path")
    check_response "$response" "Deleted $path"
}

create_remote_folder() {
    local path="$1"
    local response
    response=$(api_post "/files/create" \
        -d "pathname=$path" \
        -d "isFolder=true")
    local http_code
    http_code=$(echo "$response" | tail -n1)
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        log "Created folder: $path"
    fi
}

upload_file() {
    local local_path="$1"
    local remote_dir="$2"
    local response
    response=$(api_post "/files/upload" \
        -F "pathname=$remote_dir" \
        -F "files=@$local_path")
    check_response "$response" "Upload $(basename "$local_path") -> $remote_dir"
}

log "Starting deployment to nekoweb..."
log "Site: ${NEKOWEB_SITE}.nekoweb.org"
log "Source: $DIST_DIR"
echo ""

log "Cleaning root..."
root_items=$(api_get "/files/readfolder?pathname=/")
root_count=$(echo "$root_items" | jq 'length')

if [ "$root_count" -gt 0 ]; then
    for i in $(seq 0 $((root_count - 1))); do
        name=$(echo "$root_items" | jq -r ".[$i].name")
        delete_remote_path "/$name"
    done
else
    log "Root is already empty"
fi
echo ""

log "Creating site folder: $SITE_FOLDER"
create_remote_folder "$SITE_FOLDER"
echo ""

declare -a dirs_to_create=()
declare -a files_to_upload=()

while IFS= read -r -d '' file; do
    rel_path="${file#$DIST_DIR/}"

    if [[ "$rel_path" == .* ]]; then
        continue
    fi

    if [ -d "$file" ]; then
        dirs_to_create+=("$SITE_FOLDER/$rel_path")
    elif [ -f "$file" ]; then
        remote_dir="$SITE_FOLDER/$(dirname "$rel_path")"
        if [ "$remote_dir" = "$SITE_FOLDER/." ]; then
            remote_dir="$SITE_FOLDER"
        fi
        files_to_upload+=("$file|$remote_dir")
    fi
done < <(find "$DIST_DIR" -mindepth 1 -print0 | sort -z)

if [ ${#dirs_to_create[@]} -gt 0 ]; then
    log "Creating ${#dirs_to_create[@]} directories..."
    for dir in "${dirs_to_create[@]}"; do
        create_remote_folder "$dir"
    done
    echo ""
fi

file_count=${#files_to_upload[@]}
log "Uploading $file_count file(s)..."
echo ""

current=0
for entry in "${files_to_upload[@]}"; do
    local_path="${entry%%|*}"
    remote_dir="${entry##*|}"
    current=$((current + 1))
    log "[$current/$file_count] $(basename "$local_path")"
    upload_file "$local_path" "$remote_dir"
done

echo ""
log "Deployment complete!"
