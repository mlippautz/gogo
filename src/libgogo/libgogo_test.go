
package libgogo

import "testing"
import "./_obj/libgogo"

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

    len = libgogo.StringLength("+-*/&|!{}[].,;=<>");
    if len != 17 {
		t.Fatalf("libgogo.StringLength() = %d, want 17", len)
    } 
}

