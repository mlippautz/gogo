// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

TEXT ·Uint64ToObjectDescPtr(SB),0,$0-0 //Uint64ToObjectDescPtr: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move address from BX to return value (return value after one parameter => SP+2*64bit)
  RET

TEXT ·Uint64ToTypeDescPtr(SB),0,$0-0 //Uint64ToTypeDescPtr: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move address from BX to return value (return value after one parameter => SP+2*64bit)
  RET


