# SonarCloud 代码质量报告

> 2026-06-08

## 概览

- **项目**: HIS-Go
- **Project Key**: Tangyd893_HIS-Go
- **质量门**: ❌ ERROR
- **代码行**: 40455
- **未解决 Issue**: 524
- **技术债务**: 42h 22min

## 指标

| 指标 | 值 |
|------|-----|
| Bug | 0 |
| 漏洞 | 2 |
| 代码异味 | 522 |
| 重复率 | 6.0% |
| 认知复杂度 | 3705 |
| 测试覆盖率 | N/A% |

## 严重级别分布

| 级别 | 数量 | 占比 |
|------|------|------|
| 🔴 BLOCKER | 2 | 0.4% |
| 🟠 CRITICAL | 284 | 54.2% |
| 🟡 MAJOR | 61 | 11.6% |
| 🔵 MINOR | 177 | 33.8% |
| ⚪ INFO | 0 | 0.0% |

## 问题类型分布

| 类型 | 数量 | 占比 |
|------|------|------|
| Bug | 0 | 0.0% |
| Vulnerability | 2 | 0.4% |
| Code Smell | 522 | 99.6% |

## Issue 分布（按文件类型）

| 类别 | 数量 |
|------|------|
| Go | 322 |
| SQL | 104 |
| TypeScript | 53 |
| Other | 27 |
| Vue | 12 |
| Python | 2 |
| Shell | 2 |
| JavaScript | 1 |
| YAML | 1 |

## 关键问题

### 🔴 BLOCKER

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `backend/sql/seed_data_extended.sql` | 44 | `secrets:S8215` | Make sure this bcrypt password hash gets revoked, changed, and removed from the code. |
| `scripts/gen-sonar-report.py` | 189 | `pythonsecurity:S2083` | Change this code to not construct the path from user-controlled data. |

### 🟠 CRITICAL（Top 20）

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `backend/sql/seed_data.sql` | 11 | `plsql:S1192` | Define a constant instead of duplicating this literal 6 times. |
| `backend/sql/seed_data.sql` | 31 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 32 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 115 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 115 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |
| `backend/sql/seed_data.sql` | 116 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 116 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |
| `backend/sql/seed_data.sql` | 116 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |
| `backend/sql/seed_data.sql` | 176 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 259 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 286 | `plsql:S1192` | Define a constant instead of duplicating this literal 8 times. |
| `backend/sql/seed_data.sql` | 288 | `plsql:S1192` | Define a constant instead of duplicating this literal 8 times. |
| `backend/sql/seed_data.sql` | 305 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 307 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 309 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 323 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 324 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 370 | `plsql:S1192` | Define a constant instead of duplicating this literal 8 times. |
| `backend/sql/seed_data.sql` | 370 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |
| `backend/sql/seed_data.sql` | 371 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |

## 高频规则（Top 10）

| 规则 | 级别 | 类型 | 命中数 | 占比 |
|------|------|------|--------|------|
| `go:S1186` | 🟠 CRITICAL | CODE_SMELL | 138 | 26.3% |
| `go:S100` | 🔵 MINOR | CODE_SMELL | 121 | 23.1% |
| `plsql:S1192` | 🟠 CRITICAL | CODE_SMELL | 103 | 19.7% |
| `go:S1192` | 🟠 CRITICAL | CODE_SMELL | 38 | 7.3% |
| `powershelldre:S8626` | 🟡 MAJOR | CODE_SMELL | 26 | 5.0% |
| `typescript:S6853` | 🟡 MAJOR | CODE_SMELL | 14 | 2.7% |
| `godre:S8184` | 🔵 MINOR | CODE_SMELL | 13 | 2.5% |
| `typescript:S6759` | 🔵 MINOR | CODE_SMELL | 8 | 1.5% |
| `typescript:S1128` | 🔵 MINOR | CODE_SMELL | 7 | 1.3% |
| `typescript:S7773` | 🔵 MINOR | CODE_SMELL | 5 | 1.0% |
