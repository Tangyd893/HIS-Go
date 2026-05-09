package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

func generateRS256KeyPair(t *testing.T) (privatePEM, publicPEM string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥对失败: %v", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("序列化私钥失败: %v", err)
	}
	privatePEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}))

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("序列化公钥失败: %v", err)
	}
	publicPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}))

	return
}

const testLegacyPrivateKeyPEM = `-----BEGIN PRIVATE KEY-----
MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEAtmmQDKoUCJPtCr7X
ArmsUjKgKKDJNH0E5bn3tOn8qvY6gE39P5+eljKyjPd2HDUiMPDt8X9uuaKe+8rq
zHiApQIDAQABAkAlvrgw9qyIjdt54r1o8fSnWZRsc8DOnKP7yTxpchV3ZnxSvrTZ
d1qZyk+mwUxl8+U5LGjqO95RE1Xz01a/PaktAiEA5/F5+QvDSTmXGoHT3LxA/kTc
liQppeUGYdcG9zECJk8CIQDJVPLew9FiFSN8YvClOaoKy3JooPoHVeQf3H8JjRHg
ywIhAMY54C5yWSIZsAQddL2vvjQREhzXJyj6xSuVJATaw6WNAiAI0RWJt92VhAN3
0QVk1u+hZWNvPY11gMdqtcdCbdEYnQIhAOTLJvTnmn6iPI/+J2xdDxcFPFAIBmU6
QssxtDLRGbtO
-----END PRIVATE KEY-----`

const testLegacyPublicKeyPEM = `-----BEGIN PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALZpkAyqFAiT7Qq+1wK5rFIyoCigyTR9
BOW597Tp/Kr2OoBN/T+fnpYysoz3dhw1IjDw7fF/brminvvK6sx4gKUCAwEAAQ==
-----END PUBLIC KEY-----`

func TestNewSimpleJWTService(t *testing.T) {
	svc := NewSimpleJWTService("test-secret", 1)
	if svc == nil {
		t.Fatal("NewSimpleJWTService 返回 nil")
	}
	if !svc.useHMAC {
		t.Error("期望 useHMAC 为 true")
	}
	if svc.expireHour != 1 {
		t.Errorf("期望 expireHour=1，实际=%d", svc.expireHour)
	}
}

func TestNewSimpleJWTService_EmptySecret(t *testing.T) {
	svc := NewSimpleJWTService("", 2)
	if svc == nil {
		t.Fatal("NewSimpleJWTService 返回 nil")
	}
	if string(svc.hmacSecret) != "his-go-default-secret" {
		t.Error("期望使用默认密钥")
	}
}

func TestNewJWTService_ValidKeys(t *testing.T) {
	priv, pub := generateRS256KeyPair(t)
	svc, err := NewJWTService(priv, pub, 1)
	if err != nil {
		t.Fatalf("创建 JWTService 失败: %v", err)
	}
	if svc.useHMAC {
		t.Error("期望 RSA 模式")
	}
}

func TestNewJWTService_InvalidPrivateKey(t *testing.T) {
	_, pub := generateRS256KeyPair(t)
	_, err := NewJWTService("not-valid-pem", pub, 1)
	if err == nil {
		t.Error("期望无效私钥时返回错误")
	}
}

func TestNewJWTService_InvalidPublicKey(t *testing.T) {
	priv, _ := generateRS256KeyPair(t)
	_, err := NewJWTService(priv, "not-valid-pem", 1)
	if err == nil {
		t.Error("期望无效公钥时返回错误")
	}
}

func TestNewJWTService_NonRSAPrivateKey(t *testing.T) {
	_, pub := generateRS256KeyPair(t)
	ecKey := "-----BEGIN EC PRIVATE KEY-----\nMHQCAQEE..." // 无效的 EC 密钥
	_, err := NewJWTService(ecKey, pub, 1)
	if err == nil {
		t.Error("期望非 RSA 私钥时返回错误")
	}
}

func TestNewVerifyOnlyJWTService_Valid(t *testing.T) {
	_, pub := generateRS256KeyPair(t)
	svc, err := NewVerifyOnlyJWTService(pub)
	if err != nil {
		t.Fatalf("创建 VerifyOnly JWTService 失败: %v", err)
	}
	if svc.privateKey != nil {
		t.Error("VerifyOnly 服务不应持有私钥")
	}
}

func TestNewVerifyOnlyJWTService_InvalidKey(t *testing.T) {
	_, err := NewVerifyOnlyJWTService("invalid")
	if err == nil {
		t.Error("期望无效公钥时返回错误")
	}
}

func TestGenerateAndParseToken_HS256(t *testing.T) {
	svc := NewSimpleJWTService("test-secret", 1)
	claims := &Claims{
		UserID:   "user-001",
		Username: "testuser",
		Role:     "doctor",
		DeptID:   "dept-1",
		Perms:    []string{"patient:read", "prescription:write"},
	}

	token, err := svc.GenerateToken(claims)
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}
	if token == "" {
		t.Fatal("Token 为空")
	}

	parsed, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("解析 Token 失败: %v", err)
	}
	if parsed.UserID != "user-001" {
		t.Errorf("期望 UserID='user-001'，实际=%s", parsed.UserID)
	}
	if parsed.Username != "testuser" {
		t.Errorf("期望 Username='testuser'，实际=%s", parsed.Username)
	}
	if len(parsed.Perms) != 2 {
		t.Errorf("期望 2 个权限，实际=%d", len(parsed.Perms))
	}
}

func TestGenerateAndParseToken_RS256(t *testing.T) {
	priv, pub := generateRS256KeyPair(t)
	svc, err := NewJWTService(priv, pub, 1)
	if err != nil {
		t.Fatalf("创建 JWTService 失败: %v", err)
	}

	claims := &Claims{
		UserID:   "user-002",
		Username: "doctor1",
		Role:     "doctor",
	}

	token, err := svc.GenerateToken(claims)
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}

	parsed, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("解析 Token 失败: %v", err)
	}
	if parsed.UserID != "user-002" {
		t.Errorf("期望 UserID='user-002'，实际=%s", parsed.UserID)
	}
}

func TestParseToken_Expired(t *testing.T) {
	svc := NewSimpleJWTService("test-secret", 0) // 0 小时过期
	claims := &Claims{UserID: "user-1", Username: "test"}
	token, err := svc.GenerateToken(claims)
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	_, err = svc.ParseToken(token)
	if err == nil {
		t.Error("期望过期 Token 时返回错误")
	}
}

func TestParseToken_InvalidTokenString(t *testing.T) {
	svc := NewSimpleJWTService("test-secret", 1)
	_, err := svc.ParseToken("not.a.valid.jwt")
	if err == nil {
		t.Error("期望无效 Token 字符串时返回错误")
	}
}

func TestParseToken_WrongSigningMethod(t *testing.T) {
	// 用 HS256 签发，用 RS256 尝试解析（模拟签名方法不匹配）
	hmacToken := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, &Claims{UserID: "test"})
	hmacStr, _ := hmacToken.SignedString([]byte("test-secret"))

	priv, pub := generateRS256KeyPair(t)
	rsaSvc, _ := NewJWTService(priv, pub, 1)
	_, err := rsaSvc.ParseToken(hmacStr)
	if err == nil {
		t.Error("期望签名方法不匹配时返回错误")
	}
}

func TestParseToken_RS256_WrongSigningMethod(t *testing.T) {
	priv, pub := generateRS256KeyPair(t)
	rsaSvc, _ := NewJWTService(priv, pub, 1)
	claims := &Claims{UserID: "test"}
	rsaToken, err := rsaSvc.GenerateToken(claims)
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}

	hsSvc := NewSimpleJWTService("test-secret", 1)
	_, err = hsSvc.ParseToken(rsaToken)
	if err == nil {
		t.Error("期望 HS256 服务解析 RS256 Token 时返回错误")
	}
}

func TestParseToken_WrongHMACSecret(t *testing.T) {
	svc1 := NewSimpleJWTService("correct-secret", 1)
	svc2 := NewSimpleJWTService("wrong-secret", 1)

	claims := &Claims{UserID: "test"}
	token, _ := svc1.GenerateToken(claims)
	_, err := svc2.ParseToken(token)
	if err == nil {
		t.Error("期望错误密钥时返回错误")
	}
}

func TestParseToken_WrongRSAKey(t *testing.T) {
	priv1, pub1 := generateRS256KeyPair(t)
	svc1, _ := NewJWTService(priv1, pub1, 1)

	priv2, pub2 := generateRS256KeyPair(t)
	svc2, _ := NewJWTService(priv2, pub2, 1)

	claims := &Claims{UserID: "test"}
	token, _ := svc1.GenerateToken(claims)
	_, err := svc2.ParseToken(token)
	if err == nil {
		t.Error("期望不同密钥对时返回错误")
	}
}

func TestRefreshToken(t *testing.T) {
	svc := NewSimpleJWTService("test-secret", 1)
	claims := &Claims{
		UserID:   "user-001",
		Username: "testuser",
	}

	origToken, err := svc.GenerateToken(claims)
	if err != nil {
		t.Fatalf("生成原始 Token 失败: %v", err)
	}

	parsed, err := svc.ParseToken(origToken)
	if err != nil {
		t.Fatalf("解析原始 Token 失败: %v", err)
	}

	refreshToken, err := svc.RefreshToken(parsed)
	if err != nil {
		t.Fatalf("刷新 Token 失败: %v", err)
	}
	if refreshToken == "" {
		t.Fatal("刷新 Token 为空")
	}

	// 验证刷新后的 token 仍可解析且内容一致
	newParsed, err := svc.ParseToken(refreshToken)
	if err != nil {
		t.Fatalf("解析刷新后 Token 失败: %v", err)
	}
	if newParsed.UserID != "user-001" {
		t.Errorf("期望 UserID='user-001'，实际=%s", newParsed.UserID)
	}
}

func TestParseToken_NBFViolation(t *testing.T) {
	svc := NewSimpleJWTService("test-secret", 24)

	claims := &Claims{UserID: "test"}
	// 构造一个 nbf 在未来的 token（不使用 GenerateToken，因为它会覆盖 NBF）
	claims.RegisteredClaims = jwtlib.RegisteredClaims{
		Issuer:    "his-go",
		Subject:   "test",
		ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(2 * time.Hour)),
		IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		NotBefore: jwtlib.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}

	_, err = svc.ParseToken(tokenStr)
	if err == nil {
		t.Error("期望 nbf 在未来的 Token 解析时返回错误")
	}
}

func TestParseToken_MalformedBase64(t *testing.T) {
	svc := NewSimpleJWTService("test-secret", 1)
	_, err := svc.ParseToken("header.!!!invalid!!!.signature")
	if err == nil {
		t.Error("期望格式异常 Token 时返回错误")
	}
}

func TestVerifyOnlyJWTService_CannotGenerate(t *testing.T) {
	_, pub := generateRS256KeyPair(t)
	svc, err := NewVerifyOnlyJWTService(pub)
	if err != nil {
		t.Fatalf("创建 VerifyOnly 服务失败: %v", err)
	}

	claims := &Claims{UserID: "test"}
	_, err = svc.GenerateToken(claims)
	if err == nil {
		t.Error("期望 VerifyOnly 服务生成 Token 时返回错误")
	}
}

func TestParsePrivateKey_InvalidPEM(t *testing.T) {
	_, err := parsePrivateKey("not-valid-pem-content")
	if err == nil {
		t.Error("期望无效 PEM 时返回错误")
	}
}

func TestParsePublicKey_InvalidPEM(t *testing.T) {
	_, err := parsePublicKey("not-valid-pem-content")
	if err == nil {
		t.Error("期望无效 PEM 时返回错误")
	}
}

func TestParsePrivateKey_NonRSAType(t *testing.T) {
	// EC 密钥的 PEM 块
	ecPEM := `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIHv5I8xpH3qYJhG2oY6eXqSmPLQsMNHN8YbYZxD2YqJ8
-----END PRIVATE KEY-----`
	_, err := parsePrivateKey(ecPEM)
	if err == nil {
		t.Error("期望非 RSA 密钥时返回错误")
	}
	if err != nil && !strings.Contains(err.Error(), "类型不是 RSA") {
		t.Errorf("期望 '类型不是 RSA' 错误，实际: %v", err)
	}
}

func TestParsePublicKey_NonRSAType(t *testing.T) {
	ecPub := `-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAGbFGwGknxjb9YwU9JfsmCx8YHLNsPpQDDNOFgNOHK8A=
-----END PUBLIC KEY-----`
	_, err := parsePublicKey(ecPub)
	if err == nil {
		t.Error("期望非 RSA 公钥时返回错误")
	}
	if err != nil && !strings.Contains(err.Error(), "类型不是 RSA") {
		t.Errorf("期望 '类型不是 RSA' 错误，实际: %v", err)
	}
}

func TestParseToken_EmptyString(t *testing.T) {
	svc := NewSimpleJWTService("test-secret", 1)
	_, err := svc.ParseToken("")
	if err == nil {
		t.Error("期望空 Token 字符串时返回错误")
	}
}

func TestClaims_RegisteredClaimsOverride(t *testing.T) {
	svc := NewSimpleJWTService("test-secret", 2)
	claims := &Claims{
		UserID:   "user-003",
		Username: "demo",
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer: "custom-issuer",
		},
	}

	token, err := svc.GenerateToken(claims)
	if err != nil {
		t.Fatalf("生成 Token 失败: %v", err)
	}

	parsed, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("解析 Token 失败: %v", err)
	}

	// GenerateToken 会覆盖 RegisteredClaims，Issuer 应为 "his-go"
	if parsed.Issuer != "his-go" {
		t.Errorf("期望 Issuer='his-go'，实际=%s", parsed.Issuer)
	}
}
