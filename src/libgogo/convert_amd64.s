// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo conversion functions (ASM)
//

TEXT ·ToIntFromByte(SB),$0-16 //ToIntFromByte: 1 parameter, 1 return value
  MOVQ $0, 16(SP) //Initialize return value (return value after one parameter => SP+2*64bit)
  MOVB 8(SP), AX //Move byte parameter to AX (first parameter => SP+64bit)
  MOVB AX, 16(SP) //Move whole AX register with byte parameter to result (return value after one parameter => SP+2*64bit)
  RET

TEXT ·ToByteFromInt(SB),$0-16 //ToByteFromInt: 1 parameter, 1 return value
  MOVQ $0, 16(SP) //Clear whole return value (return value after one parameter => SP+2*64bit)
  MOVB 8(SP), AX //Move byte parameter (part) to AX (first parameter => SP+64bit)
  MOVB AX, 16(SP) //Move AL (last byte of parameter) to result (return value after one parameter => SP+2*64bit)
  RET

TEXT ·ToUint64FromBytePtr(SB),$0-16 //ToUint64FromBytePtr: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move address from BX to return value (return value after one parameter => SP+2*64bit)
  RET

TEXT ·ToUint64FromUint64Ptr(SB),$0-16 //ToUint64FromUint64Ptr: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move address from AX to return value (return value after one parameter => SP+2*64bit)
  RET
  
TEXT ·ToUint64PtrFromUint64(SB),$0-16 //ToUint64PtrFromUint64: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move address from AX to return value (return value after one parameter => SP+2*64bit)
  RET
  
TEXT ·ToUint64FromStringPtr(SB),$0-16 //ToUint64FromStringPtr: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move address from AX to return value (return value after one parameter => SP+2*64bit)
  RET

TEXT ·GetStringFromAddress(SB),$0-16 //GetStringFromAddress: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move string address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move string address from AX to return value (return value after one parameter => SP+2*64bit)
  RET
