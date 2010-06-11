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
// Allocates a new, empty byte array for a string
//
func ResetString(str *string) {
    var nullByte byte = 0; //End of string constant
    var nullByte_addr uint64;
    var n uint64;
    n = Alloc(1);
    nullByte_addr = ToUint64FromBytePtr(&nullByte);
    CopyMem(n, nullByte_addr, 1);
    SetStringAddressAndLength(str, n, 0); //New string with length 0
}

//
// Compares to strings and returns 0 for equality and 1 otherwise
//
func StringCompare(str1 string, str2 string) uint64 {
    var i uint64;
    var equal uint64 = 0; //Assume equality by default
    var strlen1 uint64;
    var strlen2 uint64;
    strlen1 = StringLength(str1);
    strlen2 = StringLength(str2);
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
// Rounds length to the next power of two
//
func GetMaxStringLength(length uint64) uint64 {
    var i uint64;
    var j uint64;
    for i = 0; length != 0; i = i + 1 { //Calculate number of divisions by 2 required to reach 0 => log2(length), rounded towards +inf
        length = length / 2;
    }
    length = 1; //Restart with length of 1
    for j = 0; j < i; j = j + 1 { //Multiply by 2 as many times as divided previously => 2^i
        length = length * 2;
    }
    if length == 1 { //Special case length = 1 => still return 1
        length = 2;
    }
    return length - 1;
}

//
// Appends a single character to a string.
// This functions creates a new string by copying the old one and then appending the character
//
func CharAppend(str *string, char byte) {
    var nullByte byte = 0; //End of string constant
    var strlen uint64;
    var new_length uint64;
    var max_length uint64;
    var new_addr uint64;
    var old_addr uint64;
    var char_addr uint64;
    var nullByte_addr uint64;
    strlen = StringLength2(str); //Get length of old string
    new_length = strlen + 1; //The length of the new string the length of the old string plus the length of the character to be appended (1)
    max_length = GetMaxStringLength(strlen + 1); //Get maximum capacity of old string, considering its trailing '\0'
    old_addr = GetStringAddress(str);
    if (old_addr < start_ptr) || ((new_length + 1) > max_length) { //If using non-managed address or the old string's capacity doesn't suffice => allocate more memory
        max_length = GetMaxStringLength(new_length + 1); //Consider trailing '\0'
        new_addr = Alloc(max_length); //Allocate memory for the new string, including the space for the end of string constant
        if strlen > 0 { //Only copy old content if there is content to copy
            CopyMem(old_addr, new_addr, strlen); //Copy the content of the old string to the newly allocated memory
        }
    } else { //Re-use old address as capacity suffices; no copy operation necessary
        new_addr = old_addr;
    }
    char_addr = ToUint64FromBytePtr(&char);
    CopyMem(char_addr, new_addr + strlen, 1); //Copy the additional character to its according position at the end of the new string
    nullByte_addr = ToUint64FromBytePtr(&nullByte);
    CopyMem(nullByte_addr, new_addr+strlen + 1, 1); //Append the end of string constant
    SetStringAddressAndLength(str, new_addr, new_length); //Update the string in order to point to the new character sequence with the new length
}

//
// Concatenates two strings
// This function creates a new string by copying the first one and appending the second one
//
func StringAppend(str *string, append_str string) {
    var nullByte byte = 0; //End of string constant
    var strlen uint64;
    var strappendlen uint64;
    var new_length uint64;
    var max_length uint64;
    var new_addr uint64;
    var old_addr uint64;
    var append_addr uint64;
    var nullByte_addr uint64;
    strlen = StringLength2(str); //Get length of first string
    max_length = GetMaxStringLength(strlen + 1); //Get maximum capacity of old string, considering its trailing '\0'
    strappendlen = StringLength(append_str); //Get length of second string
    new_length = strlen + strappendlen; //The length of the new string is the length of the first plus the length of the second string
    old_addr = GetStringAddress(str);
    if (old_addr < start_ptr) || ((new_length + 1) > max_length) { //If using non-managed address or the old string's capacity doesn't suffice => allocate more memory
        max_length = GetMaxStringLength(new_length + 1); //Consider trailing '\0'
        new_addr = Alloc(max_length); //Allocate memory for the new string, including the space for the end of string constant
        if strlen > 0 { //Only copy old content if there is content to copy
            CopyMem(old_addr, new_addr, strlen); //Copy the content of the old string to the newly allocated memory
        }
    } else { //Re-use old address as capacity suffices; no copy operation necessary
        new_addr = old_addr;
    }
    append_addr = GetStringAddress(&append_str);
    CopyMem(append_addr, new_addr + strlen, strappendlen); //Copy the second string to its according position at the end of the new string
    nullByte_addr = ToUint64FromBytePtr(&nullByte);
    CopyMem(nullByte_addr, new_addr + strlen + strappendlen + 1, 1); //Append the end of string constant
    SetStringAddressAndLength(str, new_addr, new_length); //Update the string in order to point to the new character sequence with the new length
}
