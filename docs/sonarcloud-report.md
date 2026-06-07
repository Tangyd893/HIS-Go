# SonarCloud 代码质量报告

> 2026-06-07

## 概览

- **项目**: HIS-Go
- **Project Key**: Tangyd893_HIS-Go
- **质量门**: ⚪ NOT ANALYZED
- **代码行**: 38,213
- **未解决 Issue**: 551
- **技术债务**: 42h 34min

## 指标

| 指标 | 值 |
|------|-----|
| Bug | 0 |
| 漏洞 | 5 |
| 代码异味 | 546 |
| 重复率 | 6.0% |
| 认知复杂度 | 3,405 |
| 测试覆盖率 | N/A |

## 严重级别分布

| 级别 | 数量 | 占比 |
|------|------|------|
| 🔴 BLOCKER | 4 | 0.7% |
| 🟠 CRITICAL | 240 | 43.6% |
| 🟡 MAJOR | 120 | 21.8% |
| 🔵 MINOR | 187 | 33.9% |
| ⚪ INFO | 0 | 0% |

## 问题类型分布

| 类型 | 数量 | 占比 |
|------|------|------|
| Bug | 0 | 0% |
| Vulnerability | 5 | 0.9% |
| Code Smell | 546 | 99.1% |

## Issue 分布（按文件类型）

| 类别 | 数量 |
|------|------|
| Go | 270 |
| TypeScript | 73 |
| Shell | 63 |
| SQL | 47 |
| YAML | 19 |
| CSS | 3 |
| Other | 76 |

## 关键问题

### 🔴 BLOCKER

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `backend/sql/seed_data.sql` | 18 | `secrets:S8215` | Make sure this bcrypt password hash gets revoked, changed, and removed from the code. |
| `backend/sql/seed_data.sql` | 19 | `secrets:S8215` | Make sure this bcrypt password hash gets revoked, changed, and removed from the code. |
| `backend/sql/seed_data.sql` | 20 | `secrets:S8215` | Make sure this bcrypt password hash gets revoked, changed, and removed from the code. |
| `backend/sql/seed_data.sql` | 21 | `secrets:S8215` | Make sure this bcrypt password hash gets revoked, changed, and removed from the code. |

### 🟠 CRITICAL（Top 20）

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `backend/sql/seed_data.sql` | 18 | `plsql:S1192` | Define a constant instead of duplicating this literal 14 times. |
| `backend/sql/seed_data.sql` | 18 | `plsql:S1192` | Define a constant instead of duplicating this literal 33 times. |
| `backend/sql/seed_data.sql` | 19 | `plsql:S1192` | Define a constant instead of duplicating this literal 7 times. |
| `backend/sql/seed_data.sql` | 28 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |
| `backend/sql/seed_data.sql` | 29 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |
| `backend/sql/seed_data.sql` | 35 | `plsql:S1192` | Define a constant instead of duplicating this literal 11 times. |
| `backend/sql/seed_data.sql` | 36 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 92 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 92 | `plsql:S1192` | Define a constant instead of duplicating this literal 37 times. |
| `backend/sql/seed_data.sql` | 94 | `plsql:S1192` | Define a constant instead of duplicating this literal 14 times. |
| `backend/sql/seed_data.sql` | 112 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 112 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 112 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 113 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 113 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 113 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 114 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 116 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 125 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 125 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |

## 高频规则（Top 10）

| 规则 | 级别 | 类型 | 命中数 | 占比 |
|------|------|------|--------|------|
| `go:S1186` | 🟠 CRITICAL | CODE_SMELL | 138 | 25.0% |
| `go:S100` | 🔵 MINOR | CODE_SMELL | 121 | 22.0% |
| `plsql:S1192` | 🟠 CRITICAL | CODE_SMELL | 46 | 8.3% |
| `go:S1192` | 🟠 CRITICAL | CODE_SMELL | 45 | 8.2% |
| `shelldre:S7688` | 🟡 MAJOR | CODE_SMELL | 39 | 7.1% |
| `typescript:S1128` | 🟡 MAJOR | CODE_SMELL | 21 | 3.8% |
| `kubernetes:S6596` | 🟡 MAJOR | CODE_SMELL | 19 | 3.4% |
| `shelldre:S7679` | 🟡 MAJOR | CODE_SMELL | 18 | 3.3% |
| `typescript:S6853` | 🟡 MAJOR | CODE_SMELL | 14 | 2.5% |
| `godre:S8184` | 🔵 MINOR | CODE_SMELL | 13 | 2.4% |
