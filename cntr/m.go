package cntr

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// M Interface for JSON Field of M Table
type M map[string]interface{}

// Value Marshal
//
//goland:noinspection GoMixedReceiverTypes
func (a M) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
//
//goland:noinspection GoMixedReceiverTypes
func (a *M) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// CleanZeroBytes Clean all string fields with \u0000
//
//goland:noinspection GoMixedReceiverTypes
func (a M) CleanZeroBytes() {
	for k, v := range a {
		if s, ok := v.(string); ok {
			b := []byte(s)
			detected := false
			for i, c := range b {
				if c == 0 {
					b[i] = ' '
					detected = true
				}
			}
			if detected {
				(a)[k] = string(b)
			}
		}

		// try to clean sub map
		if m, ok := v.(map[string]interface{}); ok {
			if m == nil || len(m) == 0 {
				continue
			}
			subMap := M(m)
			subMap.CleanZeroBytes()
		}
	}
}

// GetString Get string value from M
//
//goland:noinspection GoMixedReceiverTypes
func (a M) GetString(key string) string {
	if v, ok := (a)[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetInt Get int value from M
//
//goland:noinspection GoMixedReceiverTypes
func (a M) GetInt(key string) int {
	if v, ok := (a)[key]; ok {
		if s, ok := v.(int); ok {
			return s
		}
	}
	return 0
}

// GetInt64 Get int64 value from M
//
//goland:noinspection GoMixedReceiverTypes
func (a M) GetInt64(key string) int64 {
	if v, ok := (a)[key]; ok {
		if s, ok := v.(int64); ok {
			return s
		}
	}
	return 0
}

// GetFloat64 Get float64 value from M
//
//goland:noinspection GoMixedReceiverTypes
func (a M) GetFloat64(key string) float64 {
	if v, ok := (a)[key]; ok {
		if s, ok := v.(float64); ok {
			return s
		}
	}
	return 0
}

// GetBool Get bool value from M
//
//goland:noinspection GoMixedReceiverTypes
func (a M) GetBool(key string) bool {
	if v, ok := (a)[key]; ok {
		if s, ok := v.(bool); ok {
			return s
		}
	}
	return false
}

// GetM Get M value from M
//
//goland:noinspection GoMixedReceiverTypes
func (a M) GetM(key string) M {
	if v, ok := (a)[key]; ok {
		if s, ok := v.(M); ok {
			return s
		}
		// try to convert to M
		if s, ok := v.(map[string]interface{}); ok {
			return M(s)
		}
	}
	return nil
}
