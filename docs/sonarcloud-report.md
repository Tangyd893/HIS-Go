# SonarCloud 代码质量报告

> 2026-06-10

## 概览

- **项目**: HIS-Go
- **Project Key**: Tangyd893_HIS-Go
- **质量门**: ❌ ERROR
- **代码行**: 42396
- **未解决 Issue**: 555
- **技术债务**: 44h 30min

## 指标

| 指标 | 值 |
|------|-----|
| Bug | 0 |
| 漏洞 | 1 |
| 代码异味 | 554 |
| 重复率 | 5.7% |
| 认知复杂度 | 3889 |
| 测试覆盖率 | N/A% |

## 严重级别分布

| 级别 | 数量 | 占比 |
|------|------|------|
| 🔴 BLOCKER | 1 | 0.2% |
| 🟠 CRITICAL | 305 | 55.0% |
| 🟡 MAJOR | 62 | 11.2% |
| 🔵 MINOR | 187 | 33.7% |
| ⚪ INFO | 0 | 0.0% |

## 问题类型分布

| 类型 | 数量 | 占比 |
|------|------|------|
| Bug | 0 | 0.0% |
| Vulnerability | 1 | 0.2% |
| Code Smell | 554 | 99.8% |

## Issue 分布（按文件类型）

| 类别 | 数量 |
|------|------|
| Go | 326 |
| SQL | 123 |
| TypeScript | 53 |
| Other | 27 |
| Vue | 12 |
| Python | 10 |
| Shell | 2 |
| JavaScript | 1 |
| YAML | 1 |

## 关键问题

### 🔴 BLOCKER

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `scripts/gen-sonar-report.py` | 189 | `pythonsecurity:S2083` | Change this code to not construct the path from user-controlled data. |

### 🟠 CRITICAL（Top 20）

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `backend/sql/seed_data.sql` | 32 | `plsql:S1192` | Define a constant instead of duplicating this literal 8 times. |
| `backend/sql/seed_data.sql` | 33 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 34 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 35 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 36 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 37 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 38 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 39 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 102 | `plsql:S1192` | Define a constant instead of duplicating this literal 42 times. |
| `backend/sql/seed_data.sql` | 134 | `plsql:S1192` | Define a constant instead of duplicating this literal 14 times. |
| `backend/sql/seed_data.sql` | 135 | `plsql:S1192` | Define a constant instead of duplicating this literal 9 times. |
| `backend/sql/seed_data.sql` | 136 | `plsql:S1192` | Define a constant instead of duplicating this literal 10 times. |
| `backend/sql/seed_data.sql` | 137 | `plsql:S1192` | Define a constant instead of duplicating this literal 8 times. |
| `backend/sql/seed_data.sql` | 232 | `plsql:S1192` | Define a constant instead of duplicating this literal 12 times. |
| `backend/sql/seed_data.sql` | 282 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 282 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 283 | `plsql:S1192` | Define a constant instead of duplicating this literal 19 times. |
| `backend/sql/seed_data.sql` | 284 | `plsql:S1192` | Define a constant instead of duplicating this literal 15 times. |
| `backend/sql/seed_data.sql` | 285 | `plsql:S1192` | Define a constant instead of duplicating this literal 16 times. |
| `backend/sql/seed_data.sql` | 285 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |

## 高频规则（Top 10）

| 规则 | 级别 | 类型 | 命中数 | 占比 |
|------|------|------|--------|------|
| `go:S1186` | 🟠 CRITICAL | CODE_SMELL | 138 | 24.9% |
| `plsql:S1192` | 🟠 CRITICAL | CODE_SMELL | 123 | 22.2% |
| `go:S100` | 🔵 MINOR | CODE_SMELL | 121 | 21.8% |
| `go:S1192` | 🟠 CRITICAL | CODE_SMELL | 38 | 6.8% |
| `powershelldre:S8626` | 🟡 MAJOR | CODE_SMELL | 26 | 4.7% |
| `godre:S8184` | 🔵 MINOR | CODE_SMELL | 17 | 3.1% |
| `typescript:S6853` | 🟡 MAJOR | CODE_SMELL | 14 | 2.5% |
| `typescript:S6759` | 🔵 MINOR | CODE_SMELL | 8 | 1.4% |
| `typescript:S1128` | 🔵 MINOR | CODE_SMELL | 7 | 1.3% |
| `python:S117` | 🔵 MINOR | CODE_SMELL | 6 | 1.1% |
