//
// GoGo's very basic Library functions
//

TEXT ·Exit(SB),1,$0 //Exit: 1 parameter, no return value
  MOVQ $1, AX //sys_exit (1 parameter)
  MOVQ 8(SP), BX //return code (first parameter => SP+1*64bit)
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  //If sys_exit was successful this code is not called
  //TODO: Error handling?
  RET

TEXT ·Write(SB),4,$0 //Write: 3 parameters, 1 return value
  MOVQ $4, AX //sys_write (4 parameters)
  MOVQ 8(SP), BX //fd (first parameter => SP+64bit)
	MOVQ 16(SP), CX //text (second parameter => SP+2*64bit)
	MOVQ 24(SP), DX //text length (third parameter => SP+3*64bit)
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS WRITE_SUCCESS //Return result if successful
WRITE_ERROR:
  MOVQ $0, 40(SP) //Return 0 (return value after three parameters => SP+4*64bit+64bit?)
  //TODO: Error handling?
WRITE_SUCCESS:
  MOVQ AX, 40(SP) //First return value of syscall is in AX (return value after three parameters => SP+4*64bit+64bit?)
  RET

TEXT ·Read(SB),4,$0 //Write: 3 parameters, 1 return value
  MOVQ $3, AX //sys_read (4 parameters)
  MOVQ 8(SP), BX //fd (first parameter => SP+64bit)
	MOVQ 16(SP), CX //buffer (second parameter => SP+2*64bit)
	MOVQ 24(SP), DX //buffer size (third parameter => SP+3*64bit)
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS READ_SUCCESS //Return result if successful
READ_ERROR:
  MOVQ $0, 40(SP) //Return 0 (return value after three parameters => SP+4*64bit+64bit?)
  //TODO: Error handling?
READ_SUCCESS:
  MOVQ AX, 40(SP) //First return value of syscall is in AX (return value after three parameters => SP+4*64bit+64bit?)
  RET
