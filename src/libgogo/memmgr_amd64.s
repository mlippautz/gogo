// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo Memory manager functions (ASM)
//

TEXT ·Brk(SB),$0-16 //Brk: 1 parameter, 1 return value
  MOVQ $12, AX //sys_brk (1 parameter)
  MOVQ 8(SP), DI //brk (first parameter => SP+8*64bit)
  SYSCALL //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS BRK_SUCCESS //Return result if successful
BRK_ERROR:
  NEGQ AX //Get errno
  MOVQ AX, 16(SP) //Return errno to indicate that an error occurred (return value after one parameters => SP+2*64bit)
  RET
BRK_SUCCESS:
  MOVQ $0, 16(SP) //Return 0 to indicate success (return value after one parameters => SP+2*64bit)
  RET

TEXT ·GetBrk(SB),$0-8 //GetBrk: no parameters, 1 return value
  MOVQ $12, AX //sys_brk (1 parameter)
  MOVQ $0, DI //brk (first parameter => SP+64bit)
  SYSCALL //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS GETBRK_SUCCESS //Return result if successful
GETBRK_ERROR:
  MOVQ $0, 8(SP) //Return 0 to indicate that an error occurred (return value after no parameters => SP+1*64bit)
  RET
GETBRK_SUCCESS:
  MOVQ AX, 8(SP) //First return value of syscall is in AX (return value after no parameters => SP+1*64bit)
  RET

TEXT ·TestMem(SB),$0-16 //Write: 1 parameter, 1 return value
  MOVQ 8(SP), AX //Move address to AX (first parameter => SP+64bit)
  MOVQ $1234567890, (AX) //Move some value to address
  MOVQ (AX), BX //Move value back to BX
  XORQ $1234567890, BX //value XOR value => 0
  MOVQ BX, 16(SP) //Return value (0 if successful)
  RET
