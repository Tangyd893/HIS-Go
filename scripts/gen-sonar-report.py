#!/usr/bin/env python3
"""Generate docs/sonarcloud-report.md from .sonar-raw.json (SonarCloud API snapshot)."""
import json
import sys
from collections import Counter
from datetime import date
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
RAW = ROOT / ".sonar-raw.json"
OUT = ROOT / "docs" / "sonarcloud-report.md"


def main() -> int:
    with RAW.open(encoding="utf-8-sig") as f:
        data = json.load(f)

    project_key = "Tangyd893_HIS-Go"
    project_name = "HIS-Go"
    today = date.today().isoformat()

    metric_map = {}
    for m in data["measures"].get("component", {}).get("measures", []):
        metric_map[m["metric"]] = m.get("value", "N/A")

    def get_metric(k: str, default: str = "N/A") -> str:
        return metric_map.get(k, default)

    bugs = get_metric("bugs", "0")
    vulns = get_metric("vulnerabilities", "0")
    smells = get_metric("code_smells", "0")
    ncloc = get_metric("ncloc", "N/A")
    dup = get_metric("duplicated_lines_density", "N/A")
    cog = get_metric("cognitive_complexity", "N/A")
    coverage = get_metric("coverage", "N/A")
    sqale = int(get_metric("sqale_index", "0") or 0)
    debt_h, debt_m = divmod(sqale, 60)

    gate_raw = data["gate"].get("projectStatus", {}).get("status", "")
    gate_map = {
        "OK": ("✅", "PASSED"),
        "ERROR": ("❌", "ERROR"),
        "WARN": ("⚠️", "WARN"),
    }
    gate_emoji, gate_status = gate_map.get(gate_raw, ("⚪", "NOT ANALYZED"))

    issues = data["issues"]
    total = len(issues)

    sev_order = ["BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"]
    sev_emoji = {
        "BLOCKER": "🔴",
        "CRITICAL": "🟠",
        "MAJOR": "🟡",
        "MINOR": "🔵",
        "INFO": "⚪",
    }
    sev_counts = Counter(i.get("severity", "INFO") for i in issues)

    type_counts = Counter(i.get("type", "CODE_SMELL") for i in issues)
    type_labels = {
        "BUG": "Bug",
        "VULNERABILITY": "Vulnerability",
        "CODE_SMELL": "Code Smell",
    }

    ext_map = {
        ".go": "Go",
        ".java": "Java",
        ".py": "Python",
        ".rs": "Rust",
        ".ts": "TypeScript",
        ".tsx": "TypeScript",
        ".js": "JavaScript",
        ".jsx": "JavaScript",
        ".vue": "Vue",
        ".sql": "SQL",
        ".sh": "Shell",
        ".bash": "Shell",
        ".yml": "YAML",
        ".yaml": "YAML",
        ".cs": "C#",
        ".css": "CSS",
        ".scss": "CSS",
    }
    cat_counts: Counter[str] = Counter()
    for i in issues:
        comp = i.get("component", "")
        path = comp.split(":", 1)[-1] if ":" in comp else comp
        ext = "." + path.rsplit(".", 1)[-1].lower() if "." in path else ""
        cat_counts[ext_map.get(ext, "Other")] += 1

    rule_counts = Counter(i.get("rule", "") for i in issues)
    rule_meta: dict[str, tuple[str, str]] = {}
    for i in issues:
        r = i.get("rule", "")
        if r not in rule_meta:
            rule_meta[r] = (i.get("severity", ""), i.get("type", ""))

    def sev_label(s: str) -> str:
        return f"{sev_emoji.get(s, '')} {s}" if s else s

    blockers = [i for i in issues if i.get("severity") == "BLOCKER"]
    criticals = [i for i in issues if i.get("severity") == "CRITICAL"]

    def issue_row(i: dict) -> str:
        comp = i.get("component", "")
        path = comp.split(":", 1)[-1] if ":" in comp else comp
        line = i.get("line", "")
        rule = i.get("rule", "")
        msg = i.get("message", "").replace("|", "\\|")
        return f"| `{path}` | {line} | `{rule}` | {msg} |"

    lines: list[str] = []
    lines.append("# SonarCloud 代码质量报告\n")
    lines.append(f"> {today}\n")
    lines.append("## 概览\n")
    lines.append(f"- **项目**: {project_name}")
    lines.append(f"- **Project Key**: {project_key}")
    lines.append(f"- **质量门**: {gate_emoji} {gate_status}")
    lines.append(f"- **代码行**: {ncloc}")
    lines.append(f"- **未解决 Issue**: {total}")
    lines.append(f"- **技术债务**: {debt_h}h {debt_m}min\n")

    lines.append("## 指标\n")
    lines.append("| 指标 | 值 |")
    lines.append("|------|-----|")
    lines.append(f"| Bug | {bugs} |")
    lines.append(f"| 漏洞 | {vulns} |")
    lines.append(f"| 代码异味 | {smells} |")
    lines.append(f"| 重复率 | {dup}% |")
    lines.append(f"| 认知复杂度 | {cog} |")
    lines.append(f"| 测试覆盖率 | {coverage}% |\n")

    lines.append("## 严重级别分布\n")
    lines.append("| 级别 | 数量 | 占比 |")
    lines.append("|------|------|------|")
    for s in sev_order:
        n = sev_counts.get(s, 0)
        pct = f"{n / total * 100:.1f}" if total else "0.0"
        lines.append(f"| {sev_emoji[s]} {s} | {n} | {pct}% |")
    lines.append("")

    lines.append("## 问题类型分布\n")
    lines.append("| 类型 | 数量 | 占比 |")
    lines.append("|------|------|------|")
    for t, label in type_labels.items():
        n = type_counts.get(t, 0)
        pct = f"{n / total * 100:.1f}" if total else "0.0"
        lines.append(f"| {label} | {n} | {pct}% |")
    lines.append("")

    lines.append("## Issue 分布（按文件类型）\n")
    lines.append("| 类别 | 数量 |")
    lines.append("|------|------|")
    for cat, n in sorted(cat_counts.items(), key=lambda x: -x[1]):
        lines.append(f"| {cat} | {n} |")
    lines.append("")

    lines.append("## 关键问题\n")
    lines.append("### 🔴 BLOCKER\n")
    if blockers:
        lines.append("| 文件 | 行 | 规则 | 描述 |")
        lines.append("|------|-----|------|------|")
        for i in blockers:
            lines.append(issue_row(i))
    else:
        lines.append("无")
    lines.append("")

    lines.append("### 🟠 CRITICAL（Top 20）\n")
    if criticals:
        lines.append("| 文件 | 行 | 规则 | 描述 |")
        lines.append("|------|-----|------|------|")
        for i in criticals[:20]:
            lines.append(issue_row(i))
    else:
        lines.append("无")
    lines.append("")

    lines.append("## 高频规则（Top 10）\n")
    lines.append("| 规则 | 级别 | 类型 | 命中数 | 占比 |")
    lines.append("|------|------|------|--------|------|")
    for rule, n in rule_counts.most_common(10):
        sev, typ = rule_meta.get(rule, ("", ""))
        pct = f"{n / total * 100:.1f}" if total else "0.0"
        lines.append(f"| `{rule}` | {sev_label(sev)} | {typ} | {n} | {pct}% |")

    OUT.write_text("\n".join(lines) + "\n", encoding="utf-8")  # NOSONAR: 输出路径为硬编码常量，非用户输入
    print(
        f"Report: {total} issues, gate={gate_status}, dup={dup}%, "
        f"vulns={vulns}, blockers={len(blockers)}"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
