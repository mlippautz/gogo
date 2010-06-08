package main

type foo struct {
    bar1 uint64;
    bar2 string;
    bar3 *uint64;
};

type test struct {
    bar foo;
    bazz [20]string;
    barPtr *foo;
};

var initTest byte = 1 + 2 * 3;
var initTest2 uint64 = initTest + 1;
var fooArrInst [10]foo;
var testInst test;
var testPtr *test;

func test(bla uint64, blub string) {
    var a uint64 = 5;
    var b uint64 = 1;
    var c string = "Test";
    var d *string;
    var e byte;
    var f byte;
    b = e + f + 1 - a;
    a = 1 + b * a;
    c = fooArrInst[a].bar2;
    c[1] = testPtr.barPtr.bar2[1];
    d[2] = 'x';
    a = fooArrInst[b].bar1 + testInst.barPtr.bar1;
    fooArrInst[0].bar3 = &a;
    fooArrInst[0].bar3 = testPtr.barPtr.bar3;
    fooArrInst[a].bar3 = &a;
    fooArrInst[a].bar3 = testPtr.barPtr.bar3;

    if (a < 10) && (b > 10) || (a == 5) {
        if (b > 15) {
            b = 4;
        } else {
            b = 3;
        }
    } else {
        b = 2;
    }
}

func foo() uint64 {
    return 3;
}

func bar(a uint64) uint64 {
    return a;
}

func foobar(a uint64, b uint64, c uint64) uint64 {
    return a + b + c;
}

func bazz(a uint64) uint64 {
    var ret uint64;
    var tmp uint64 = 3;
    foo();
    tmp = bar(tmp);
    ret = foobar(a, tmp, 3);
    return ret;
}

func blub(r string) string {
    var x string = r;
    return x;
}

func main() {
    var ret uint64;
    var u string = "ABC";
    var x string;
    ret = bazz(1);
    x = blub(u);
    ret = bazz_fwd(1);
}

func foo_fwd() {
    var ret uint64;
    ret = bazz_fwd(1);
}

func bazz_fwd(b byte) uint64 {
    return b;
}

func muh() {
    var x string;
    x = maeh(x);
}

func maeh(y string) string {
}
