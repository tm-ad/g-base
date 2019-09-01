// Package i18n provides a simple implementation of i18n interface
// And a basic implementation (no specific multilingual management logic is implemented, only default values are entered)
// All outputs in base follow i18n output requirements
package i18n

import "fmt"

// _i18nProvider provider
var _i18nProvider Provider

// _locale default locale
var _locale = "en"

// Provider is the implementation definition interface providing i18n
type Provider interface {
	L(locale, catalog, key, defaultFormat string, args ...interface{}) string
	AddCatalog(catalog string) error
	AddKey(catalog, key, format string) error
}

// L returns text processed in the default locale
func L(catalog, key, pattern string, args ...interface{}) string {
	return LocalizeWithLocale(_locale, catalog, key, pattern, args...)
}

// LocalizeWithLocale returns text processed in the specified locale
func LocalizeWithLocale(locale, catalog, key, pattern string, args ...interface{}) string {
	if _i18nProvider == nil {
		return fmt.Sprintf(pattern, args...)
	}
	return _i18nProvider.L(locale, catalog, key, pattern, args...)
}

// SetLocale sets default locale
func SetLocale(locale string) {
	_locale = locale
}

// SetProvider sets a implementation of i18n provider
func SetProvider(provider Provider) {
	_i18nProvider = provider
}
