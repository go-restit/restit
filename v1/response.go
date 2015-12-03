// Copyright (c) 2014 Yeung Shu Hung (Koala Yeung)
//
//  This file is part of RESTit.
//
//  RESTit is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  RESTit is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  Use of this source code is governed by the GPL v3 license. A copy
//  of the licence can be found in the LICENSE.md file along with RESTit.
//  If not, see <http://www.gnu.org/licenses/>.

package restit

// Response needed to fulfill this interface
// in order to be tested by RESTit
type Response interface {

	// count the number of items in the result
	Count() int

	// test if the nth item is valid
	NthValid(int) error

	// test get the nth item
	GetNth(int) (interface{}, error)

	// test if the given 2 items matches
	Match(interface{}, interface{}) error

	// reset the result, prevent the problem of null result
	Reset()
}
