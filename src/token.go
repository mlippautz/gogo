// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

//
// Set of recognized tokens 
//
var TOKEN_IDENTIFIER uint64 = 1;  // Identifier
var TOKEN_STRING uint64 = 2;      // String using "..."
var TOKEN_EOS uint64 = 3;         // End of Scan
var TOKEN_LBRAC uint64 = 4;       // Left bracket '('
var TOKEN_RBRAC uint64 = 5;       // Right bracket ')'
var TOKEN_LSBRAC uint64 = 6;      // Left square bracket '['
var TOKEN_RSBRAC uint64 = 7;      // Right square bracket ']'
var TOKEN_INTEGER uint64 = 8;     // Integer number
var TOKEN_LCBRAC uint64 = 9;      // Left curly bracket '{'
var TOKEN_RCBRAC uint64 = 10;     // Right curly bracket '}'
var TOKEN_PT uint64 = 11;         // Point '.'
var TOKEN_NOT uint64 = 12;        // Boolean negation '!'
var TOKEN_NOTEQUAL uint64 = 13;   // Comparison, not equal '!='
var TOKEN_SEMICOLON uint64 = 14;  // Semi-colon ';'
var TOKEN_COLON uint64 = 15;      // Colon ','
var TOKEN_ASSIGN uint64 = 16;     // Assignment '='
var TOKEN_EQUALS uint64 = 17;     // Equal comparison '=='
var TOKEN_CHAR uint64 = 18;       // Single Quoted Character 'x'
var TOKEN_REL_AND uint64 = 19;    // AND Relation '&&'
var TOKEN_REL_OR uint64 = 20;     // OR Relation '||'
var TOKEN_REL_GTOE uint64 = 21;   // Greather-Than or Equal '>='
var TOKEN_REL_GT uint64 = 22;     // Greather-Than '>'
var TOKEN_REL_LTOE uint64 = 23;   // Less-Than or Equal '<='
var TOKEN_REL_LT uint64 = 24;     // Less-Than '<'
var TOKEN_ARITH_PLUS uint64 = 25; // Arith. Plus '+'
var TOKEN_ARITH_MINUS uint64 = 26;// Arith. Minus '-'
var TOKEN_ARITH_MUL uint64 = 27;  // Arith. Multiplication '*'
var TOKEN_ARITH_DIV uint64 = 28;  // Arith. Division '/'
var TOKEN_OP_ADR uint64 = 29;     // Address operator '&'

//
// Advanced tokens, that are generated in the 2nd step from identifiers
// The tokens represent the corresponding language keywords.
//
var TOKEN_FOR uint64 = 101;
var TOKEN_IF uint64 = 102;
var TOKEN_TYPE uint64 = 103;
var TOKEN_CONST uint64 = 104;
var TOKEN_VAR uint64 = 105;
var TOKEN_STRUCT uint64 = 106;
var TOKEN_RETURN uint64 = 107;
var TOKEN_FUNC uint64 = 108;
var TOKEN_PACKAGE uint64 = 109;
var TOKEN_IMPORT uint64 = 110;

//
// Helper functions
//

func TokenToString (id uint64) string {
    var retStr string;

    if id == TOKEN_IDENTIFIER {
        retStr = "<identifier>";
    }
    if id == TOKEN_STRING {
        retStr = "<string>";
    }
    if id == TOKEN_EOS {
        retStr = "<END-OF-SCAN>";
    }
    if id == TOKEN_LBRAC {
        retStr = "(";
    }
    if id == TOKEN_RBRAC {
        retStr = ")";
    }
    if id == TOKEN_LSBRAC {
        retStr = "[";
    }
    if id == TOKEN_RSBRAC {
        retStr = "]";
    }
    if id == TOKEN_INTEGER {
        retStr = "<integer>";
    }
    if id == TOKEN_LCBRAC {
        retStr = "{";
    }
    if id == TOKEN_RCBRAC {
        retStr = "}";
    }
    if id == TOKEN_PT {
        retStr = ".";
    }
    if id == TOKEN_NOT {
        retStr = "!";
    }
    if id == TOKEN_NOTEQUAL {
        retStr = "!=";
    }
    if id == TOKEN_SEMICOLON {
        retStr = ";";
    }
    if id == TOKEN_COLON {
        retStr = ",";
    }
    if id == TOKEN_ASSIGN {
        retStr = "=";
    }
    if id == TOKEN_EQUALS {
        retStr = "==";
    }
    if id == TOKEN_CHAR {
        retStr = "'<char>'";
    }
    if id == TOKEN_REL_AND {
        retStr = "&&";
    }
    if id == TOKEN_REL_OR {
        retStr = "||"
    }
    if id == TOKEN_REL_GTOE {
        retStr = ">=";
    }
    if id == TOKEN_REL_GT {
        retStr = ">";
    }
    if id == TOKEN_REL_LTOE {
        retStr = "<=";
    }
    if id == TOKEN_REL_LT {
        retStr = "<";
    }
    if id == TOKEN_ARITH_PLUS {
        retStr = "+";
    }
    if id == TOKEN_ARITH_MINUS {
        retStr = "-";
    }
    if id == TOKEN_ARITH_MUL {
        retStr = "*";
    }
    if id == TOKEN_ARITH_DIV {
        retStr = "/";
    }
    if id == TOKEN_OP_ADR {
        retStr = "&";
    }
    if id == TOKEN_FOR {
        retStr = "for";
    }
    if id == TOKEN_IF {
        retStr = "if";
    }
    if id == TOKEN_TYPE {
        retStr = "type";
    }
    if id == TOKEN_CONST {
        retStr = "const";
    }
    if id == TOKEN_VAR {
        retStr = "var";
    }
    if id == TOKEN_STRUCT {
        retStr = "struct";
    }
    if id == TOKEN_RETURN {
        retStr = "return";
    }
    if id == TOKEN_FUNC {
        retStr = "func";
    }
    if id == TOKEN_PACKAGE {
        retStr = "package";
    }
    if id == TOKEN_IMPORT {
        retStr = "import";
    }

    return retStr;
}
