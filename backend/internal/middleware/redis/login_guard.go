package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// 登录失败自动封禁参数：15 分钟窗口内失败达 10 次封禁 1 小时（到点自动解封）。
// 共享出口 IP（NAT）下可能误伤正常用户，因此自动封禁带 TTL 而非永久
const (
	LoginFailThreshold = 10
	loginFailWindow    = 15 * time.Minute
	LoginBanTTL        = time.Hour
)

const loginFailKeyPrefix = "security:login_fail:"

// incrLoginFail 原子执行 INCR 与首次 EXPIRE，避免 key 因 TTL 缺失而永不过期
var incrLoginFail = redis.NewScript(`
local v = redis.call('INCR', KEYS[1])
if v == 1 then
	redis.call('EXPIRE', KEYS[1], ARGV[1])
end
return v
`)

// RecordLoginFailure 记录一次登录失败，返回当前窗口内累计失败次数
func RecordLoginFailure(client *redis.Client, ip string) (int64, error) {
	return incrLoginFail.Run(context.Background(), client,
		[]string{loginFailKeyPrefix + ip}, int64(loginFailWindow.Seconds())).Int64()
}

// ResetLoginFailures 登录成功后清除该 IP 的失败计数
func ResetLoginFailures(client *redis.Client, ip string) {
	client.Del(context.Background(), loginFailKeyPrefix+ip)
}

// AutoBanIfAbusive 失败次数达到阈值时自动封禁该 IP（带 TTL），返回是否触发
func AutoBanIfAbusive(client *redis.Client, ip string, fails int64) bool {
	if fails < LoginFailThreshold {
		return false
	}
	return BanIP(client, ip, LoginBanTTL) == nil
}
