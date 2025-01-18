package config

type Jwt struct {
	Secret                  string `mapstructure:"secret" json:"secret" yaml:"secret"`
	JwtTtl                  int64  `mapstructure:"jwt_ttl" json:"jwt_ttl" yaml:"jwt_ttl"`
	JwtBlacklistGracePeriod int64  `mapstructure:"jwt_blacklist_grace_period" json:"jwt_blacklist_grace_period" yaml:"jwt_blacklist_grace_period"`
	RefreshGracePeriod      int64  `mapstructure:"refresh_grace_period" json:"refresh_grace_period" yaml:"refresh_grace_period"` // token 自动刷新宽限时间（秒）
}
