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
  MOVQ $4, AX //sys_write (3 parameters)
  MOVQ 8(SP), BX //fd (first parameter => SP+64bit)
  MOVQ 16(SP), CX //text (second parameter => SP+2*64bit)
  MOVQ 24(SP), DX //text length (third parameter => SP+3*64bit)
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  JLS WRITE_SUCCESS //Return result if successful
WRITE_ERROR:
  MOVQ $0, 40(SP) //Return 0 (return value after three parameters => SP+4*64bit+64bit?)
  //TODO: Error handling?
  RET
WRITE_SUCCESS:
  MOVQ AX, 40(SP) //First return value of syscall is in AX (return value after three parameters => SP+4*64bit+64bit?)
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
  MOVQ $0, 40(SP) //Return 0 (return value after three parameters => SP+4*64bit+64bit?)
  //TODO: Error handling?
  RET
READ_SUCCESS:
  MOVQ AX, 40(SP) //First return value of syscall is in AX (return value after three parameters => SP+4*64bit+64bit?)
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
  //MOVQ $0, 32(SP) //Return 0 (return value after three parameters => SP+3*64bit+64bit?)
  //TODO: Error handling?
  RET
FILEOPEN_SUCCESS:
  MOVQ AX, 32(SP) //First return value of syscall is in AX (return value after three parameters => SP+3*64bit+64bit?)
  RET

TEXT ·FileClose(SB),2,$0 //FileClose: 1 parameter, 1 return value
  MOVQ $6, AX //sys_close (3 parameters)
  MOVQ 8(SP), BX //filename (first parameter => SP+64bit)
  MOVQ $0, CX //flags (second parameter => SP+2*64bit)
  MOVQ $0, DX //not used (third parameter => SP+3*64bit)
  INT $0x80 //Linux syscall
  CMPQ AX, $0xFFFFFFFFFFFFF001 //Check for success
  //TODO: Error handling?
  RET

TEXT ·StringLength(SB),2,$0 //StringLength: 1 parameter, 1 return value
BBLEN_START:
  MOVQ 8(SP), AX //String (first parameter => SP+64bit)
  MOVQ $0, 24(SP) //Initialize length with 0
LOOP_BBLEN:
  CMPB (AX), $0 //Compare character with '\0'
  JE END_BBLEN //Terminate when '\0' has been found
  INCQ 24(SP) //Increase length (return value after one parameter => SP+2*64bit+64bit?)
  INCQ AX //Next character
  JMP LOOP_BBLEN //Continue
END_BBLEN:
  RET

TEXT ·InternalByteBufToString(SB),2,$0 //InternalByteBufToString: 2 parameters, no return value
  MOVQ 8(SP), AX //Buffer (first parameter => SP+64bit)
  MOVQ 16(SP), BX //String (second parameter => SP+2*64bit)
LOOP_BUFTOSTR:
  MOVB (AX), DX //Move byte...
  MOVB DX, (BX) //...to string
  CMPB (AX), $0 //Compare byte with 0
  JE END_BUFTOSTR //Terminate when 0 has been found
  INCQ AX //Next byte
  INCQ BX //Next character
  JMP LOOP_BUFTOSTR //Continue
END_BUFTOSTR:
  RET
