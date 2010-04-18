// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package libgogo

import "testing"
import "./_obj/libgogo"

//
// Perform some basic checks on the libgogo.StringLength function
// This function is not capable of handling utf strings!
//
func TestStringLength(t *testing.T) { 
    var len uint64;

    len = libgogo.StringLength("test");
    if len != 4 {
		t.Fatalf("libgogo.StringLength() = %d, want 4", len)
    }

    len = libgogo.StringLength("");
    if len != 0 {
		t.Fatalf("libgogo.StringLength() = %d, want 0", len)
    }

    len = libgogo.StringLength("1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890");
    if len != 100 {
		t.Fatalf("libgogo.StringLength() = %d, want 100", len)
    }

    len = libgogo.StringLength("+-*/><=!&.");
    if len != 10 {
		t.Fatalf("libgogo.StringLength() = %d, want 10", len)
    } 
}

//
// Perform basic comparison checks on libgogo.StringCompare function
//
func TestStringCompare(t *testing.T)  {
    var cmpFlag uint64;
    
    cmpFlag = libgogo.StringCompare("a","a");
    if cmpFlag != 0 {
        t.Fatalf("libgogo.StringCompare(\"a\") != 0, received %d", cmpFlag);
    }

    cmpFlag = libgogo.StringCompare("these strings","are not equal");
    if cmpFlag == 0 {
        t.Fatalf("libgogo.StringCompare(\"these strings\",\"are not equal\") == 0");
    }
}

