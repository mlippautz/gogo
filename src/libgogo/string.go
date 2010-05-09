// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo string functions
//

package libgogo

//
// Returns the length of an ASCII string
// Implemented in assembler (see corresponding .s file)
//
func StringLength(str string) uint64;

//
// Returns the length of an ASCII string referred to by the string pointer given
// Implemented in assembler (see corresponding .s file)
//
func StringLength2(str *string) uint64;

//
// Returns the internal address of the character sequence referred to by the string pointer given
// Implemented in assembler (see corresponding .s file)
//
func GetStringAddress(str *string) uint64;

//
// Compares to strings and returns 0 for equality and 1 otherwise
//
func StringCompare(str1 string, str2 string) uint64 {
    var i uint64;
    var equal uint64 = 0; //Assume equality by default
    var strlen1 uint64 = StringLength(str1);
    var strlen2 uint64 = StringLength(str2);
    if strlen1 != strlen2 { //Return inequality if lengths are not equal
       equal = 1;
    } else {
        for i = 0; i < strlen1; i = i + 1 { //If lengths are equal compare every character
            if str1[i] != str2[i] { //If two characters differ => inequality
                equal = 1;
            }
        }
    }
    return equal;
}

//
// Sets a string's internal address and length; the string given by the specified string pointer referring to it
// Implemented in assembler (see corresponding .s file)
//
func SetStringAddressAndLength(str *string, new_addr uint64, new_length uint64);

//
// Appends a single character to a string.
// This functions creates a new string by copying the old one and then appending the character
//
func CharAppend(str *string, char byte) {
    var nullByte byte = 0; //End of string constant
    var strlen uint64 = StringLength2(str); //Get length of old string
    var new_length uint64 = strlen + 1; //The length of the new string the length of the old string plus the length of the character to be appended (1)
    var new_addr uint64 = Alloc(new_length + 1); //Allocate memory for the new string, including the space for the end of string constant
    var old_addr uint64 = GetStringAddress(str);
    CopyMem(old_addr, new_addr, strlen); //Copy the content of the old string to the newly allocated memory
    CopyMem(ToUint64FromBytePtr(&char), new_addr + strlen, 1); //Copy the additional character to its according position at the end of the new string
    CopyMem(ToUint64FromBytePtr(&nullByte), new_addr+strlen + 1, 1); //Append the end of string constant
    SetStringAddressAndLength(str, new_addr, new_length); //Update the string in order to point to the new character sequence with the new length
}

//
// Concatenates two strings
// This function creates a new string by copying the first one and appending the second one
//
func StringAppend(str *string, append_str string) {
    var nullByte byte = 0; //End of string constant
    var strlen uint64 = StringLength2(str); //Get length of first string
    var strappendlen uint64 = StringLength(append_str); //Get length of second string
    var new_length uint64 = strlen + strappendlen; //The length of the new string is the length of the first plus the length of the second string
    var new_addr uint64 = Alloc(new_length + 1); //Allocate memory for the new string, including the space for the end of string constant
    var old_addr uint64 = GetStringAddress(str);
    var append_addr uint64 = GetStringAddress(&append_str);
    CopyMem(old_addr, new_addr, strlen); //Copy the content of the first string to the newly allocated memory
    CopyMem(append_addr, new_addr + strlen, strappendlen); //Copy the second string to its according position at the end of the new string
    CopyMem(ToUint64FromBytePtr(&nullByte), new_addr + strlen + strappendlen + 1, 1); //Append the end of string constant
    SetStringAddressAndLength(str, new_addr, new_length); //Update the string in order to point to the new character sequence with the new length
}