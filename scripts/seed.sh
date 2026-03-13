#!/bin/bash
set -euo pipefail

SEED_DIR="seed"
DSN="${1:-${DSN:-pagoda.db}}"

exec_sql_file() {
  sqlite3 "$DSN" < "$1"
}

query_sql() {
  sqlite3 "$DSN" "$1"
}

EXISTING=$(query_sql "SELECT COUNT(*) FROM users;")
if [ -n "$EXISTING" ] && [ "$EXISTING" -gt 0 ] 2>/dev/null; then
  echo "Database already has $EXISTING users, skipping seed"
  exit 0
fi

HASH=$(htpasswd -nbBC 10 "" password | cut -d: -f2)

FIRST_NAMES=(
  Aiko Bram Cleo Dax Elara Fern Gideon Hana Idris Jade
  Kael Luna Milo Nyx Orion Piper Quinn Rune Sage Thorn
  Uma Vesper Wren Xara Yuki Zara Asher Blaze Cedar Dove
  Echo Flint Haven Ivy Jet Kit Lark Moss Nova Onyx
  Pearl Quill Rain Sky Terra Vale Wolf Yew Zen Aria
  Bay Coral Drift Elm Fox Gale Hawk Iris Juniper Kai
  Leaf Maple Nimbus Oak Pine River Storm Thistle Umber Violet
  Willow Birch Cliff Dawn Eve Frost Glen Heath Kite Lotus
  Meadow Night Olive Petal Reed Sol Tide Vine Wynn Briar
  Cloud Dune Fable Grove Haze Isle Kelp Lyric Mars Nectar
  Opal Plum Robin Star Twig Veil Wave Yarrow Zinnia Aster
  Brook Cinder Delta Ember Fjord Garnet Holly Indigo Jasper Lichen
  Mist Nettle Orbit Poppy Reef Sable Talon Vapor Wisp Aspen
  Brine Cobalt Dew Eon Flame Glint Heather Ink Knoll Lumen
  Moth Nile Pixel Quartz Ripple Shade Torrent Umbra Verge Zephyr
)

LAST_NAMES=(
  Ashworth Blackwood Clearwater Duskfall Evergreen Foxglove Goldleaf Hawthorn
  Ironbark Junewood Keelshore Larkspur Moonvale Nightbloom Oakshade Pinecrest
  Quillstone Ravenhill Silverbrook Thornwick Underhill Verdant Westwind Zenith
  Ashford Birchwood Copperfield Dawnstar Elderwood Fernside Greymist Holloway
  Icemere Jadeheart Kestrel Lakeshore Mossglen Northlight Overbrook Peakwood
  Quicksilver Rosewood Stonebridge Tidewater Umbervale Vineheart Wildmere Ashcroft
  Brightwell Coldspring Deepwood Evenfall Frostmere Glendale Hartwood Ivybridge
  Jasperfall Kirkwood Longmeadow Millbrook Newgrove Oldfield Ridgemont Sundale
  Thistlewood Upland Valewood Wintermere Brookfield Creekside Dellwood Eastmore
  Fairhaven Greenwood Highcliff Inkwell Lakewood Nightvale Pondcress Rustwood
)

PRONOUNS=(
  "zir/zem" "voi/vol" "qui/quem" "sol/solis" "lux/luxis"
  "nim/nir" "dra/drem" "pix/pixis" "nyx/nyxis" "vex/vir"
  "kal/kis" "ryn/rynis" "zel/zem" "orb/orbis" "mox/moxis"
  "jin/jinis" "tav/tavis" "wren/wrenis" "flux/fluxis" "ash/ashis"
  "glyph/glyphir" "hex/haxis" "nil/nilis" "arc/arcis" "cyr/cyris"
  "dex/dexis" "fen/fenis" "kez/kezis" "nym/nymis" "thri/threm"
)

LOCATIONS=(
  "Portland, Oregon" "London, England" "Tokyo, Japan" "Berlin, Germany" "Melbourne, Australia"
  "Toronto, Canada" "Brooklyn, New York" "Seattle, Washington" "Amsterdam, Netherlands" "Seoul, South Korea"
  "Austin, Texas" "Glasgow, Scotland" "Lisbon, Portugal" "Chicago, Illinois" "Stockholm, Sweden"
  "Vancouver, Canada" "Dublin, Ireland" "Paris, France" "Denver, Colorado" "Cape Town, South Africa"
  "San Francisco, California" "Helsinki, Finland" "Singapore" "Oslo, Norway" "Montreal, Canada"
  "Osaka, Japan" "Manchester, England" "Barcelona, Spain" "Wellington, New Zealand" "Reykjavik, Iceland"
  "Buenos Aires, Argentina" "Prague, Czech Republic" "Bangkok, Thailand" "Taipei, Taiwan" "Nairobi, Kenya"
  "Lima, Peru" "Kyoto, Japan" "Edinburgh, Scotland" "Bologna, Italy" "Zurich, Switzerland"
  "Copenhagen, Denmark" "Marrakech, Morocco" "Hanoi, Vietnam" "Bogota, Colombia" "Tallinn, Estonia"
  "Krakow, Poland" "Kuala Lumpur, Malaysia" "Mumbai, India" "Mexico City, Mexico" "Christchurch, New Zealand"
)

BIOS=()
while IFS= read -r line; do
  [ -n "$line" ] && BIOS+=("$line")
done < "$SEED_DIR/bios.html"

SIGS=()
while IFS= read -r line; do
  [ -n "$line" ] && SIGS+=("$line")
done < "$SEED_DIR/signatures.html"

BIO_COUNT=${#BIOS[@]}
SIG_COUNT=${#SIGS[@]}

TOTAL_USERS=$(( RANDOM % 31 + 120 ))
CITIZEN_COUNT=$(( TOTAL_USERS - 1 ))

ADMIN_COUNT=$(( CITIZEN_COUNT * 2 / 100 ))
MOD_COUNT=$(( CITIZEN_COUNT * 3 / 100 ))
BANNED_COUNT=$(( CITIZEN_COUNT * 4 / 100 ))
DISABLED_COUNT=$(( CITIZEN_COUNT * 3 / 100 ))
UNVERIFIED_COUNT=$(( CITIZEN_COUNT * 5 / 100 ))
[ "$ADMIN_COUNT" -lt 1 ] && ADMIN_COUNT=1
[ "$MOD_COUNT" -lt 1 ] && MOD_COUNT=1
[ "$BANNED_COUNT" -lt 1 ] && BANNED_COUNT=1
[ "$DISABLED_COUNT" -lt 1 ] && DISABLED_COUNT=1
[ "$UNVERIFIED_COUNT" -lt 1 ] && UNVERIFIED_COUNT=1

SLOTS=()
for ((i=0; i<CITIZEN_COUNT; i++)); do SLOTS+=("$i"); done
for ((i=${#SLOTS[@]}-1; i>0; i--)); do
  j=$(( RANDOM % (i + 1) ))
  tmp=${SLOTS[$i]}; SLOTS[$i]=${SLOTS[$j]}; SLOTS[$j]=$tmp
done

BIO_MAP=""
for ((i=0; i<BIO_COUNT && i<${#SLOTS[@]}; i++)); do
  BIO_MAP="${BIO_MAP},${SLOTS[$i]}:${i},"
done

SIG_POOL=()
for ((i=0; i<BIO_COUNT && i<${#SLOTS[@]}; i++)); do SIG_POOL+=("${SLOTS[$i]}"); done
for ((i=${#SIG_POOL[@]}-1; i>0; i--)); do
  j=$(( RANDOM % (i + 1) ))
  tmp=${SIG_POOL[$i]}; SIG_POOL[$i]=${SIG_POOL[$j]}; SIG_POOL[$j]=$tmp
done

SIG_MAP=""
for ((i=0; i<SIG_COUNT && i<${#SIG_POOL[@]}; i++)); do
  SIG_MAP="${SIG_MAP},${SIG_POOL[$i]}:${i},"
done

random_date_between() {
  local start_epoch=$1
  local end_epoch=$2
  local range=$(( end_epoch - start_epoch ))
  local offset=$(( RANDOM * RANDOM % range ))
  echo $(( start_epoch + offset ))
}

NOW_EPOCH=$(date +%s)
THREE_MONTHS_AGO=$(date -v-3m +%s 2>/dev/null || date -d "3 months ago" +%s)
ONE_MONTH_AGO=$(date -v-1m +%s 2>/dev/null || date -d "1 month ago" +%s)
ONE_MONTH_AHEAD=$(date -v+1m +%s 2>/dev/null || date -d "1 month" +%s)
FOUR_MONTHS_AGO=$(date -v-4m +%s 2>/dev/null || date -d "4 months ago" +%s)

OWNER_DATE=$(date -r $FOUR_MONTHS_AGO +"%Y-%m-%d %H:%M:%S+00:00" 2>/dev/null || date -d "@$FOUR_MONTHS_AGO" +"%Y-%m-%d %H:%M:%S+00:00")
OWNER_SEEN=$(date -u +"%Y-%m-%d %H:%M:%S+00:00")
OWNER_BIO='<p class="editor-paragraph"><i><em class="editor-italic" style="white-space: pre-wrap;">A really awesome cool slick ninja dinosaur thingy</em></i></p>'
OWNER_SIG='<p class="editor-paragraph"><span style="white-space: pre-wrap;">Love and Ciao</span></p>'

SQL_FILE=$(mktemp)
trap "rm -f $SQL_FILE" EXIT

escape_sql() {
  printf '%s' "$1" | sed "s/'/''/g"
}

echo "BEGIN TRANSACTION;" > "$SQL_FILE"

OWNER_BIO_ESC=$(escape_sql "$OWNER_BIO")
OWNER_SIG_ESC=$(escape_sql "$OWNER_SIG")

printf "INSERT OR IGNORE INTO users (username, email, password_hash, display_name, role, email_verified, jade, honor, pronouns, location, bio, signature, birthday, last_seen_at, ip, created_at, updated_at) VALUES ('master', 'master@pagoda.local', '%s', 'Master', 'owner', 1, 1000, 500, 'sol/solis', 'The Cloud', '%s', '%s', '1904-03-15 00:00:00+00:00', '%s', '127.0.0.1', '%s', '%s');\n" \
  "$HASH" "$OWNER_BIO_ESC" "$OWNER_SIG_ESC" "$OWNER_SEEN" "$OWNER_DATE" "$OWNER_DATE" >> "$SQL_FILE"

echo "Generating $CITIZEN_COUNT citizens..."

USED_USERNAMES="|"
COUNT=0

for ((i=0; i<CITIZEN_COUNT; i++)); do
  USERNAME=""
  while true; do
    FIRST=${FIRST_NAMES[$(( RANDOM % ${#FIRST_NAMES[@]} ))]}
    USERNAME=$(echo "$FIRST" | tr '[:upper:]' '[:lower:]')
    if [[ "$USED_USERNAMES" != *"|${USERNAME}|"* ]]; then
      break
    fi
    LAST=${LAST_NAMES[$(( RANDOM % ${#LAST_NAMES[@]} ))]}
    USERNAME=$(echo "${FIRST}${LAST}" | tr '[:upper:]' '[:lower:]')
    if [[ "$USED_USERNAMES" != *"|${USERNAME}|"* ]]; then
      break
    fi
  done
  USED_USERNAMES="${USED_USERNAMES}${USERNAME}|"

  DISPLAY_FIRST=${FIRST_NAMES[$(( RANDOM % ${#FIRST_NAMES[@]} ))]}
  DISPLAY_LAST=${LAST_NAMES[$(( RANDOM % ${#LAST_NAMES[@]} ))]}
  DISPLAY_NAME="$DISPLAY_FIRST $DISPLAY_LAST"

  PRONOUN=${PRONOUNS[$(( RANDOM % ${#PRONOUNS[@]} ))]}
  LOC=${LOCATIONS[$(( RANDOM % ${#LOCATIONS[@]} ))]}

  JOIN_EPOCH=$(random_date_between $THREE_MONTHS_AGO $NOW_EPOCH)
  JOIN_DATE=$(date -r $JOIN_EPOCH +"%Y-%m-%d %H:%M:%S+00:00" 2>/dev/null || date -d "@$JOIN_EPOCH" +"%Y-%m-%d %H:%M:%S+00:00")

  SEEN_OFFSET=$(( RANDOM % (90 * 24 * 3600) ))
  SEEN_EPOCH=$(( NOW_EPOCH - SEEN_OFFSET ))
  LAST_SEEN=$(date -r $SEEN_EPOCH +"%Y-%m-%d %H:%M:%S+00:00" 2>/dev/null || date -d "@$SEEN_EPOCH" +"%Y-%m-%d %H:%M:%S+00:00")

  BDAY_EPOCH=$(random_date_between $ONE_MONTH_AGO $ONE_MONTH_AHEAD)
  BDAY_MONTH=$(date -r $BDAY_EPOCH +"%m" 2>/dev/null || date -d "@$BDAY_EPOCH" +"%m")
  BDAY_DAY=$(date -r $BDAY_EPOCH +"%d" 2>/dev/null || date -d "@$BDAY_EPOCH" +"%d")
  BIRTHDAY="1904-${BDAY_MONTH}-${BDAY_DAY} 00:00:00+00:00"

  JADE=$(( RANDOM % 500 ))
  HONOR=$(( RANDOM % 200 ))

  BIO=""
  SIG=""
  if [[ "$BIO_MAP" == *",${i}:"* ]]; then
    BIO_IDX=$(echo "$BIO_MAP" | sed "s/.*,${i}:\([0-9]*\),.*/\1/")
    BIO=$(escape_sql "${BIOS[$BIO_IDX]}")
  fi
  if [[ "$SIG_MAP" == *",${i}:"* ]]; then
    SIG_IDX=$(echo "$SIG_MAP" | sed "s/.*,${i}:\([0-9]*\),.*/\1/")
    SIG=$(escape_sql "${SIGS[$SIG_IDX]}")
  fi

  ROLE="member"
  EMAIL_VERIFIED=1
  BANNED=0
  BANNED_AT=""
  BANNED_REASON=""
  BANNED_BY=""
  DISABLED=0
  DISABLED_AT=""
  DISABLED_REASON=""
  DISABLED_BY=""
  LINE_NUM=$(( i + 1 ))

  if [ "$LINE_NUM" -le "$ADMIN_COUNT" ]; then
    ROLE="admin"
  elif [ "$LINE_NUM" -le "$((ADMIN_COUNT + MOD_COUNT))" ]; then
    ROLE="moderator"
  elif [ "$LINE_NUM" -le "$((ADMIN_COUNT + MOD_COUNT + BANNED_COUNT))" ]; then
    BANNED=1
    BANNED_AT="'${JOIN_DATE}'"
    BANNED_REASON="'Violated community guidelines'"
    BANNED_BY=1
  elif [ "$LINE_NUM" -le "$((ADMIN_COUNT + MOD_COUNT + BANNED_COUNT + DISABLED_COUNT))" ]; then
    DISABLED=1
    DISABLED_AT="'${JOIN_DATE}'"
    DISABLED_REASON="'Account temporarily suspended'"
    DISABLED_BY=1
  elif [ "$LINE_NUM" -le "$((ADMIN_COUNT + MOD_COUNT + BANNED_COUNT + DISABLED_COUNT + UNVERIFIED_COUNT))" ]; then
    EMAIL_VERIFIED=0
  fi

  DISPLAY_ESC=$(escape_sql "$DISPLAY_NAME")
  LOC_ESC=$(escape_sql "$LOC")

  printf "INSERT OR IGNORE INTO users (username, email, password_hash, display_name, role, email_verified, jade, honor, pronouns, location, bio, signature, birthday, last_seen_at, account_banned, banned_at, banned_reason, banned_by, account_disabled, disabled_at, disabled_reason, disabled_by, ip, created_at, updated_at) VALUES ('%s', '%s@pagoda.local', '%s', '%s', '%s', %s, %s, %s, '%s', '%s', '%s', '%s', '%s', '%s', %s, %s, %s, %s, %s, %s, %s, %s, '127.0.0.1', '%s', '%s');\n" \
    "$USERNAME" "$USERNAME" "$HASH" "$DISPLAY_ESC" "$ROLE" "$EMAIL_VERIFIED" \
    "$JADE" "$HONOR" "$PRONOUN" "$LOC_ESC" "$BIO" "$SIG" \
    "$BIRTHDAY" "$LAST_SEEN" \
    "$BANNED" "${BANNED_AT:-NULL}" "${BANNED_REASON:-NULL}" "${BANNED_BY:-NULL}" \
    "$DISABLED" "${DISABLED_AT:-NULL}" "${DISABLED_REASON:-NULL}" "${DISABLED_BY:-NULL}" \
    "$JOIN_DATE" "$JOIN_DATE" >> "$SQL_FILE"

  COUNT=$((COUNT + 1))
  STATUS=""
  if [ "$BANNED" -eq 1 ]; then STATUS=" [banned]";
  elif [ "$DISABLED" -eq 1 ]; then STATUS=" [disabled]";
  elif [ "$EMAIL_VERIFIED" -eq 0 ]; then STATUS=" [unverified]";
  elif [ "$ROLE" != "member" ]; then STATUS=" [$ROLE]";
  fi
  echo "  [$COUNT/$CITIZEN_COUNT] $USERNAME ($DISPLAY_NAME)$STATUS"
done

echo "COMMIT;" >> "$SQL_FILE"

echo "Inserting into database..."
exec_sql_file "$SQL_FILE"

TOTAL=$(query_sql "SELECT COUNT(*) FROM users;")
echo "Done. Total users: $TOTAL (admins: $ADMIN_COUNT, mods: $MOD_COUNT, banned: $BANNED_COUNT, disabled: $DISABLED_COUNT, unverified: $UNVERIFIED_COUNT)"