#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://127.0.0.1:3333/api"
COOKIES_FILE="cookies.txt"

USERNAME="${USERNAME:-}"
PASSWORD="${PASSWORD:-}"

# Articles: file path | category name
ARTICLES=(
  "content/articles/营养健康—0-6岁分龄营养实用指南.md|营养健康"
  "content/articles/早期教育—敏感期与家庭可操作教养法.md|早期教育"
  "content/articles/心理发展—依恋与情绪调节的家庭方案.md|心理发展"
  "content/articles/安全防护—家庭与外出一站式安全清单.md|安全防护"
)

if [[ -z "$USERNAME" || -z "$PASSWORD" ]]; then
  echo "ERROR: Please set USERNAME and PASSWORD env vars before running." >&2
  echo "Example: USERNAME=zhangzhi PASSWORD='yourpass' scripts/publish_articles.sh" >&2
  exit 1
fi

echo "Logging in as ${USERNAME} ..."
curl -sS -X POST \
  -H 'Content-Type: application/json' \
  -c "${COOKIES_FILE}" \
  -d "{\"username\":\"${USERNAME}\",\"password\":\"${PASSWORD}\"}" \
  "${BASE_URL}/auth/login" > /dev/null

echo "Fetching categories ..."
CATS_JSON=$(curl -sS -X GET -b "${COOKIES_FILE}" "${BASE_URL}/categories/all")
CAT_FILE=$(mktemp)
printf "%s" "$CATS_JSON" > "$CAT_FILE"

get_category_id() {
  local name="$1"
  python3 - "$name" "$CAT_FILE" <<'PY'
import sys, json
name = sys.argv[1]
cat_file = sys.argv[2]
with open(cat_file, 'r', encoding='utf-8') as f:
    data = json.load(f)
cats = data.get('data') or []
mapping = { c.get('name'): c.get('id') for c in cats }
cid = mapping.get(name)
if cid is None:
    sys.exit(2)
print(cid)
PY
}

post_article() {
  local file="$1"; shift
  local category_id="$1"; shift
  tmp_json=$(mktemp)
  python3 - "$file" "$category_id" >"$tmp_json" <<'PY'
import sys, json, re, time
file = sys.argv[1]
category_id = int(sys.argv[2])

with open(file, 'r', encoding='utf-8') as f:
    lines = f.read().splitlines()

title = ''
summary = ''
tags = ''
start_idx = 0
for i, line in enumerate(lines):
    if line.startswith('Title:'):
        title = line.split(':', 1)[1].strip()
    elif line.startswith('Summary:'):
        summary = line.split(':', 1)[1].strip()
    elif line.startswith('Tags:'):
        tags = line.split(':', 1)[1].strip()
        start_idx = i + 1
        break

content = '\n'.join(lines[start_idx:]).strip()

base = re.sub(r'[^a-zA-Z0-9]+', '-', title).strip('-').lower()
if not base:
    base = 'article'
slug = f"{base}-{int(time.time()*1000)}"

body = {
    'title': title,
    'slug': slug,
    'summary': summary,
    'content': content,
    'cover_image': '',
    'category_id': category_id,
    'tags': tags,
    'is_top': False,
    'is_recommend': False,
    'status': 1,
}
print(json.dumps(body, ensure_ascii=False))
PY

  resp=$(curl -sS -X POST \
    -H 'Content-Type: application/json' \
    -b "${COOKIES_FILE}" \
    --data-binary @"$tmp_json" \
    "${BASE_URL}/articles")
  rm -f "$tmp_json"

  echo "$resp" | python3 - <<'PY'
import sys, json
resp = sys.stdin.read()
try:
    data = json.loads(resp)
    code = data.get('code', 200)
    if code not in (0, 200):
        print(f"[FAIL] API error: {data.get('message')}")
        sys.exit(1)
    art = data.get('data') or {}
    print(f"[OK] Created article ID={art.get('id')} title={art.get('title')}")
except Exception:
    print('[FAIL] Invalid response:', resp)
    sys.exit(1)
PY
}

for item in "${ARTICLES[@]}"; do
  file="${item%%|*}"
  cname="${item##*|}"
  if [[ ! -f "$file" ]]; then
    echo "Skip: $file not found" >&2
    continue
  fi
  cid=$(get_category_id "$cname") || { echo "Category '$cname' not found" >&2; exit 1; }
  echo "Publishing: $file -> ${cname}(ID=$cid) ..."
  post_article "$file" "$cid"
done

echo "All done."
