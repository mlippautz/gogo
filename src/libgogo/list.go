// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo stack functions
//

package libgogo

//
// List data structure
//
type List struct {
    baseAddress uint64; //Where the list starts
    itemCount uint64; //How many items there currently are in the list
    capacity uint64; //How many items there can be max. (not to be changed from outside the library)
};

//
// Initializes the given list
//
func InitializeList(uninitializedList *List) {
    uninitializedList.baseAddress = Alloc(16 * 8); //Allocate 16 items by default
    uninitializedList.capacity = 16;
    uninitializedList.itemCount = 0; //Reset item count (to zero)
}

//
// Adds an item to the given list, increasing its capacity if required
//
func AddItem(list *List, value uint64) {
    var newAddress uint64;
    if (list.capacity == list.itemCount) { //Grow list if its capacity doesn't suffice to add another item
        newAddress = Alloc(list.capacity * 2 * 8); //Double the capacity
        CopyMem(list.baseAddress, newAddress, list.capacity * 8); //Copy old list items
        list.baseAddress = newAddress; //Set new address as base address
        list.capacity = list.capacity * 2; //Update (increase) capacity
    }
    CopyMem(ToUint64FromUint64Ptr(&value), list.baseAddress + 8 * list.itemCount, 8); //Append the new value by copying its value into the memory of the corresponding list item
    list.itemCount = list.itemCount + 1; //Update item count
}

//
// Removes the i-th item of the given list and returns it
// Note that this function call fails if there are less than i items in the list or if the list is empty
//
func RemoveItemAt(list *List, index uint64) uint64 {
    var returnValue uint64 = GetItemAt(list, index); //Get the correspondig value value from the list
    var i uint64;
    for i = index; i < list.itemCount - 1; i = i + 1 { //Remove item by moving the following ones "backwards" in order to fill the gap caused by the deleted item
        CopyMem(list.baseAddress + (i + 1) * 8, list.baseAddress + 8 * i, 8); //Move item at position i + 1 to position i
    }
    list.itemCount = list.itemCount - 1; //Update (decrease) item count
    return returnValue; //Return value removed from list
}

//
// Returns the i-th item of a given list without removing it
// Note that this function call fails if there are less than i items in the list or if the list is empty
//
func GetItemAt(list *List, index uint64) uint64 {
    var returnValue uint64;
    if (list.itemCount <= index) { //Check if there are enough items in the list
        PrintString("Tried to GetItemAt(");
        PrintNumber(index);
        PrintString(") from a list with only ");
        PrintNumber(list.itemCount);
        ExitError(" items - out of bounds error", 125);
    }
    CopyMem(list.baseAddress + index * 8, ToUint64FromUint64Ptr(&returnValue), 8); //Copy the value on position i of the list (e.g. its corresponding memory) into the return variable
    return returnValue;
}

//
// Returns the number of items in the given list
//
func GetListItemCount(list *List) uint64 {
    return list.itemCount;
}

//
// Returns the current capacity of the given list
//
func GetListCapacity(list *List) uint64 {
    return list.capacity;
}
