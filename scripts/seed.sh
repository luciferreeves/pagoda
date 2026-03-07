#!/bin/bash
set -euo pipefail

DB_PATH="shrine/pagoda.db"

if [ ! -f "$DB_PATH" ]; then
  echo "Database not found at $DB_PATH"
  exit 1
fi

HASH=$(npx -y bcryptjs-cli password 2>/dev/null)
OWNER_DATE=$(date -v-4m +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -d "4 months ago" +"%Y-%m-%dT%H:%M:%SZ")

echo "Generating 99 unique citizens..."
npx -y @faker-js/cli firstName >/dev/null 2>&1
FAKER_MODULES=$(find ~/.npm/_npx -path "*/@faker-js/cli/bin/faker.js" -print -quit 2>/dev/null)
FAKER_MODULES="${FAKER_MODULES%/@faker-js/cli/bin/faker.js}"

NAMES=$(NODE_PATH="$FAKER_MODULES" node -e "
var f = require('@faker-js/faker').faker;
var now = Date.now();
var threeMonthsAgo = new Date(now);
threeMonthsAgo.setMonth(threeMonthsAgo.getMonth() - 3);
var range = now - threeMonthsAgo.getTime();
var seen = new Set();
while (seen.size < 99) {
  var first = f.person.firstName();
  var username = first.toLowerCase().replace(/[^a-z]/g, '');
  if (username.length >= 3 && username.length <= 32 && !seen.has(username)) {
    seen.add(username);
    var last = f.person.lastName();
    var joinDate = new Date(threeMonthsAgo.getTime() + Math.random() * range).toISOString().replace(/\.\d{3}Z/, 'Z');
    console.log(username + '|' + first + ' ' + last + '|' + joinDate);
  }
}
")

SQL_FILE=$(mktemp)
trap "rm -f $SQL_FILE" EXIT

printf '%s' "INSERT OR IGNORE INTO users (username, email, password_hash, display_name, role, email_verified, created_at, updated_at) VALUES " > "$SQL_FILE"
printf '%s' "('cr', 'cr@pagoda.local', '${HASH}', 'Bobby', 'owner', 1, '${OWNER_DATE}', '${OWNER_DATE}')" >> "$SQL_FILE"

COUNT=0

while IFS='|' read -r USERNAME DISPLAY JOIN_DATE; do
  DISPLAY=$(echo "$DISPLAY" | sed "s/'/''/g")

  printf '%s' ", ('${USERNAME}', '${USERNAME}@pagoda.local', '${HASH}', '${DISPLAY}', 'member', 1, '${JOIN_DATE}', '${JOIN_DATE}')" >> "$SQL_FILE"

  COUNT=$((COUNT + 1))
  echo "  [$COUNT/99] $USERNAME ($DISPLAY)"
done <<< "$NAMES"

echo ";" >> "$SQL_FILE"

echo "Inserting into database..."
sqlite3 "$DB_PATH" < "$SQL_FILE"

TOTAL=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM users;")
echo "Done. Total users: $TOTAL"