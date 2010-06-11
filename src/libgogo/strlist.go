// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo list functions
//

package libgogo

//
// String list data structure
//
type StringList struct {
    baseAddress uint64; //Where the list starts
    itemCount uint64; //How many items there currently are in the list
    capacity uint64; //How many items there can be max. (not to be changed from outside the library)
};

//
// Initializes the given string list
//
func InitializeStringList(uninitializedList *StringList) {
    uninitializedList.baseAddress = Alloc(16 * 2); //Allocate 2 items by default
    uninitializedList.capacity = 2;
    uninitializedList.itemCount = 0; //Reset item count (to zero)
}

//
// Adds an item to the given string list, increasing its capacity if required
//
func AddStringItem(list *StringList, value string) {
    var newAddress uint64;
    if (list.capacity == list.itemCount) { //Grow list if its capacity doesn't suffice to add another item
        newAddress = Alloc(list.capacity * 2 * 16); //Double the capacity
        CopyMem(list.baseAddress, newAddress, list.capacity * 16); //Copy old list items
        list.baseAddress = newAddress; //Set new address as base address
        list.capacity = list.capacity * 2; //Update (increase) capacity
    }
    newAddress = ToUint64FromStringPtr(&value);
    CopyMem(newAddress, list.baseAddress + 16 * list.itemCount, 16); //Append the new value by copying its value into the memory of the corresponding list item
    list.itemCount = list.itemCount + 1; //Update item count
}

//
// Removes the i-th item of the given string list and returns it
// Note that this function call fails if there are less than i items in the list or if the list is empty
//
func RemoveStringItemAt(list *StringList, index uint64) string {
    var returnValue string;
    var i uint64;
    returnValue = GetStringItemAt(list, index); //Get the correspondig value value from the list
    for i = index; i < list.itemCount - 1; i = i + 1 { //Remove item by moving the following ones "backwards" in order to fill the gap caused by the deleted item
        CopyMem(list.baseAddress + (i + 1) * 16, list.baseAddress + 16 * i, 16); //Move item at position i + 1 to position i
    }
    list.itemCount = list.itemCount - 1; //Update (decrease) item count
    return returnValue; //Return value removed from list
}

//
// Returns the i-th item of a given string list without removing it
// Note that this function call fails if there are less than i items in the list or if the list is empty
//
func GetStringItemAt(list *StringList, index uint64) string {
    var returnValue string;
    var temp uint64;
    if (list.itemCount <= index) { //Check if there are enough items in the list
        PrintString("Tried to GetItemAt(");
        PrintNumber(index);
        PrintString(") from a list with only ");
        PrintNumber(list.itemCount);
        ExitError(" items - out of bounds error", 125);
    }
    temp = ToUint64FromStringPtr(&returnValue);
    CopyMem(list.baseAddress + index * 16, temp, 16); //Copy the value on position i of the list (e.g. its corresponding memory) into the return variable
    return returnValue;
}

//
// Returns the number of items in the given list
//
func GetStringListItemCount(list *StringList) uint64 {
    return list.itemCount;
}

//
// Returns the current capacity of the given list
//
func GetStringListCapacity(list *StringList) uint64 {
    return list.capacity;
}
