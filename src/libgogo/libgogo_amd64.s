// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo Library functions (ASM)
//

TEXT ·StringLength(SB),$0-24 //StringLength: 1 parameter (2), 1 return value
  MOVQ $0, 24(SP) //Set return value to 0
  MOVW 16(SP), AX //String length is stored together with the string (first parameter = SP+64bit -> +64bit = SP+2*64bit)
  MOVW AX, 24(SP) //Move length to result with only 32 bits (return value after one parameter => SP+3*64bit)
  RET

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

TEXT ·Exit(SB),$0-8 //Exit: 1 parameter, no return value
  MOVQ $60, AX //sys_exit (1 parameter)
  MOVQ 8(SP), DI //return code (first parameter => SP+1*64bit)
  SYSCALL //Linux syscall
  RET //Just to be sure (should never be reached)

TEXT ·Write(SB),$0-40 //Write: 3 parameters (4), 1 return value
  MOVQ $1, AX //sys_write (3 parameters)
  MOVQ 8(SP), DI //fd (first parameter => SP+64bit)
  MOVQ 16(SP), SI //text (second parameter => SP+2*64bit)
  MOVQ 32(SP), DX //text length (third parameter => SP+4*64bit)
  SYSCALL //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS WRITE_SUCCESS //Return result if successful
WRITE_ERROR:
  MOVQ $0, 40(SP) //Return 0 to indicate that an error occurred (return value after three parameters => SP+5*64bit)
  RET
WRITE_SUCCESS:
  MOVQ AX, 40(SP) //First return value of syscall is in AX (return value after three parameters => SP+5*64bit)
  RET

TEXT ·PrintChar(SB),$0-8 //PrintChar: 1 parameter, no return value
  MOVQ $1, AX //sys_write (3 parameters)
  MOVQ $1, DI //fd (1 = stdout)
  LEAQ 8(SP), SI //text (address of second parameter => SP+64bit)
  MOVQ $1, DX //text length (1)
  SYSCALL //Linux syscall
  RET

TEXT ·Read(SB),$0-40 //Read: 3 parameters (4), 1 return value
  MOVQ $0, AX //sys_read (3 parameters)
  MOVQ 8(SP), DI //fd (first parameter => SP+64bit)
  MOVQ 16(SP), SI //buffer (second parameter => SP+2*64bit)
  MOVQ 32(SP), DX //buffer size (third parameter => SP+4*64bit)
  SYSCALL //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS READ_SUCCESS //Return result if successful
READ_ERROR:
  MOVQ $0, 40(SP) //Return 0 to indicate that an error occurred (return value after three parameters => SP+5*64bit)
  RET
READ_SUCCESS:
  MOVQ AX, 40(SP) //First return value of syscall is in AX (return value after three parameters => SP+5*64bit)
  RET

TEXT ·GetChar(SB),$0-16 //Read: 1 parameter (2), 1 return value
  MOVQ $0, AX //sys_read (3 parameters)
  MOVQ 8(SP), DI //fd (first parameter => SP+64bit)
  MOVQ $0, 16(SP) //Initialize result with 0
  LEAQ 16(SP), SI //buffer (return value after one parameter => SP+2*64bit)
  MOVQ $1, DX //buffer size (size 1)
  SYSCALL //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS READ_SUCCESS //Return result if successful
GETCHAR_ERROR:
  MOVQ $0, 16(SP) //Return 0 (return value after one parameter => SP+2*64bit)
GETCHAR_SUCCESS:
  RET

TEXT ·FileOpen(SB),$0-32 //FileOpen: 2 parameters (3), 1 return value
  MOVQ $2, AX //sys_open (2 parameters)
  MOVQ 8(SP), DI //filename (first parameter => SP+64bit)
  MOVQ 24(SP), SI //flags (second parameter => SP+3*64bit)
  MOVQ $0, DX //not used
  SYSCALL //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS FILEOPEN_SUCCESS //Return result if successful
FILEOPEN_ERROR:
  MOVQ $0, 32(SP) //Return 0 to indicate that an error occured (return value after three parameters => SP+4*64bit)
  RET
FILEOPEN_SUCCESS:
  MOVQ AX, 32(SP) //First return value of syscall is in AX (return value after three parameters => SP+4*64bit)
  RET

TEXT ·FileClose(SB),$0-16 //FileClose: 1 parameter, 1 return value
  MOVQ $3, AX //sys_close (3 parameters)
  MOVQ 8(SP), DI //filename (first parameter => SP+64bit)
  MOVQ $0, SI //not used
  MOVQ $0, DX //not used
  SYSCALL //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS CLOSE_SUCCESS //Return result if successful
CLOSE_ERROR:
  NEGQ AX //Negate AX to get errno
  MOVQ AX, 16(SP) //Return errno
  RET
CLOSE_SUCCESS:
  MOVQ $0, 16(SP) //Set to zero to indicate successful call (return value after one parameter => SP+2*64bit)
  RET
