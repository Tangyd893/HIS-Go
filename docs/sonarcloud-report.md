# SonarCloud 代码质量报告

> 2026-06-08

## 概览

- **项目**: HIS-Go
- **Project Key**: Tangyd893_HIS-Go
- **质量门**: ❌ ERROR
- **代码行**: 38897
- **未解决 Issue**: 458
- **技术债务**: 38h 24min

## 指标

| 指标 | 值 |
|------|-----|
| Bug | 0 |
| 漏洞 | 1 |
| 代码异味 | 457 |
| 重复率 | 6.0% |
| 认知复杂度 | 3646 |
| 测试覆盖率 | N/A% |

## 严重级别分布

| 级别 | 数量 | 占比 |
|------|------|------|
| 🔴 BLOCKER | 0 | 0.0% |
| 🟠 CRITICAL | 228 | 49.8% |
| 🟡 MAJOR | 59 | 12.9% |
| 🔵 MINOR | 171 | 37.3% |
| ⚪ INFO | 0 | 0.0% |

## 问题类型分布

| 类型 | 数量 | 占比 |
|------|------|------|
| Bug | 0 | 0.0% |
| Vulnerability | 1 | 0.2% |
| Code Smell | 457 | 99.8% |

## Issue 分布（按文件类型）

| 类别 | 数量 |
|------|------|
| Go | 328 |
| TypeScript | 52 |
| SQL | 42 |
| Other | 28 |
| Vue | 5 |
| Shell | 2 |
| YAML | 1 |

## 关键问题

### 🔴 BLOCKER

无

### 🟠 CRITICAL（Top 20）

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `backend/internal/auth/service/auth_service.go` | 95 | `go:S1192` | Define a constant instead of duplicating this literal "auth:token:" 4 times. |
| `backend/internal/registration/repository/registration_repo.go` | 192 | `go:S1192` | Define a constant instead of duplicating this literal "queue:" 3 times. |
| `backend/internal/outpatient/assistant/department_resolver.go` | 71 | `go:S3776` | Refactor this method to reduce its Cognitive Complexity from 21 to the 15 allowed. |
| `backend/internal/outpatient/assistant/retriever.go` | 78 | `go:S3776` | Refactor this method to reduce its Cognitive Complexity from 21 to the 15 allowed. |
| `backend/internal/outpatient/assistant/retriever.go` | 134 | `go:S3776` | Refactor this method to reduce its Cognitive Complexity from 18 to the 15 allowed. |
| `backend/internal/outpatient/assistant/service.go` | 49 | `go:S3776` | Refactor this method to reduce its Cognitive Complexity from 34 to the 15 allowed. |
| `backend/sql/seed_data.sql` | 29 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |
| `backend/sql/seed_data.sql` | 30 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |
| `backend/sql/seed_data.sql` | 36 | `plsql:S1192` | Define a constant instead of duplicating this literal 11 times. |
| `backend/sql/seed_data.sql` | 37 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 93 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 93 | `plsql:S1192` | Define a constant instead of duplicating this literal 37 times. |
| `backend/sql/seed_data.sql` | 95 | `plsql:S1192` | Define a constant instead of duplicating this literal 14 times. |
| `backend/sql/seed_data.sql` | 113 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 113 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 113 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 114 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 114 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `backend/sql/seed_data.sql` | 114 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |
| `backend/sql/seed_data.sql` | 115 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |

## 高频规则（Top 10）

| 规则 | 级别 | 类型 | 命中数 | 占比 |
|------|------|------|--------|------|
| `go:S1186` | 🟠 CRITICAL | CODE_SMELL | 138 | 30.1% |
| `go:S100` | 🔵 MINOR | CODE_SMELL | 121 | 26.4% |
| `plsql:S1192` | 🟠 CRITICAL | CODE_SMELL | 42 | 9.2% |
| `go:S1192` | 🟠 CRITICAL | CODE_SMELL | 41 | 9.0% |
| `powershelldre:S8626` | 🟡 MAJOR | CODE_SMELL | 26 | 5.7% |
| `typescript:S6853` | 🟡 MAJOR | CODE_SMELL | 14 | 3.1% |
| `godre:S8184` | 🔵 MINOR | CODE_SMELL | 13 | 2.8% |
| `typescript:S6759` | 🔵 MINOR | CODE_SMELL | 8 | 1.7% |
| `go:S3776` | 🟠 CRITICAL | CODE_SMELL | 7 | 1.5% |
| `typescript:S1128` | 🔵 MINOR | CODE_SMELL | 7 | 1.5% |
