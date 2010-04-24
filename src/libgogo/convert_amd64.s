// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo conversion functions (ASM)
//

TEXT ·ToIntFromByte(SB),$0-16 //ToIntFromByte: 1 parameter, 1 return value
  MOVQ $0, AX //Set AX to 0
  MOVB 8(SP), AL //Move byte parameter to AL (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move whole AX register with byte parameter to result (return value after one parameter => SP+2*64bit)
  RET

TEXT ·ToByteFromInt(SB),$0-16 //ToByteFromInt: 1 parameter, 1 return value
  MOVQ $0, 16(SP) //Clear whole return value (return value after one parameter => SP+2*64bit)
  MOVQ 8(SP), AX //Move whole parameter to AX (first parameter => SP+64bit)
  MOVB AL, 16(SP) //Move AL (last byte of parameter) to result (return value after one parameter => SP+2*64bit)
  RET

TEXT ·ToUint64FromBytePtr(SB),$0-16 //ToUint64FromBytePtr: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move address from BX to return value (return value after one parameter => SP+2*64bit)
  RET

TEXT ·ToUint64FromUint64Ptr(SB),$0-16 //ToUint64FromUint64Ptr: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move address from BX to return value (return value after one parameter => SP+2*64bit)
  RET

TEXT ·GetStringFromAddress(SB),$0-16 //GetStringFromAddress: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move string address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move string address from AX to return value (return value after one parameter => SP+2*64bit)
  RET
