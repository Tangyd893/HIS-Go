#!/usr/bin/env python3
"""
build-triage-embeddings.py — 就诊助手知识库嵌入构建脚本
用法:
  基本:  python scripts/build-triage-embeddings.py
  指定:  SILICONFLOW_API_KEY=sk-xxx SILICONFLOW_EMBED_MODEL=BAAI/bge-m3 python scripts/build-triage-embeddings.py

输入:  backend/data/triage/knowledge.json
输出:  backend/data/triage/chunks.json   (切分后的文本块)
       backend/data/triage/embeddings.json (BGE-M3 向量嵌入)
"""

import json
import os
import sys
import time
import urllib.request
import urllib.error

# ===== 配置 =====
KNOWLEDGE_PATH = os.path.join(os.path.dirname(__file__), "..", "backend", "data", "triage", "knowledge.json")
CHUNKS_PATH = os.path.join(os.path.dirname(__file__), "..", "backend", "data", "triage", "chunks.json")
EMBEDDINGS_PATH = os.path.join(os.path.dirname(__file__), "..", "backend", "data", "triage", "embeddings.json")

API_KEY = os.environ.get("SILICONFLOW_API_KEY", "")
API_BASE = os.environ.get("DEEPSEEK_BASE_URL", "https://api.siliconflow.cn")
EMBED_MODEL = os.environ.get("SILICONFLOW_EMBED_MODEL", "BAAI/bge-m3")

BATCH_SIZE = 16  # 每批最多 16 条
REQUEST_DELAY = 0.3  # 请求间隔（秒）


def load_knowledge(path: str) -> list[dict]:
    """加载知识库 JSON"""
    if not os.path.exists(path):
        print(f"[ERROR] 知识库文件不存在: {path}")
        sys.exit(1)
    with open(path, "r", encoding="utf-8") as f:
        return json.load(f)


def chunk_knowledge(entries: list[dict]) -> list[dict]:
    """
    将知识库条目切分为可检索的文本块
    每个条目生成多个块: 关键词块 + 分类描述块 + 科室推荐块 + 注意事项块
    """
    chunks = []
    for entry in entries:
        kid = entry["id"]
        # 块1: 关键词（用于关键词匹配检索）
        chunks.append({
            "id": f"{kid}-kw",
            "source_id": kid,
            "text": f"症状关键词: {', '.join(entry['keywords'])}",
            "type": "keywords",
            "dept_types": entry["dept_types"],
            "urgency": entry["urgency"],
            "category": entry["category"],
            "notes": entry["notes"],
        })
        # 块2: 病症分类描述
        chunks.append({
            "id": f"{kid}-desc",
            "source_id": kid,
            "text": f"病症分类: {entry['category']}。常见症状包括: {', '.join(entry['keywords'][:8])}等。建议挂号科室: {', '.join(entry['dept_types'])}。",
            "type": "description",
            "dept_types": entry["dept_types"],
            "urgency": entry["urgency"],
            "category": entry["category"],
            "notes": entry["notes"],
        })
        # 块3: 注意事项（如果有）
        if entry.get("notes"):
            chunks.append({
                "id": f"{kid}-note",
                "source_id": kid,
                "text": f"{entry['category']}注意事项: {entry['notes']}",
                "type": "notes",
                "dept_types": entry["dept_types"],
                "urgency": entry["urgency"],
                "category": entry["category"],
                "notes": entry["notes"],
            })
    return chunks


def get_embeddings(texts: list[str]) -> list[list[float]]:
    """调用 SiliconFlow BGE-M3 API 获取批量向量嵌入"""
    if not API_KEY:
        print("[WARN] SILICONFLOW_API_KEY 未设置，将生成随机向量作为占位")
        print("[WARN] 请在 demo.env 中配置 SILICONFLOW_API_KEY 以启用语义检索")
        import random
        random.seed(42)
        return [[random.uniform(-0.1, 0.1) for _ in range(1024)] for _ in texts]

    url = f"{API_BASE}/v1/embeddings"
    headers = {
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json",
    }
    payload = {
        "model": EMBED_MODEL,
        "input": texts,
        "encoding_format": "float",
    }

    req = urllib.request.Request(url, data=json.dumps(payload).encode("utf-8"), headers=headers, method="POST")
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            result = json.loads(resp.read().decode("utf-8"))
            embeddings = [item["embedding"] for item in sorted(result["data"], key=lambda x: x["index"])]
            dim = len(embeddings[0]) if embeddings else 0
            print(f"[OK] 获取 {len(embeddings)} 个嵌入向量，维度={dim}")
            return embeddings
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        print(f"[ERROR] API 请求失败 ({e.code}): {body}")
        sys.exit(1)
    except Exception as e:
        print(f"[ERROR] 网络请求异常: {e}")
        sys.exit(1)


def compute_embeddings_batched(chunks: list[dict]) -> list[dict]:
    """批量计算所有块的嵌入向量"""
    texts = [c["text"] for c in chunks]
    all_embeddings = []

    total = len(texts)
    for i in range(0, total, BATCH_SIZE):
        batch = texts[i : i + BATCH_SIZE]
        print(f"  处理批次 {i // BATCH_SIZE + 1}/{(total + BATCH_SIZE - 1) // BATCH_SIZE} ({len(batch)} 条)...")
        batch_embeddings = get_embeddings(batch)
        all_embeddings.extend(batch_embeddings)
        if i + BATCH_SIZE < total:
            time.sleep(REQUEST_DELAY)

    # 将嵌入附到块上
    for chunk, emb in zip(chunks, all_embeddings):
        chunk["embedding"] = emb
    return chunks


def main():
    print("=" * 50)
    print("HIS-Go 就诊助手 — 知识库嵌入构建")
    print("=" * 50)
    print(f"知识库: {KNOWLEDGE_PATH}")
    print(f"嵌入模型: {EMBED_MODEL}")
    print(f"API Key: {'已设置' if API_KEY else '未设置（将使用占位向量）'}")
    print()

    # 1. 加载知识库
    print("[1/3] 加载知识库...")
    entries = load_knowledge(KNOWLEDGE_PATH)
    print(f"  加载 {len(entries)} 条知识条目")

    # 2. 切分
    print("[2/3] 切分文本块...")
    chunks = chunk_knowledge(entries)
    print(f"  生成 {len(chunks)} 个文本块")

    # 保存 chunks
    chunks_for_save = [{k: v for k, v in c.items() if k != "embedding"} for c in chunks]
    os.makedirs(os.path.dirname(CHUNKS_PATH), exist_ok=True)
    with open(CHUNKS_PATH, "w", encoding="utf-8") as f:
        json.dump(chunks_for_save, f, ensure_ascii=False, indent=2)
    print(f"  已保存: {CHUNKS_PATH}")

    # 3. 嵌入
    print("[3/3] 计算嵌入向量...")
    chunks_with_emb = compute_embeddings_batched(chunks)
    print(f"  生成 {len(chunks_with_emb)} 个嵌入向量")

    # 保存 embeddings
    os.makedirs(os.path.dirname(EMBEDDINGS_PATH), exist_ok=True)
    with open(EMBEDDINGS_PATH, "w", encoding="utf-8") as f:
        json.dump(chunks_with_emb, f, ensure_ascii=False)
    print(f"  已保存: {EMBEDDINGS_PATH}")

    print()
    print("[DONE] 嵌入构建完成！")
    print(f"  chunks:     {CHUNKS_PATH}")
    print(f"  embeddings: {EMBEDDINGS_PATH}")


if __name__ == "__main__":
    main()
