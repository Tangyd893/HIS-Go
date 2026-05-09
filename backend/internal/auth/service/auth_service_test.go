package service

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/auth/model"
	"his-go/internal/auth/repository"
	"his-go/pkg/redis"
	jwtpkg "his-go/pkg/security/jwt"
)

func createTestUser(t *testing.T, db *gorm.DB, username, password, role string) *model.AuthUser {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt 加密失败: %v", err)
	}
	user := &model.AuthUser{
		ID:       username + "-id",
		Username: username,
		Password: string(hashed),
		Role:     role,
		Status:   1,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	return user
}

func createTestRole(t *testing.T, db *gorm.DB, roleID, roleName, roleCode string) {
	t.Helper()
	role := model.Role{
		ID:       roleID,
		RoleName: roleName,
		RoleCode: roleCode,
		Status:   1,
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("创建测试角色失败: %v", err)
	}
}

func createTestPermission(t *testing.T, db *gorm.DB, permID, permName, permCode string) {
	t.Helper()
	perm := model.Permission{
		ID:       permID,
		PermName: permName,
		PermCode: permCode,
	}
	if err := db.Create(&perm).Error; err != nil {
		t.Fatalf("创建测试权限失败: %v", err)
	}
}

func createTestRolePermission(t *testing.T, db *gorm.DB, roleID, permID string) {
	t.Helper()
	rp := model.RolePermission{
		RoleID: roleID,
		PermID: permID,
	}
	if err := db.Create(&rp).Error; err != nil {
		t.Fatalf("创建角色权限关联失败: %v", err)
	}
}

func newTestAuthService(t *testing.T) (*AuthService, *gorm.DB, *jwtpkg.JWTService, *redis.Client, *miniredis.Miniredis) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	// SQLite 不支持 PostgreSQL 的 gen_random_uuid()，手动建表
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			real_name TEXT,
			phone TEXT,
			email TEXT,
			avatar TEXT,
			role TEXT DEFAULT 'doctor',
			dept_id TEXT,
			status INTEGER DEFAULT 1,
			last_login_time DATETIME,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS roles (
			id TEXT PRIMARY KEY,
			role_name TEXT NOT NULL,
			role_code TEXT NOT NULL UNIQUE,
			description TEXT,
			status INTEGER DEFAULT 1,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS permissions (
			id TEXT PRIMARY KEY,
			perm_name TEXT NOT NULL,
			perm_code TEXT NOT NULL UNIQUE,
			parent_id TEXT,
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS role_permissions (
			role_id TEXT NOT NULL,
			perm_id TEXT NOT NULL,
			PRIMARY KEY (role_id, perm_id)
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v\nSQL: %s", err, ddl)
		}
	}

	mr := miniredis.RunT(t)
	rdb := &redis.Client{
		Client: goredis.NewClient(&goredis.Options{Addr: mr.Addr()}),
	}

	jwtSvc := jwtpkg.NewSimpleJWTService("test-auth-secret", 24)
	repo := repository.NewAuthRepository(db)
	svc := NewAuthService(repo, jwtSvc, rdb)

	return svc, db, jwtSvc, rdb, mr
}

func TestAuthService_Login_Success(t *testing.T) {
	svc, db, jwtSvc, _, mr := newTestAuthService(t)
	defer mr.Close()

	user := createTestUser(t, db, "doctor1", "pass123", "doctor")
	createTestRole(t, db, "role-1", "医生", "doctor")
	createTestPermission(t, db, "perm-1", "患者查看", "patient:read")
	createTestPermission(t, db, "perm-2", "处方开具", "prescription:write")
	createTestRolePermission(t, db, "role-1", "perm-1")
	createTestRolePermission(t, db, "role-1", "perm-2")

	result, err := svc.Login("doctor1", "pass123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if result.Token == "" {
		t.Error("期望返回非空 Token")
	}
	if result.RefreshToken == "" {
		t.Error("期望返回非空 RefreshToken")
	}
	if result.UserInfo.Username != "doctor1" {
		t.Errorf("期望 Username='doctor1'，实际=%s", result.UserInfo.Username)
	}
	if result.UserInfo.Role != "doctor" {
		t.Errorf("期望 Role='doctor'，实际=%s", result.UserInfo.Role)
	}
	if len(result.UserInfo.Perms) != 2 {
		t.Errorf("期望 2 个权限，实际=%d", len(result.UserInfo.Perms))
	}
	if result.ExpiresIn <= 0 {
		t.Errorf("期望 ExpiresIn > 0，实际=%d", result.ExpiresIn)
	}

	parsed, err := jwtSvc.ParseToken(result.Token)
	if err != nil {
		t.Fatalf("解析返回的 Token 失败: %v", err)
	}
	if parsed.UserID != user.ID {
		t.Errorf("期望 UserID=%s，实际=%s", user.ID, parsed.UserID)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc, db, _, _, mr := newTestAuthService(t)
	defer mr.Close()

	createTestUser(t, db, "doctor2", "correct-pass", "doctor")

	_, err := svc.Login("doctor2", "wrong-pass")
	if err == nil {
		t.Error("期望密码错误时返回错误")
	}
	if err.Error() != "密码错误" {
		t.Errorf("期望 '密码错误'，实际=%s", err.Error())
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc, _, _, _, mr := newTestAuthService(t)
	defer mr.Close()

	_, err := svc.Login("nonexistent", "any-pass")
	if err == nil {
		t.Error("期望用户不存在时返回错误")
	}
	if err.Error() != "用户不存在" {
		t.Errorf("期望 '用户不存在'，实际=%s", err.Error())
	}
}

func TestAuthService_Logout(t *testing.T) {
	svc, db, _, rdb, mr := newTestAuthService(t)
	defer mr.Close()

	createTestUser(t, db, "doctor3", "pass123", "doctor")

	result, err := svc.Login("doctor3", "pass123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	err = svc.Logout(result.UserInfo.UserID)
	if err != nil {
		t.Fatalf("登出失败: %v", err)
	}

	ctx := context.Background()
	_, err = rdb.Get(ctx, "auth:token:"+result.UserInfo.UserID)
	if err == nil {
		t.Error("期望登出后 token 被删除")
	}
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	svc, db, _, _, mr := newTestAuthService(t)
	defer mr.Close()

	createTestUser(t, db, "doctor4", "pass123", "doctor")

	result, err := svc.Login("doctor4", "pass123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	claims, err := svc.ValidateToken(result.Token)
	if err != nil {
		t.Fatalf("验证 Token 失败: %v", err)
	}
	if claims.UserID != result.UserInfo.UserID {
		t.Errorf("期望 UserID=%s，实际=%s", result.UserInfo.UserID, claims.UserID)
	}
}

func TestAuthService_ValidateToken_NotInRedis(t *testing.T) {
	svc, _, jwtSvc, _, mr := newTestAuthService(t)
	defer mr.Close()

	claims := &jwtpkg.Claims{
		UserID:   "doctor5-id",
		Username: "doctor5",
		Role:     "doctor",
	}
	token, err := jwtSvc.GenerateToken(claims)
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}

	_, err = svc.ValidateToken(token)
	if err == nil {
		t.Error("期望 Token 不在 Redis 时验证失败")
	}
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	svc, db, _, _, mr := newTestAuthService(t)
	defer mr.Close()

	createTestUser(t, db, "doctor6", "pass123", "doctor")
	createTestRole(t, db, "role-r", "医生", "doctor")
	createTestPermission(t, db, "perm-r", "患者查看", "patient:read")
	createTestRolePermission(t, db, "role-r", "perm-r")

	result, err := svc.Login("doctor6", "pass123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	newResult, err := svc.RefreshToken(result.RefreshToken)
	if err != nil {
		t.Fatalf("刷新 Token 失败: %v", err)
	}
	if newResult.Token == "" {
		t.Error("期望返回非空 Token")
	}
	if newResult.UserInfo.Username != "doctor6" {
		t.Errorf("期望 Username='doctor6'，实际=%s", newResult.UserInfo.Username)
	}
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	svc, _, _, _, mr := newTestAuthService(t)
	defer mr.Close()

	_, err := svc.RefreshToken("invalid-token-string")
	if err == nil {
		t.Error("期望无效 Token 时返回错误")
	}
}

func TestAuthService_RefreshToken_NotInRedis(t *testing.T) {
	svc, _, jwtSvc, _, mr := newTestAuthService(t)
	defer mr.Close()

	claims := &jwtpkg.Claims{
		UserID:   "some-user",
		Username: "some-user",
		Role:     "doctor",
	}
	token, err := jwtSvc.GenerateToken(claims)
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}

	_, err = svc.RefreshToken(token)
	if err == nil {
		t.Error("期望 Token 不在 Redis 时返回错误")
	}
}

func TestAuthService_Login_NoPermissions(t *testing.T) {
	svc, db, _, _, mr := newTestAuthService(t)
	defer mr.Close()

	createTestUser(t, db, "norole", "pass123", "unknown-role")

	result, err := svc.Login("norole", "pass123")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if len(result.UserInfo.Perms) != 0 {
		t.Errorf("期望无权限，实际=%d 个", len(result.UserInfo.Perms))
	}
}
