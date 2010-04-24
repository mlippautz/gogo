// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo stack functions (ASM)
//

TEXT ·ToUint64FromUint64Ptr(SB),$0-16 //ToUint64FromUint64Ptr: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move address from BX to return value (return value after one parameter => SP+2*64bit)
  RET
