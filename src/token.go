// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

//
// Set of recognized tokens 
//
const TOKEN_IDENTIFIER uint64 = 1;  // Identifier
const TOKEN_STRING uint64 = 2;      // String using "..."
const TOKEN_EOS uint64 = 3;         // End of Scan
const TOKEN_LBRAC uint64 = 4;       // Left bracket '('
const TOKEN_RBRAC uint64 = 5;       // Right bracket ')'
const TOKEN_LSBRAC uint64 = 6;      // Left square bracket '['
const TOKEN_RSBRAC uint64 = 7;      // Right square bracket ']'
const TOKEN_INTEGER uint64 = 8;     // Integer number
const TOKEN_LCBRAC uint64 = 9;      // Left curly bracket '{'
const TOKEN_RCBRAC uint64 = 10;     // Right curly bracket '}'
const TOKEN_PT uint64 = 11;         // Point '.'
const TOKEN_NOT uint64 = 12;        // Boolean negation '!'
const TOKEN_NOTEQUAL uint64 = 13;   // Comparison, not equal '!='
const TOKEN_SEMICOLON uint64 = 14;  // Semi-colon ';'
const TOKEN_COLON uint64 = 15;      // Colon ','
const TOKEN_ASSIGN uint64 = 16;     // Assignment '='
const TOKEN_EQUALS uint64 = 17;     // Equal comparison '=='
const TOKEN_CHAR uint64 = 18;       // Single Quoted Character 'x'
const TOKEN_REL_AND uint64 = 19;    // AND Relation '&&'
const TOKEN_REL_OR uint64 = 20;     // OR Relation '||'
const TOKEN_REL_GTOE uint64 = 21;   // Greather-Than or Equal '>='
const TOKEN_REL_GT uint64 = 22;     // Greather-Than '>'
const TOKEN_REL_LTOE uint64 = 23;   // Less-Than or Equal '<='
const TOKEN_REL_LT uint64 = 24;     // Less-Than '<'
const TOKEN_ARITH_PLUS uint64 = 25; // Arith. Plus '+'
const TOKEN_ARITH_MINUS uint64 = 26;// Arith. Minus '-'
const TOKEN_ARITH_MUL uint64 = 27;  // Arith. Multiplication '*'
const TOKEN_ARITH_DIV uint64 = 28;  // Arith. Division '/'
const TOKEN_OP_ADR uint64 = 29;     // Address operator '&'

//
// Advanced tokens, that are generated in the 2nd step from identifiers
// The tokens represent the corresponding language keywords.
//
const TOKEN_FOR uint64 = 101;
const TOKEN_IF uint64 = 102;
const TOKEN_TYPE uint64 = 103;
const TOKEN_CONST uint64 = 104;
const TOKEN_VAR uint64 = 105;
const TOKEN_STRUCT uint64 = 106;
const TOKEN_RETURN uint64 = 107;
const TOKEN_FUNC uint64 = 108;
const TOKEN_PACKAGE uint64 = 109;
const TOKEN_IMPORT uint64 = 110;
