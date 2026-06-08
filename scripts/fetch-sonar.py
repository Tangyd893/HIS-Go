#!/usr/bin/env python3
"""Fetch SonarCloud API snapshot to .sonar-raw.json."""
import json
import os
import sys
import time
import urllib.request
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
RAW = ROOT / ".sonar-raw.json"


def main() -> int:
    token = os.environ.get("SONAR_TOKEN")
    if not token:
        print("SONAR_TOKEN is not set", file=sys.stderr)
        return 1

    key = "Tangyd893_HIS-Go"
    headers = {"Authorization": f"Bearer {token}"}

    def get(url: str) -> dict:
        req = urllib.request.Request(url, headers=headers)
        with urllib.request.urlopen(req) as resp:
            return json.loads(resp.read().decode())

    metrics = get(
        "https://sonarcloud.io/api/measures/component"
        f"?component={key}&metricKeys=bugs,vulnerabilities,code_smells,ncloc,"
        "duplicated_lines_density,cognitive_complexity,sqale_index,coverage,quality_gate_details"
    )
    time.sleep(0.25)

    gate = get(
        f"https://sonarcloud.io/api/qualitygates/project_status?projectKey={key}"
    )
    time.sleep(0.25)

    issues: list = []
    page = 1
    while True:
        data = get(
            "https://sonarcloud.io/api/issues/search"
            f"?componentKeys={key}&ps=500&resolved=false"
            f"&types=BUG,VULNERABILITY,CODE_SMELL&p={page}"
        )
        issues.extend(data.get("issues", []))
        total = data.get("paging", {}).get("total", 0)
        if page * 500 >= total:
            break
        page += 1
        time.sleep(0.25)

    raw = {"measures": metrics, "gate": gate, "issues": issues}
    RAW.write_text(json.dumps(raw, ensure_ascii=False, indent=2), encoding="utf-8")
    status = gate.get("projectStatus", {}).get("status", "UNKNOWN")
    print(f"Fetched {len(issues)} issues, gate={status}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
