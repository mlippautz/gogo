// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

//
// Main parsing function. Corresponds to the EBNF main structure called 
// go_program.
//
func Parse( fd uint64 ) {
    var tok Token;
    tok.id = 0;
    tok.nextChar = 0;
    tok.nextToken = 0;    

    ParsePackageStatement(fd, &tok);    
    ParseImportStatementList(fd, &tok);
    ParseStructDeclList(fd, &tok);
    ParseVarDeclList(fd, &tok);
    ParseFuncDeclList(fd, &tok);

    // Scan the rest for debugging purposes
    // To be removed when parser is able to parse the complete EBNF
    // Is actually not EBNF compliant!
    GetNextToken(fd,&tok)
    for ; tok.id != TOKEN_EOS; {
        debugToken(&tok);
        GetNextToken(fd,&tok)
    }
}

//
// Parses: package identifier
// This is enforced by the go language as first statement in a source file.
//
func ParsePackageStatement(fd uint64, tok *Token) {    
    AssertNextToken(fd, tok, TOKEN_PACKAGE);
    AssertNextToken(fd, tok, TOKEN_IDENTIFIER);
    // package ok, value in tok.strValue
}

//
// Parses: { import_stmt }
// Function parsing the whole import block (which is optional) of a go program
//
func ParseImportStatementList(fd uint64, tok *Token) {
    var validImport uint64;
    for validImport = ParseImportStatement(fd, tok);
        validImport == 0;
        validImport = ParseImportStatement(fd, tok) { }
}


//
// Parses: "import" string
// This function parses a single import line.
// Returning 0 if import statement is valid, 1 otherwise.
//
func ParseImportStatement(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(fd,tok);
    if tok.id == TOKEN_IMPORT {
        AssertNextToken(fd, tok, TOKEN_STRING);
        // import ok, value in tok.strValue
        boolFlag = 0;
    } else {
        boolFlag = 1;
        SyncToken(tok);
    }    
    return boolFlag;
}

//
// Parses: { struct_decl }
//
func ParseStructDeclList(fd uint64, tok *Token) {
    var boolFlag uint64;
    for boolFlag = ParseStructDecl(fd, tok);
        boolFlag == 0;
        boolFlag = ParseStructDecl(fd, tok) { }
}

//
// Parses: "type" identifier "struct" "{" struct_var_decl_list "}" ";"
//
func ParseStructDecl(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(fd,tok);
    if tok.id == TOKEN_TYPE {
        AssertNextToken(fd, tok, TOKEN_IDENTIFIER);
        // identifier of struct in tok.strValue
        AssertNextToken(fd, tok, TOKEN_STRUCT);
        AssertNextToken(fd, tok, TOKEN_LCBRAC);
        ParseStructVarDeclList(fd, tok);
        AssertNextToken(fd, tok, TOKEN_RCBRAC);
        AssertNextToken(fd, tok, TOKEN_SEMICOLON);
        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }    
    return boolFlag;
}

//
// Parses: { struct_var_decl }
//
func ParseStructVarDeclList(fd uint64, tok *Token) {
    var boolFlag uint64;
    for boolFlag = ParseStructVarDecl(fd, tok);
        boolFlag == 0;
        boolFlag = ParseStructVarDecl(fd, tok) { }
}

//
// Parses: identifier type ";"
//
func ParseStructVarDecl(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;

    GetNextTokenSafe(fd,tok);
    if tok.id == TOKEN_IDENTIFIER {
        ParseType(fd, tok);
        AssertNextToken(fd, tok, TOKEN_SEMICOLON);
        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }    
    return boolFlag;    
}

//
// Parses: [ "[" integer "]" ] identifier 
//
func ParseType(fd uint64, tok *Token) {
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_LSBRAC {   
        AssertNextToken(fd, tok, TOKEN_INTEGER);
        // value of integer in tok.intValue
        AssertNextToken(fd, tok, TOKEN_RSBRAC);        
    } else {
        SyncToken(tok);
    }

    AssertNextToken(fd, tok, TOKEN_IDENTIFIER);
    // typename in tok.strValue
}

//
// Parses: [ "[" integer "]" ] identifier 
//
func ParseTypeOptional(fd uint64, tok *Token) {
    var es [255]uint64;

    GetNextTokenSafe(fd, tok);

    if tok.id == TOKEN_LSBRAC {        
        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_INTEGER  {
            es[0] = TOKEN_INTEGER;
            ParseError(tok.id,es,1);
        }       

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_RSBRAC  {
            es[0] = TOKEN_RSBRAC;
            ParseError(tok.id,es,1);
        }
    } else {
        tok.nextToken = tok.id;
    }

    GetNextTokenSafe(fd, tok);
    if tok.id != TOKEN_IDENTIFIER  {
        tok.nextToken = tok.id;
    }
}

//
//
//
func ParseVarDeclList(fd uint64, tok *Token) {
    var boolFlag uint64;
    for boolFlag = ParseVarDecl(fd, tok);
        boolFlag == 0;
        boolFlag = ParseVarDecl(fd, tok) { }
}

//
//
//
func ParseVarDecl(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(fd,tok);
    if tok.id == TOKEN_VAR {
        AssertNextToken(fd, tok, TOKEN_IDENTIFIER);
        // variable name in tok.strValue
        ParseType(fd, tok);

        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_ASSIGN {
            ParseExpression(fd, tok);        
        } else {
            SyncToken(tok);
        } 

        AssertNextToken(fd, tok, TOKEN_SEMICOLON);
        boolFlag = 0;
    } else {
        boolFlag = 1;
        SyncToken(tok);
    }    
    return boolFlag;
}

//
//
//
func ParseExpression(fd uint64, tok *Token) {
    ParseSimpleExpression(fd, tok);    
    ParseCmpOp(fd, tok);
}

//
//
//
func ParseCmpOp(fd uint64, tok *Token) {
    GetNextTokenSafe(fd, tok);
    if (tok.id == TOKEN_EQUALS) || (tok.id == TOKEN_NOTEQUAL) || 
        (tok.id == TOKEN_REL_LT) || (tok.id == TOKEN_REL_LTOE) || 
        (tok.id == TOKEN_REL_GT) || (tok.id == TOKEN_REL_GTOE) {
        ParseSimpleExpression(fd, tok);
    } else {
        SyncToken(tok);
    }
}

//
//
//
func ParseSimpleExpression(fd uint64, tok *Token) {
    var boolFlag uint64;
    ParseUnaryArithOp(fd, tok);
    ParseTerm(fd, tok);
    for boolFlag = ParseSimpleExpressionOp(fd, tok);
        boolFlag == 0;
        boolFlag = ParseSimpleExpressionOp(fd, tok) { }
}

//
//
//
func ParseSimpleExpressionOp(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    
    boolFlag = ParseUnaryArithOp(fd, tok);
    if boolFlag == 0 {
        // read +/-
    } else {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_REL_OR {
            // read ||
            boolFlag = 0;
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;   
        }
    }

    if boolFlag == 0 {
        ParseTerm(fd, tok);
    }

    return boolFlag;
}

//
//
//
func ParseUnaryArithOp(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_ARITH_PLUS {
        boolFlag = 0;
    } else {
        if tok.id == TOKEN_ARITH_MINUS {
            boolFlag = 0;            
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;
        }
    }

    return boolFlag;   
}

//
//
//
func ParseBinaryArithOp(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_ARITH_MUL {
        boolFlag = 0;
    } else {
        if tok.id == TOKEN_ARITH_DIV {
            boolFlag = 0;            
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;
        }
    }

    return boolFlag; 
}

//
//
//
func ParseTerm(fd uint64, tok *Token) {
    var boolFlag uint64;
    ParseFactor(fd, tok);
    for boolFlag = ParseTermOp(fd, tok);
        boolFlag == 0;
        boolFlag = ParseTermOp(fd, tok) { }      
}

//
//
//
func ParseTermOp(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    boolFlag = ParseBinaryArithOp(fd, tok);
    if boolFlag == 0 {
        // read *//
    } else {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_REL_AND {
            // read &&
            boolFlag = 0;
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;   
        }
    }

    if boolFlag == 0 {
        ParseFactor(fd, tok);
    }

    return boolFlag;
}

//
//
//
func ParseFactor(fd uint64, tok *Token) uint64 {
    var es [255]uint64;
    var doneFlag uint64 = 1;
    var boolFlag uint64;

    GetNextTokenSafe(fd, tok);

    if (doneFlag == 1) && (tok.id == TOKEN_OP_ADR) {
        GetNextTokenSafe(fd, tok); 

        if tok.id == TOKEN_IDENTIFIER {
            ParseSelector(fd ,tok);
            ParseFunctionCallOptional(fd, tok);
            doneFlag = 0;
        } else {
            es[0] = TOKEN_IDENTIFIER;
            ParseError(tok.id, es, 1);
        }
    }
    if (doneFlag == 1) && (tok.id == TOKEN_IDENTIFIER) {
        ParseSelector(fd ,tok);
        ParseFunctionCallOptional(fd, tok);
        doneFlag = 0;
    } 
    if (doneFlag == 1) && (tok.id == TOKEN_INTEGER) {
        doneFlag = 0;
    }
    if (doneFlag) == 1 && (tok.id == TOKEN_STRING) {
        doneFlag = 0;
    }
    if (doneFlag) == 1 && (tok.id == TOKEN_LBRAC) {

        ParseExpression(fd, tok);
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_RBRAC {
            doneFlag = 0;
        } else {
            es[0] = TOKEN_RBRAC;
            ParseError(tok.id,es,1);
        }
    }
    if (doneFlag == 1) && (tok.id == TOKEN_NOT) {
        ParseFactor(fd, tok);
        doneFlag = 0;
    }

    if doneFlag != 0 {
        boolFlag = 1;
        tok.nextToken = tok.id;
    } else {
        boolFlag = 0;
    }
    return boolFlag;
}

//
//
//
func ParseSelector(fd uint64, tok *Token) {
    var boolFlag uint64;
    for boolFlag = ParseSelectorSub(fd, tok);boolFlag == 0; boolFlag = ParseSelectorSub(fd, tok) {
    }
}

//
//
//
func ParseSelectorSub(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;

    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_PT {
        AssertNextToken(fd, tok, TOKEN_IDENTIFIER);
        // value in tok.strValue
        boolFlag = 0;
    } else {
        if tok.id == TOKEN_LSBRAC {
            GetNextTokenSafe(fd, tok);
            if tok.id == TOKEN_INTEGER {
                
            } else {
                if tok.id == TOKEN_IDENTIFIER {  
                    ParseSelector(fd, tok);
                }
            } 

            AssertNextToken(fd, tok, TOKEN_RSBRAC);
            boolFlag = 0;
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;
        }
    }
    return boolFlag;
}

func ParseFuncDeclList(fd uint64, tok *Token) {
    var boolFlag uint64; 
    for boolFlag = ParseFuncDeclListSub(fd, tok);
        boolFlag == 0; 
        boolFlag = ParseFuncDeclListSub(fd, tok) { }
}

func ParseFuncDeclListSub(fd uint64, tok *Token) uint64 {
    var es [255]uint64;
    var boolFlag uint64;    
    boolFlag = ParseFuncDeclHead(fd, tok);
    if boolFlag == 0 {

        boolFlag = ParseFuncDeclRaw(fd, tok);
        if boolFlag != 0 {
            boolFlag = ParseFuncDecl(fd, tok);
        }
        if boolFlag != 0 {
            es[0] = TOKEN_SEMICOLON;
            es[1] = TOKEN_LCBRAC;
            ParseError(tok.id,es,2);
        }
    }
    
    return boolFlag;
}

func ParseFuncDeclHead(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_FUNC {
        AssertNextToken(fd, tok, TOKEN_IDENTIFIER);
        // function name in tok.strValue
        AssertNextToken(fd, tok, TOKEN_LBRAC);
        ParseIdentifierTypeList(fd, tok);
        AssertNextToken(fd, tok, TOKEN_RBRAC);
        ParseTypeOptional(fd, tok);
        boolFlag = 0;
    } else {    
        SyncToken(tok);
        boolFlag = 1;
    }
    return boolFlag;
}

func ParseFuncDeclRaw(fd uint64, tok *Token) uint64 {
    var boolFlag uint64 = 1;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_SEMICOLON {
        boolFlag = 0;
    } else {
        SyncToken(tok);
    }
    return boolFlag;
}

func ParseFuncDecl(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;

    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_LCBRAC {
        ParseVarDeclList(fd, tok);
        ParseStatementSequence(fd, tok);
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_RETURN {
            ParseExpression(fd, tok);
            AssertNextToken(fd, tok, TOKEN_SEMICOLON);
        } else {
            SyncToken(tok);
        }
        AssertNextToken(fd, tok, TOKEN_RCBRAC);
        boolFlag = 0;
    } else {
        SyncToken(tok);
        boolFlag = 1;
    }
    return boolFlag;
}

func ParseIdentifierTypeList(fd uint64, tok *Token) {
    var boolFlag uint64;

    boolFlag = ParseIdentifierType(fd, tok);
    if boolFlag == 0 {
        for boolFlag = ParseIdentifierTypeListSub(fd, tok);boolFlag == 0; boolFlag = ParseIdentifierTypeListSub(fd, tok) {
        }   
    }
}

func ParseIdentifierTypeListSub(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_COLON {
        boolFlag = ParseIdentifierType(fd, tok);
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }
    return boolFlag;
}

func ParseIdentifierType(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;

    GetNextTokenSafe(fd, tok);
    if tok.id != TOKEN_IDENTIFIER {
        tok.nextToken = tok.id;
        boolFlag = 1;        
    } else {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_ARITH_MUL {

        } else {
            tok.nextToken = tok.id;
        }
        ParseType(fd, tok);
        boolFlag = 0;
    }

    return boolFlag;
}

func ParseStatementSequence(fd uint64, tok *Token) {
    var boolFlag uint64; 
    for boolFlag = ParseStatement(fd, tok);boolFlag == 0; boolFlag = ParseStatement(fd, tok) {
    }
}

func ParseStatement(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    var doneFlag uint64;
    var es [255]uint64; 

    doneFlag = 1;

    GetNextTokenSafe(fd, tok);
    if (doneFlag == 1) && (tok.id == TOKEN_IDENTIFIER) {
        // Could be assignment or a function call.
        // Cannnot be resolved until selectors are all parsed
        // To be improved!
        ParseSelector(fd, tok);
        tok.nextToken = tok.id;
        boolFlag = ParseAssignment(fd, tok);
        if boolFlag != 0 {
            tok.nextToken = tok.id;
            boolFlag = ParseFunctionCallStatement(fd, tok);
        }        
        if boolFlag != 0 {
            ParseError(tok.id,es,0);
        }
        doneFlag = 0;
    }
    
    if (doneFlag == 1) && (tok.id == TOKEN_IF) {
        tok.nextToken = tok.id;
        ParseIfStatement(fd, tok);
        doneFlag = 0;
    }

    if (doneFlag == 1) && (tok.id == TOKEN_FOR) {
        tok.nextToken = tok.id;
        ParseForStatement(fd, tok);
        doneFlag = 0;
    }

    if (doneFlag == 1) && (tok.id == TOKEN_SEMICOLON) {
        // NOP    
        doneFlag = 0;
    }

    if doneFlag != 0 {
        tok.nextToken = tok.id;
        boolFlag = 1;
    } else {
        boolFlag = 0;
    }

    return boolFlag;
}

func ParseAssignment(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    var es [255]uint64;

    //GetNextTokenSafe(fd, tok);
    //if tok.id == TOKEN_IDENTIFIER {
    //    ParseSelector(fd, tok);
        
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_ASSIGN {
            
            ParseExpression(fd, tok);

            GetNextTokenSafe(fd, tok);
            if tok.id != TOKEN_SEMICOLON {
                es[0] = TOKEN_SEMICOLON;
                ParseError(tok.id, es, 1);
            }   
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;
        }
//    } 
//else {
//        tok.nextToken = tok.id;
 //       boolFlag = 1;
  //  }

    return boolFlag;
}

func ParseAssignmentWithoutSC(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
        
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_ASSIGN {      
        ParseExpression(fd, tok);
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }

    return boolFlag;
}

func ParseFunctionCallOptional(fd uint64, tok *Token) {
    var es [255]uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_LBRAC {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_RBRAC {

        } else {
            tok.nextToken = tok.id;
            ParseExpressionList(fd, tok);
            GetNextTokenSafe(fd, tok);
            if tok.id != TOKEN_RBRAC {
                es[0] = TOKEN_RBRAC;
                ParseError(tok.id, es, 1);
            }        
        }
    } else {
        SyncToken(tok);
    }
}

func ParseFunctionCall(fd uint64, tok *Token) {
    var es [255]uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_LBRAC {
        if tok.id == TOKEN_RBRAC {

        } else {
            ParseExpressionList(fd, tok);
            GetNextTokenSafe(fd, tok);
            if tok.id != TOKEN_RBRAC {
                es[0] = TOKEN_RBRAC;
                ParseError(tok.id, es, 1);
            }        
        }     
    } else {
        es[0] = TOKEN_LBRAC;
        ParseError(tok.id,es,1);
    }
}

func ParseExpressionList(fd uint64, tok *Token) {
    var boolFlag uint64;

    ParseExpression(fd, tok);
    for boolFlag = ParseExpressionListSub(fd, tok);boolFlag == 0; boolFlag = ParseExpressionListSub(fd, tok) {
    }   
}

func ParseExpressionListSub(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_COLON {
        ParseExpression(fd, tok);
        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }
    return boolFlag;   
}

func ParseFunctionCallStatement(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_LBRAC {
        tok.nextToken = tok.id;
        ParseFunctionCall(fd, tok);                
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    return boolFlag;
}

func ParseForStatement(fd uint64, tok *Token) {
    var es [255]uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_FOR {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_SEMICOLON {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            GetNextTokenSafe(fd,tok);
            if tok.id == TOKEN_IDENTIFIER {
                ParseSelector(fd, tok);
            } else {
                es[0] = TOKEN_IDENTIFIER;
                ParseError(tok.id,es,1);
            }
            ParseAssignmentWithoutSC(fd, tok);
        }
        
        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_SEMICOLON {
            es[0] = TOKEN_SEMICOLON;
            ParseError(tok.id,es,1);
        }

        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_SEMICOLON {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            ParseExpression(fd, tok);
        }

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_SEMICOLON {
            es[0] = TOKEN_SEMICOLON;
            ParseError(tok.id,es,1);
        }

        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_LCBRAC {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            GetNextTokenSafe(fd,tok);
            if tok.id == TOKEN_IDENTIFIER {
                ParseSelector(fd, tok);
            } else {
                es[0] = TOKEN_IDENTIFIER;
                ParseError(tok.id,es,1);
            }
            ParseAssignmentWithoutSC(fd, tok);
        }

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_LCBRAC {
            es[0] = TOKEN_LCBRAC;
            ParseError(tok.id,es,1);
        }
        
        ParseStatementSequence(fd, tok);

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_RCBRAC {
            es[0] = TOKEN_RCBRAC;
            ParseError(tok.id,es,1);
        }        

    } else {
        tok.nextToken = tok.id;
    }   
}

func ParseIfStatement(fd uint64, tok *Token) {
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_IF {
        ParseExpression(fd, tok);
        AssertNextToken(fd, tok, TOKEN_LCBRAC);
        ParseStatementSequence(fd, tok);
        AssertNextToken(fd, tok, TOKEN_RCBRAC);

        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_ELSE {
            ParseElseStatement(fd, tok);
        } else {
            SyncToken(tok);
        }

    } else {
        SyncToken(tok);
    }
}

func ParseElseStatement(fd uint64, tok *Token) {
    AssertNextToken(fd, tok, TOKEN_LCBRAC);
    ParseStatementSequence(fd, tok);
    AssertNextToken(fd, tok, TOKEN_RCBRAC);
}
