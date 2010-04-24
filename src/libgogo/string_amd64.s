// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo string functions (ASM)
//

TEXT ·StringLength(SB),$0-24 //StringLength: 1 parameter (2), 1 return value
  MOVQ $0, 24(SP) //Set return value to 0
  MOVW 16(SP), AX //String length is stored together with the string (first parameter = SP+64bit -> +64bit = SP+2*64bit)
  MOVW AX, 24(SP) //Move length to result with only 32 bits (return value after one parameter => SP+3*64bit)
  RET

TEXT ·StringLength2(SB),$0-16 //StringLength: 1 parameter, 1 return value
  MOVQ $0, 16(SP) //Set return value to 0
  MOVQ 8(SP), AX //Load string address to AX (first parameter = SP+64bit)
  MOVW 8(AX), BX //String length is stored together with the string
  MOVW BX, 16(SP) //Move length to result with only 32 bits (return value after one parameter => SP+3*64bit)
  RET

TEXT ·GetStringAddress(SB),$0-16 //ModifyString: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move string pointer address to AX (first parameter => SP+64bit)
  MOVQ (AX), BX //Load effective string address to BX
  MOVQ BX, 16(SP) //Move string address from BX to return value (return value after one parameter => SP+2*64bit)
  RET

TEXT ·GetStringFromAddress(SB),$0-16 //GetStringFromAddress: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move string address to AX (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move string address from AX to return value (return value after one parameter => SP+2*64bit)
  RET

TEXT ·SetStringAddressAndLength(SB),$0-24 //SetStringAddressAndLength: 3 parameters, no return value
  MOVQ 8(SP), AX //Move string address to AX (first parameter => SP+64bit)
  MOVQ 16(SP), BX //Move string address to BX (second parameter => SP+2*64bit)
  MOVQ 24(SP), CX //Move string length to CX (third parameter => SP+3*64bit)
  MOVQ BX, (AX) //Set string address
  MOVQ CX, 8(AX) //Set string length
  RET
