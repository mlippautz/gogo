//
// GoGo Library functions (ASM)
//

TEXT ·Exit(SB),1,$0 //Exit: 1 parameter, no return value
  MOVQ $1, AX //sys_exit (1 parameter)
  MOVQ 8(SP), BX //return code (first parameter => SP+1*64bit)
  INT $0x80 //Linux syscall
  RET //Just to be sure (should never be reached)

TEXT ·StringLength(SB),3,$0 //StringLength: 1 parameter, 1 return value
  MOVQ $0, 24(SP) //Set return value to 0
  MOVW 16(SP), AX //String length is stored together with the string (first parameter = SP+64bit -> +64bit = SP+2*64bit)
  MOVW AX, 24(SP) //Move length to result with only 32 bits (return value after one parameter => SP+3*64bit)
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

TEXT ·Write(SB),5,$0 //Write: 3 parameters, 1 return value
  MOVQ $4, AX //sys_write (3 parameters)
  MOVQ 8(SP), BX //fd (first parameter => SP+64bit)
  MOVQ 16(SP), CX //text (second parameter => SP+2*64bit)
  MOVQ 32(SP), DX //text length (third parameter => SP+4*64bit)
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS WRITE_SUCCESS //Return result if successful
WRITE_ERROR:
  MOVQ $0, 40(SP) //Return 0 to indicate that an error occured (return value after three parameters => SP+5*64bit)
  RET
WRITE_SUCCESS:
  MOVQ AX, 40(SP) //First return value of syscall is in AX (return value after three parameters => SP+5*64bit)
  RET

//Prototype (does not work yet fully as expected)
/*TEXT ·PrintChar(SB),1,$0 //PrintChar: 1 parameter, no return value
  MOVQ $4, AX //sys_write (3 parameters)
  MOVQ $1, BX //fd (1 = stdout)
  LEAQ 8(SP), CX //text (address of second parameter => SP+64bit)
  MOVQ $1, DX //text length (1)
  INT $0x80 //Linux syscall
  RET*/

//Prototype (does not work yet fully as expected)
/*TEXT ·GetChar(SB),2,$0 //Read: 1 parameter, 1 return value
  MOVQ $3, AX //sys_read (3 parameters)
  MOVQ 8(SP), BX //fd (first parameter => SP+64bit)
  LEAQ 8(SP), CX //buffer (reuse first parameter => SP+64bit)
  MOVQ $1, DX //buffer size (size 1)
  INT $0x80 //Linux syscall
  MOVQ $0, DX //Overwrite DX (initialize with 0)
  MOVB (CX), DL //Move buffer to DL
  MOVQ DX, 16(SP) //Move whole DX register to result (return value after one parameter => SP+2*64bit)
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS READ_SUCCESS //Return result if successful
GETCHAR_ERROR:
  MOVQ $0, 16(SP) //Return 0 (return value after one parameter => SP+2*64bit)
GETCHAR_SUCCESS:
  RET*/

TEXT ·FileClose(SB),2,$0 //FileClose: 1 parameter, 1 return value
  MOVQ $6, AX //sys_close (3 parameters)
  MOVQ 8(SP), BX //filename (first parameter => SP+64bit)
  MOVQ $0, CX //not used
  MOVQ $0, DX //not used
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS CLOSE_SUCCESS //Return result if successful
CLOSE_ERROR:
  NEGQ AX //Negate AX to get errno
  MOVQ AX, 16(SP) //Return errno
  RET
CLOSE_SUCCESS:
  MOVQ $0, 16(SP) //Set to zero to indicate successful call (return value after one parameter => SP+2*64bit)
  RET

//--- Cleanup necessary from here onwards (most functions don't work properly!)

TEXT ·Read(SB),5,$0 //Read: 3 parameters, 1 return value
  MOVQ $3, AX //sys_read (3 parameters)
  MOVQ 8(SP), BX //fd (first parameter => SP+64bit)
  MOVQ 16(SP), CX //buffer (second parameter => SP+2*64bit)
  MOVQ 32(SP), DX //buffer size (third parameter => SP+4*64bit)
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

TEXT ·FileOpen(SB),4,$0 //FileOpen: 2 parameters, 1 return value
  MOVQ $5, AX //sys_open (2 parameters)
  MOVQ 8(SP), BX //filename (first parameter => SP+64bit)
  MOVQ 24(SP), CX //flags (second parameter => SP+3*64bit)
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
