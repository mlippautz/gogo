package test


var bump_ptr uint64;
var end_ptr uint64;

func test(bla uint64, blub string) {
    var a uint64;
    var b uint64;
    var c uint64;
    var d uint64;
    var e uint64;
    var f uint64;



    if (a != b) && (c < d) || ((a > b) || (c <d) && (e > f)) || (a > b)  {
        a = 1;
        if (c > 4) {
            c = 3;
        }
    } else {
        a = 5;
    }

}
