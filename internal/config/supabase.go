package config

import (
	"github.com/febry3/gamingin/internal/infra/storage"
	"github.com/spf13/viper"
)

func NewSupabaseConfig(config *viper.Viper) storage.SupabaseConfig {
	return storage.SupabaseConfig{
		ProjectRef: config.GetString("SUPABASE_PROJECT_REF"),
		ApiKey:     config.GetString("SUPABASE_API_KEY"),
	}
}
