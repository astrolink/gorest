// Tideland Go REST Server Library - JSON Web Token - Claims
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package jwt

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tideland/golib/errors"
)

//--------------------
// CLAIMS
//--------------------

// Claims contains the claims of a token payload. The type
// also provides getters and setters for the reserved claims.
type Claims map[string]interface{}

// NewClaims returns an empty set of claims.
func NewClaims() Claims {
	return Claims{}
}

// Len returns the number of entries in the claims.
func (c Claims) Len() int {
	if c == nil {
		return 0
	}
	return len(c)
}

// Get retrieves a value from the claims.
func (c Claims) Get(key string) (interface{}, bool) {
	if c == nil {
		return nil, false
	}
	value, ok := c[key]
	return value, ok
}

// GetString retrieves a string value. If it is no string it
// will be converted into a string.
func (c Claims) GetString(key string) (string, bool) {
	value, ok := c.Get(key)
	if !ok {
		return "", false
	}
	if str, ok := value.(string); ok {
		return str, true
	}
	return fmt.Sprintf("%v", value), true
}

// GetInt retrieves an integer value.
func (c Claims) GetInt(key string) (int, bool) {
	value, ok := c.Get(key)
	if !ok {
		return 0, false
	}
	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return int(i), true
		}
	}
	return 0, false
}

// GetFloat64 retrieves a float value.
func (c Claims) GetFloat64(key string) (float64, bool) {
	value, ok := c.Get(key)
	if !ok {
		return 0, false
	}
	switch v := value.(type) {
	case int:
		return float64(v), true
	case float64:
		return v, true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0.0, false
}

// GetTime retrieves a time value. Int, int32, int64,
// and float64 are valid types for the conversion. In case
// a string it is interpreted as RFC 3339 formatted time.
func (c Claims) GetTime(key string) (time.Time, bool) {
	value, ok := c.Get(key)
	if !ok {
		return time.Time{}, false
	}
	switch v := value.(type) {
	case int:
		return time.Unix(int64(v), 0), true
	case int32:
		return time.Unix(int64(v), 0), true
	case int64:
		return time.Unix(v, 0), true
	case float64:
		return time.Unix(int64(v), 0), true
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Time{}, false
		}
		return t, true
	default:
		return time.Time{}, false
	}
}

// Set sets a value in the claims. It returns a potential
// old value.
func (c Claims) Set(key string, value interface{}) interface{} {
	if c == nil {
		return nil
	}
	old, _ := c.Get(key)
	c[key] = value
	return old
}

// SetTime sets a time value in the claims. It returns a
// potential old value.
func (c Claims) SetTime(key string, t time.Time) time.Time {
	old, _ := c.GetTime(key)
	c.Set(key, t.Unix())
	return old
}

// Delete deletes a value from the claims. It returns a potential
// old value.
func (c Claims) Delete(key string) interface{} {
	old, _ := c.Get(key)
	delete(c, key)
	return old
}

// Contains checks if the claims contain a given key.
func (c Claims) Contains(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// Audience retrieves the reserved "aud" claim.
func (c Claims) Audience() ([]string, bool) {
	// Little helper to cast interfaces to strings.
	makeStrings := func(auds ...interface{}) ([]string, bool) {
		if len(auds) == 0 {
			return nil, false
		}
		strs := make([]string, len(auds))
		for i, aud := range auds {
			str, ok := aud.(string)
			if !ok {
				return nil, false
			}
			strs[i] = str
		}
		return strs, true
	}
	// Now retrieve the audience claim.
	value, ok := c.Get("aud")
	if !ok {
		return nil, false
	}
	switch aud := value.(type) {
	case []string:
		return aud, true
	case string:
		return []string{aud}, true
	case []interface{}:
		return makeStrings(aud...)
	case interface{}:
		return makeStrings(aud)
	}
	return nil, false
}

// SetAudience sets the reserved "aud" claim. It returns a
// potential old value.
func (c Claims) SetAudience(auds ...string) []string {
	old, _ := c.Audience()
	switch len(auds) {
	case 0:
		c.Delete("aud")
	case 1:
		c.Set("aud", auds[0])
	default:
		c.Set("aud", auds)
	}
	return old
}

// DeleteAudience deletes the reserved "aud" claim. It returns
// a potential old value.
func (c Claims) DeleteAudience() []string {
	old, _ := c.Audience()
	c.Delete("aud")
	return old
}

// Expiration retrieves the reserved "exp" claim.
func (c Claims) Expiration() (time.Time, bool) {
	return c.GetTime("exp")
}

// SetExpiration sets the reserved "exp" claim. It returns a
// potential old value.
func (c Claims) SetExpiration(t time.Time) time.Time {
	return c.SetTime("exp", t)
}

// DeleteExpiration deletes the reserved "exp" claim. It returns
// a potential old value.
func (c Claims) DeleteExpiration() time.Time {
	old, _ := c.Expiration()
	c.Delete("exp")
	return old
}

// IssuedAt retrieves the reserved "iat" claim.
func (c Claims) IssuedAt() (time.Time, bool) {
	return c.GetTime("iat")
}

// SetIssuedAt sets the reserved "iat" claim. It returns a
// potential old value.
func (c Claims) SetIssuedAt(t time.Time) time.Time {
	return c.SetTime("iat", t)
}

// DeleteIssuedAt deletes the reserved "iat" claim. It returns
// a potential old value.
func (c Claims) DeleteIssuedAt() time.Time {
	old, _ := c.IssuedAt()
	c.Delete("iat")
	return old
}

// Issuer retrieves the reserved "iss" claim.
func (c Claims) Issuer() (string, bool) {
	return c.GetString("iss")
}

// SetIssuer sets the reserved "iss" claim. It returns a
// potential old value.
func (c Claims) SetIssuer(issuer string) string {
	old, _ := c.GetString("iss")
	c.Set("iss", issuer)
	return old
}

// DeleteIssuer deletes the reserved "iss" claim. It returns
// a potential old value.
func (c Claims) DeleteIssuer() string {
	old, _ := c.Issuer()
	c.Delete("iss")
	return old
}

// Identifier retrieves the reserved "jti" claim.
func (c Claims) Identifier() (string, bool) {
	return c.GetString("jti")
}

// SetIdentifier sets the reserved "jti" claim. It returns a
// potential old value.
func (c Claims) SetIdentifier(id string) string {
	old, _ := c.GetString("jti")
	c.Set("jti", id)
	return old
}

// DeleteIdentifier deletes the reserved "jti" claim. It returns
// a potential old value.
func (c Claims) DeleteIdentifier() string {
	old, _ := c.Identifier()
	c.Delete("jti")
	return old
}

// NotBefore retrieves the reserved "nbf" claim.
func (c Claims) NotBefore() (time.Time, bool) {
	return c.GetTime("nbf")
}

// SetNotBefore sets the reserved "nbf" claim. It returns a
// potential old value.
func (c Claims) SetNotBefore(t time.Time) time.Time {
	return c.SetTime("nbf", t)
}

// DeleteNotBefore deletes the reserved "nbf" claim. It returns
// a potential old value.
func (c Claims) DeleteNotBefore() time.Time {
	old, _ := c.NotBefore()
	c.Delete("nbf")
	return old
}

// Subject retrieves the reserved "sub" claim.
func (c Claims) Subject() (string, bool) {
	return c.GetString("sub")
}

// SetSubject sets the reserved "sub" claim. It returns a
// potential old value.
func (c Claims) SetSubject(subject string) string {
	old, _ := c.GetString("sub")
	c.Set("sub", subject)
	return old
}

// DeleteSubject deletes the reserved "sub" claim. It returns
// a potential old value.
func (c Claims) DeleteSubject() string {
	old, _ := c.Subject()
	c.Delete("sub")
	return old
}

// MarshalJSON implements the json.Marshaller interface
// even for nil or empty claims.
func (c Claims) MarshalJSON() ([]byte, error) {
	if c.Len() == 0 {
		return []byte("{}"), nil
	}
	b, err := json.Marshal(map[string]interface{}(c))
	if err != nil {
		return nil, errors.Annotate(err, ErrJSONMarshalling, errorMessages)
	}
	return b, nil
}

// MarshalJSON implements the json.Marshaller interface.
func (c *Claims) UnmarshalJSON(b []byte) error {
	if b == nil {
		return nil
	}
	raw := map[string]interface{}(*c)
	if err := json.Unmarshal(b, &raw); err != nil {
		return errors.Annotate(err, ErrJSONUnmarshalling, errorMessages)
	}
	*c = Claims(raw)
	return nil
}

// EOF
