// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo Library functions (ASM)
//

TEXT ·CopyMem(SB),$0-24 //CopyMem: 3 parameters, no return value
  MOVQ 8(SP), AX //Move source address to AX (first parameter => SP+64bit)
  MOVQ 16(SP), BX //Move destination address to BX (second parameter => SP+2*64bit)
  MOVQ 24(SP), CX //Move length to CX (third parameter => SP+3*64bit)
  JCXZ COPYMEM_END //Return right away to end if length is 0
COPYMEM_LOOP:
  MOVB (AX), DX //Move source (one byte) to DX
  MOVB DX, (BX) //Move DX to destination (one byte)
  INCQ AX //Next source address
  INCQ BX //Next destination address
  LOOP COPYMEM_LOOP //Move to next address (and decrement length)
COPYMEM_END:
  RET

TEXT ·Exit(SB),$0-8 //Exit: 1 parameter, no return value
  MOVQ $60, AX //sys_exit (1 parameter)
  MOVQ 8(SP), DI //return code (first parameter => SP+1*64bit)
  SYSCALL //Linux syscall
  RET //Just to be sure (should never be reached)
