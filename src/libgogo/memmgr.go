// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo memory manager functions
//

package libgogo

func Brk(brk uint64) uint64;

func GetBrk() uint64;

func TestMem(address uint64) uint64;
