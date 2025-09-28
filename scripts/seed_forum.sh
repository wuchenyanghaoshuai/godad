#!/usr/bin/env bash
set -euo pipefail

# Seed the forum (community) with posts and replies using provided accounts.
# - Creates 5–8 posts per topic in /api/forum/topics
# - Adds 2–4 replies per post from mixed accounts (some nested)

BASE_URL="${BASE_URL:-http://127.0.0.1:3333/api}"

# Default accounts (username|password). Override by exporting ACCOUNTS env var with same format
ACCOUNTS_DEFAULT=(
  "chenyangwu|543521"
  "dongxueli|543521"
  "zhangzhi|WUchenyang@123"
  "admin|WUchenyang@123"
)

if [[ -n "${ACCOUNTS:-}" ]]; then
  # If user provides ACCOUNTS as newline or space separated username|password entries
  # shellcheck disable=SC2206
  ACCOUNTS=( ${ACCOUNTS} )
else
  ACCOUNTS=( "${ACCOUNTS_DEFAULT[@]}" )
fi

log() { printf "[%s] %s\n" "$(date +'%H:%M:%S')" "$*"; }

# Simple random helpers
rand_between() { # inclusive
  local min=$1 max=$2
  echo $(( RANDOM % (max - min + 1) + min ))
}

json_field() { # read JSON from stdin, print a.b.c field path
  local path="$1"
  python3 - "$path" <<'PY' || true
import sys, json
path = sys.argv[1].split('.')
try:
    data = json.load(sys.stdin)
    for p in path:
        if isinstance(data, dict):
            data = data.get(p)
        else:
            data = None
            break
    if data is None:
        raise KeyError
    if isinstance(data, (dict, list)):
        import json as _json
        print(_json.dumps(data, ensure_ascii=False))
    else:
        print(str(data))
except Exception:
    sys.exit(1)
PY
}

# Read JSON field from file
json_field_file() {
  local file="$1"; shift
  local path="$1"
  python3 - "$file" "$path" <<'PY' || true
import sys, json, pathlib
file, path = sys.argv[1], sys.argv[2]
parts = path.split('.')
try:
    s = pathlib.Path(file).read_text(encoding='utf-8')
    data = json.loads(s)
    for p in parts:
        if isinstance(data, dict):
            data = data.get(p)
        else:
            data = None
            break
    if data is None:
        raise KeyError
    if isinstance(data, (dict, list)):
        import json as _json
        print(_json.dumps(data, ensure_ascii=False))
    else:
        print(str(data))
except Exception:
    sys.exit(1)
PY
}

# Login all accounts and capture tokens
declare -a USERS=()
declare -a TOKENS=()

log "Logging in with ${#ACCOUNTS[@]} accounts ..."
for entry in "${ACCOUNTS[@]}"; do
  username="${entry%%|*}"
  password="${entry##*|}"

  log "→ Login as ${username}"
  tmp_login=$(mktemp)
  curl -sS -X POST \
    -H 'Content-Type: application/json' \
    --data "{\"username\":\"${username}\",\"password\":\"${password}\"}" \
    "${BASE_URL}/auth/login" -o "$tmp_login" || { log "curl failed for login"; exit 1; }

  code=$(json_field_file "$tmp_login" code || echo 0)
  if [[ "$code" != "200" ]]; then
    msg=$(json_field_file "$tmp_login" message || echo "")
    log "!! Login failed for ${username}: ${msg}"
    rm -f "$tmp_login"
    continue
  fi

  token=$(json_field_file "$tmp_login" 'data.token' || true)
  if [[ -z "$token" || "$token" == "null" ]]; then
    log "!! No token in login response for ${username}"
    rm -f "$tmp_login"
    continue
  fi

  USERS+=("$username")
  TOKENS+=("$token")
  rm -f "$tmp_login"
done

if [[ ${#TOKENS[@]} -eq 0 ]]; then
  log "ERROR: No accounts could log in successfully. Abort."
  exit 1
fi

log "Fetched tokens for ${#TOKENS[@]} accounts."

# Fetch forum topics
log "Fetching forum topics ..."
tmp_topics=$(mktemp)
curl -sS -X GET "${BASE_URL}/forum/topics" -o "$tmp_topics" || { log "Failed to fetch topics"; exit 1; }
code=$(json_field_file "$tmp_topics" code || echo 0)
if [[ "$code" != "200" ]]; then
  msg=$(json_field_file "$tmp_topics" message || echo "")
  log "ERROR: Cannot get topics: $msg"
  rm -f "$tmp_topics"
  exit 1
fi

# Parse topics array (macOS default bash doesn't support mapfile)
TOPICS=()
tmp_list=$(mktemp)
python3 - "$tmp_topics" <<'PY' > "$tmp_list"
import sys, json, pathlib
path = sys.argv[1]
try:
    d = json.loads(pathlib.Path(path).read_text(encoding='utf-8'))
    for t in (d.get('data') or []):
        if isinstance(t, str):
            print(t)
except Exception as e:
    # write nothing, handled by shell
    pass
PY
while IFS= read -r line; do
  [[ -n "$line" ]] && TOPICS+=("$line")
done < "$tmp_list"
rm -f "$tmp_list"

if [[ ${#TOPICS[@]} -eq 0 ]]; then
  log "ERROR: No topics returned. Abort."
  rm -f "$tmp_topics"
  exit 1
fi

log "Got ${#TOPICS[@]} topics. Will create 5-8 posts per topic."
rm -f "$tmp_topics"

# Helper: map English topic to Chinese label for nicer titles/content
topic_label_zh() {
  case "$1" in
    "Baby Care") echo "婴儿护理" ;;
    "Feeding") echo "喂养" ;;
    "Sleep") echo "睡眠" ;;
    "Health") echo "健康" ;;
    "Development") echo "发育" ;;
    "Activities") echo "活动" ;;
    "Gear") echo "用品" ;;
    "Parenting") echo "育儿" ;;
    "Family Life") echo "家庭生活" ;;
    "Work & Life Balance") echo "工作生活平衡" ;;
    "Relationships") echo "人际关系" ;;
    "Mental Health") echo "心理健康" ;;
    "Finances") echo "财务" ;;
    "Legal") echo "法律" ;;
    *) echo "其他" ;;
  esac
}

# Generate a post payload (title, content) for a given topic
gen_post_payload() {
  local topic_en="$1" topic_zh="$2"
  python3 - "$topic_en" "$topic_zh" <<'PY'
import sys, json, random
topic_en, topic_zh = sys.argv[1], sys.argv[2]
title_tpls = [
  f"求助：关于{topic_zh}的困扰",
  f"分享：我的{topic_zh}经验与做法",
  f"讨论：大家如何处理{topic_zh}？",
  f"新手问题：{topic_zh}有哪些注意事项？",
  f"请教：{topic_zh}方面有没有推荐的方法",
  f"经验交流：{topic_zh}的实用技巧",
  f"避坑指南：{topic_zh}常见误区",
  f"好物推荐：{topic_zh}相关用品/资源",
]
title = random.choice(title_tpls)

paragraphs = []
paragraphs.append(f"背景：最近在实践{topic_zh}时遇到一些问题，想听听大家的意见和经验。")
paragraphs.append(f"我的情况：1）目前尝试过几种做法；2）参考了文章和医生建议；3）仍存在疑惑，比如细节节奏、频次安排等。")
paragraphs.append(f"问题点：A）{topic_zh}的关键步骤；B）不同年龄阶段的差异；C）如何建立稳定习惯。")
paragraphs.append(f"求助：如果您在{topic_zh}方面有成熟做法，能否分享具体流程、节奏安排以及注意事项？非常感谢！")

content = "\n\n".join(paragraphs)
print(json.dumps({"title": title, "content": content, "topic": topic_en}, ensure_ascii=False))
PY
}

# Generate a reply payload for a post
gen_reply_payload() {
  local post_id="$1" topic_zh="$2"
  local parent_id="${3:-}"
  python3 - "$post_id" "$topic_zh" "$parent_id" <<'PY'
import sys, json, random
post_id = int(sys.argv[1])
topic_zh = sys.argv[2]
parent_id = sys.argv[3] if len(sys.argv) > 3 else ''

snippets = [
  f"我家也遇到过类似的{topic_zh}问题，可以尝试逐步过渡的方法。",
  f"建议先建立稳定的节奏，再慢慢微调；{topic_zh}确实需要耐心。",
  f"可以记录一周数据对比，{topic_zh}有时候受环境影响很大。",
  f"我的经验：把关键步骤拆小，一次只优化一个点，会更容易坚持。",
  f"如果条件允许，可以咨询下专业人士，尤其是{topic_zh}方面的细节。",
]
content = random.choice(snippets)
body = {"post_id": post_id, "content": content}
if parent_id and parent_id.lower() != 'none':
    try:
        body["parent_id"] = int(parent_id)
    except Exception:
        pass
print(json.dumps(body, ensure_ascii=False))
PY
}

# Create a forum post, return new post id on stdout
create_post() {
  local token="$1" topic_en="$2" topic_zh="$3"
  tmp=$(mktemp)
  gen_post_payload "$topic_en" "$topic_zh" >"$tmp"
  resp_file=$(mktemp)
  curl -sS -X POST \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer ${token}" \
    --data-binary @"$tmp" \
    "${BASE_URL}/forum/posts" -o "$resp_file"
  rm -f "$tmp"

  code=$(json_field_file "$resp_file" code || echo 0)
  if [[ "$code" != "200" ]]; then
    msg=$(json_field_file "$resp_file" message || echo "")
    log "!! Create post failed: ${msg}"
    rm -f "$resp_file"
    return 1
  fi
  post_id=$(json_field_file "$resp_file" 'data.id' || true)
  if [[ -z "$post_id" ]]; then
    log "!! No post id returned"
    rm -f "$resp_file"
    return 1
  fi
  rm -f "$resp_file"
  printf "%s" "$post_id"
}

# Create a reply, return reply id on stdout
create_reply() {
  local token="$1" post_id="$2" topic_zh="$3" parent_id="${4:-}"
  tmp=$(mktemp)
  gen_reply_payload "$post_id" "$topic_zh" "${parent_id:-}" >"$tmp"
  resp_file=$(mktemp)
  curl -sS -X POST \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer ${token}" \
    --data-binary @"$tmp" \
    "${BASE_URL}/forum/replies" -o "$resp_file"
  rm -f "$tmp"

  code=$(json_field_file "$resp_file" code || echo 0)
  if [[ "$code" != "200" ]]; then
    msg=$(json_field_file "$resp_file" message || echo "")
    log "!! Create reply failed: ${msg}"
    rm -f "$resp_file"
    return 1
  fi
  rid=$(json_field_file "$resp_file" 'data.id' || true)
  if [[ -z "$rid" ]]; then
    log "!! No reply id returned"
    rm -f "$resp_file"
    return 1
  fi
  rm -f "$resp_file"
  printf "%s" "$rid"
}

declare -a CREATED_POSTS=() # elements: "postId|authorIndex|topicEn|topicZh"
total_posts=0
total_replies=0

for topic in "${TOPICS[@]}"; do
  topic_zh=$(topic_label_zh "$topic")
  n=$(rand_between 5 8)
  log "Topic '${topic}'(${topic_zh}): creating ${n} posts ..."

  for ((i=0; i<n; i++)); do
    idx=$(rand_between 0 $(( ${#TOKENS[@]} - 1 )))
    token="${TOKENS[$idx]}"
    author="${USERS[$idx]}"
    if post_id=$(create_post "$token" "$topic" "$topic_zh"); then
      log "  [OK] Post #$((i+1)) by ${author}: ID=${post_id}"
      CREATED_POSTS+=("${post_id}|${idx}|${topic}|${topic_zh}")
      total_posts=$((total_posts+1))
    else
      log "  [FAIL] Post #$((i+1)) by ${author}"
    fi
  done
done

log "Created ${total_posts} posts across ${#TOPICS[@]} topics."

# Create replies
log "Creating replies (2-4 per post, mixed accounts) ..."
for entry in "${CREATED_POSTS[@]}"; do
  post_id="${entry%%|*}"
  rest="${entry#*|}"
  author_idx="${rest%%|*}"
  rest2="${rest#*|}"
  topic_en="${rest2%%|*}"
  topic_zh="${rest2#*|}"

  # Fetch post detail for contextual replies
  POST_TITLE=""; POST_CONTENT=""
  tmp_post=$(mktemp)
  curl -sS -X GET "${BASE_URL}/forum/posts/${post_id}" -o "$tmp_post" || true
  # Only try to parse if file non-empty
  if [[ -s "$tmp_post" ]]; then
    POST_TITLE=$(json_field_file "$tmp_post" 'data.title' || echo "")
    POST_CONTENT=$(json_field_file "$tmp_post" 'data.content' || echo "")
  fi
  rm -f "$tmp_post"

  rcount=$(rand_between 2 4)
  first_reply_id=""
  for ((r=0; r<rcount; r++)); do
    # Choose a replier (prefer different from author)
    ridx=$(rand_between 0 $(( ${#TOKENS[@]} - 1 )))
    if [[ $ridx -eq $author_idx && ${#TOKENS[@]} -gt 1 ]]; then
      ridx=$(( (ridx + 1) % ${#TOKENS[@]} ))
    fi
    rtoken="${TOKENS[$ridx]}"
    ruser="${USERS[$ridx]}"

    parent=""
    # 30% chance to reply to first reply to form a thread
    if [[ -n "$first_reply_id" && $((RANDOM % 10)) -lt 3 ]]; then
      parent="$first_reply_id"
    fi

    # Build contextual reply body (fallback to generic if generation fails)
    tmp_body=$(mktemp)
    python3 - "$post_id" "$topic_en" "$topic_zh" "$POST_TITLE" "$POST_CONTENT" "$parent" <<'PY' > "$tmp_body" || true
import sys, json, random, re
post_id = int(sys.argv[1])
topic_en, topic_zh = sys.argv[2], sys.argv[3]
title, content, parent = sys.argv[4], sys.argv[5], sys.argv[6]

rules = {
  'Sleep': {
    'kw': ['夜醒','午睡','哄睡','作息','入睡','早醒','奶睡','抱睡','白噪音','黑暗','睡眠倒退','过度疲劳'],
    'tips': [
      '固定就寝流程（洗漱-灯光-安静活动）建立入睡暗示',
      '拉平白天清醒/午睡间隔，避免过度疲劳',
      '就寝环境尽量黑暗安静，适度白噪音帮助入睡',
      '逐步减少喂奶/抱睡依赖，改为轻拍+口头安抚',
      '记录一周作息，按15分钟粒度微调就寝时间'
    ]
  },
  'Feeding': {
    'kw': ['奶量','配方奶','母乳','加餐','断奶','喂养间隔','添加辅食','胃口','胀气','返流'],
    'tips': [
      '按体重与月龄评估目标奶量，避免强喂',
      '固定喂养间隔，观察饥饱信号而非只看时钟',
      '若有胀气返流，减少每次量、放慢节奏并拍嗝',
      '添加辅食遵循单一到混合、由稀到稠、由少到多',
      '记录三天摄入与反应，调整配方/时间点'
    ]
  },
  'Baby Care': {
    'kw': ['洗澡','皮疹','湿疹','脐带','黄疸','体温','抱姿','喂药'],
    'tips': [
      '皮肤护理以保湿为主，尽量选择温和无香配方',
      '洗澡水温37-40℃，时长控制在5-10分钟',
      '脐部保持干燥清洁，渗出或红肿及时就医',
      '体温异常或精神差时，优先排查感染/脱水'
    ]
  },
  'Health': {
    'kw': ['发烧','咳嗽','感冒','腹泻','疫苗','湿疹','过敏'],
    'tips': [
      '高热伴精神差/呼吸急促请尽快就医',
      '充分补液与休息，警惕持续高热超过48小时',
      '遵医嘱用药，避免重复含同成分退热药',
      '记录发作诱因以评估过敏/哮喘可能'
    ]
  },
  'Development': {
    'kw': ['抬头','翻身','爬','走','精细动作','语言','里程碑'],
    'tips': [
      '创设安全可探索环境，增加趴玩/地面活动时间',
      '以游戏方式鼓励高频短时练习，避免过度纠正',
      '和同月龄比起点差异很大，更看重连续进步',
      '若多个里程碑明显滞后，建议评估排查'
    ]
  }
}

def pick(seq, n):
  seq = list(seq)
  random.shuffle(seq)
  return seq[:n]

bucket = rules.get(topic_en, {
  'kw': ['计划','节奏','记录','观察'],
  'tips': ['先记录现状再小步调整','一次只改一个变量便于复盘','关注孩子信号而非只看表格']
})

text = (title or '') + ' ' + (content or '')
hits = [k for k in bucket['kw'] if k in text]
key_phr = '、'.join(pick(hits or bucket['kw'], min(2, len(bucket['kw']))))

intro_opts = [
  f"关于你提到的「{key_phr}」，我这边的做法是：",
  f"围绕「{key_phr}」这个点，可以参考下面思路：",
  f"我之前也遇到过类似情况（涉及{key_phr}），给你一些可执行建议："
]

tips = pick(bucket['tips'], min(3, len(bucket['tips'])))

closing_opts = [
  '建议连续记录3-7天再评估效果，耐心很重要～',
  '每个孩子差异很大，稳步微调比一步到位更可靠。',
  '若出现异常信号（精神差、进食显著下降等）尽快就医。'
]

lines = [random.choice(intro_opts)]
for t in tips:
  lines.append('· ' + t)
lines.append(random.choice(closing_opts))

body = { 'post_id': post_id, 'content': "\n".join(lines) }
if parent and parent.lower()!='none':
  try:
    body['parent_id'] = int(parent)
  except Exception:
    pass
print(json.dumps(body, ensure_ascii=False))
PY

    # Actually call create_reply for consistent parsing and id extraction
    if rid=$(create_reply "$rtoken" "$post_id" "$topic_zh" "${parent}"); then
      [[ -z "$first_reply_id" ]] && first_reply_id="$rid"
      total_replies=$((total_replies+1))
      :
    else
      log "  [FAIL] Reply on post ${post_id} by ${ruser}"
    fi
    rm -f "$tmp_body"
  done
done

log "Seeding completed: ${total_posts} posts, ${total_replies} replies."
