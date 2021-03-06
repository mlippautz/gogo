// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package libgogo

type Item struct {
    Mode uint64;
    Itemtype *TypeDesc;
    PtrType uint64; //If 1, item is of type *Itemtype, otherwise it is of type Itemtype
    A uint64; //Mode = MODE_VAR: variable address; Mode = MODE_CONST: const value; Mode = MODE_REG: 0 = register contains value, 1 = register contains address
    R uint64; //Mode = MODE_REGISTER: register number
    C uint64; // Mode = MODE_COND: Type of comparison; Mode = MODE_REG: Optional second register occupied
    Global uint64; //If 1, the variable is global, 0 if it is function-local, 2 if it is a function parameter
    LinkerInformation string; //Mode = MODE_VAR and Global = 0: If non-empty, indicates unknown parameter offset to be replaced by the linker
};

//
// Pseudo constants that specify the descriptor sizes 
//
var ITEM_SIZE uint64 = 72; //9*8 bytes space for an object

//
// Modes for items
//
var MODE_VAR uint64 = 1;
var MODE_CONST uint64 = 2;
var MODE_REG uint64 = 3;
var MODE_COND uint64 = 4;
var MODE_UNDEF uint64 = 5;

//
// Relation types
//
var REL_EQ uint64 = 0;
var REL_NEQ uint64 = 1;
var REL_LT uint64 = 2;
var REL_GTE uint64 = 3;
var REL_LTE uint64 = 4;
var REL_GT uint64 = 5; 

//
// Convert the uint64 value (returned from malloc) to a real item address
//
func Uint64ToItemPtr(adr uint64) *Item;

//
//
//
func ItemPtrToUint64(item *Item) uint64;

//
// Creates a new, uninitialized item
//
func NewItem() *Item {
    var adr uint64;
    var item *Item;
    adr = Alloc(ITEM_SIZE);
    item = Uint64ToItemPtr(adr);
    item.Mode = MODE_UNDEF;
    return item;
}

//
// Sets the given item's properties
//
func SetItem(item *Item, mode uint64, itemtype *TypeDesc, ptrtype uint64, a uint64, r uint64, global uint64) {
    item.Mode = mode;
    item.Itemtype = itemtype;
    item.PtrType = ptrtype;
    item.A = a;
    item.C = 0;
    item.R = r;
    item.Global = global;
    item.LinkerInformation = "";
}
