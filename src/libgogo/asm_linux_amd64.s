//
// GoGo Library functions (ASM)
//

TEXT ·Exit(SB),1,$0 //Exit: 1 parameter, no return value
  MOVQ $1, AX //sys_exit (1 parameter)
  MOVQ 8(SP), BX //return code (first parameter => SP+1*64bit)
  INT $0x80 //Linux syscall
  RET //Just to be sure (should never be reached)

TEXT ·StringLength(SB),2,$0 //StringLength: 1 parameter, 1 return value
STRLEN_START:
  MOVQ 8(SP), AX //String (first parameter => SP+64bit)
  MOVQ $0, 24(SP) //Initialize length with 0
LOOP_STRLEN:
  CMPB (AX), $0 //Compare character with '\0'
  JE END_STRLEN //Terminate when '\0' has been found
  INCQ 24(SP) //Increase length (return value after one parameter => SP+3*64bit)
  INCQ AX //Next character
  JMP LOOP_STRLEN //Continue
END_STRLEN:
  RET

TEXT ·ToIntFromByte(SB),2,$0 //ToIntFromByte: 1 parameter, 1 return value
  MOVQ $0, AX //Set AX to 0
  MOVB 8(SP), AL //Move byte parameter to AL (first parameter => SP+64bit)
  MOVQ AX, 16(SP) //Move whole AX register with byte parameter to result (return value after one parameter => SP+2*64bit)
  RET

TEXT ·ToByteFromInt(SB),2,$0 //ToByteFromInt: 1 parameter, 1 return value
  MOVQ $0, 16(SP) //Clear whole return value (return value after one parameter => SP+2*64bit)
  MOVQ 8(SP), AX //Move whole parameter to AX (first parameter => SP+64bit)
  MOVB AL, 16(SP) //Move AL (last byte of parameter) to result (return value after one parameter => SP+2*64bit)
  RET

//--- Cleanup necessary from here onwards (most functions don't work properly!)

TEXT ·Write(SB),4,$0 //Write: 3 parameters, 1 return value
  MOVQ $4, AX //sys_write (3 parameters)
  MOVQ 8(SP), BX //fd (first parameter => SP+64bit)
  MOVQ 16(SP), CX //text (second parameter => SP+2*64bit)
  MOVQ 24(SP), DX //text length (third parameter => SP+3*64bit)
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS WRITE_SUCCESS //Return result if successful
WRITE_ERROR:
  MOVQ $0, 40(SP) //Return 0 (return value after three parameters => SP+5*64bit)
  //TODO: Error handling?
  RET
WRITE_SUCCESS:
  MOVQ AX, 40(SP) //First return value of syscall is in AX (return value after three parameters => SP+5*64bit)
  RET

TEXT ·Read(SB),4,$0 //Read: 3 parameters, 1 return value
  MOVQ $3, AX //sys_read (3 parameters)
  MOVQ 8(SP), BX //fd (first parameter => SP+64bit)
  MOVQ 16(SP), CX //buffer (second parameter => SP+2*64bit)
  MOVQ 24(SP), DX //buffer size (third parameter => SP+3*64bit)
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS READ_SUCCESS //Return result if successful
READ_ERROR:
  MOVQ $0, 40(SP) //Return 0 (return value after three parameters => SP+5*64bit)
  //TODO: Error handling?
  RET
READ_SUCCESS:
  MOVQ AX, 40(SP) //First return value of syscall is in AX (return value after three parameters => SP+5*64bit)
  RET

TEXT ·FileOpen(SB),3,$0 //FileOpen: 2 parameters, 1 return value
  MOVQ $5, AX //sys_open (2 parameters)
  MOVQ 8(SP), BX //filename (first parameter => SP+64bit)
  MOVQ 16(SP), CX //flags (second parameter => SP+2*64bit)
  MOVQ $0, DX //not used
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS FILEOPEN_SUCCESS //Return result if successful
FILEOPEN_ERROR:
  //MOVQ $0, 32(SP) //Return 0 (return value after three parameters => SP+4*64bit)
  //TODO: Error handling?
  RET
FILEOPEN_SUCCESS:
  MOVQ AX, 32(SP) //First return value of syscall is in AX (return value after three parameters => SP+4*64bit)
  RET

TEXT ·FileClose(SB),2,$0 //FileClose: 1 parameter, 1 return value
  MOVQ $6, AX //sys_close (3 parameters)
  MOVQ 8(SP), BX //filename (first parameter => SP+64bit)
  MOVQ $0, CX //flags (second parameter => SP+2*64bit)
  MOVQ $0, DX //not used
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  //TODO: Error handling?
  RET
