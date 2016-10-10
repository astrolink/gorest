// Tideland Go REST Server Library - JSON Web Token
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// The Tideland Go REST Server Library jwt provides the generation,
// verification, and analyzing of JSON Web Tokens.
package jwt

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/golib/version"
)

//--------------------
// VERSION
//--------------------

// Version returns the version of the JSON Web Token package.
func Version() version.Version {
	return version.New(2, 4, 2)
}

// EOF
