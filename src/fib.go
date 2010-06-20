// fib test

package main

func main() uint64 {
    var retValue uint64;
    var tofib1 uint64;
    var tofib2 uint64;
    var tmpStr string;
    tofib1 = fib1(11);
    tofib2 = fib2(11);
    libgogo.PrintString("Fibonacci\n");
    libgogo.PrintString("=========\n");
    libgogo.PrintString("Rekursiv(11): ");
    tmpStr = libgogo.IntToString(tofib1);
    libgogo.PrintString(tmpStr);
    libgogo.PrintString("\n");
    libgogo.PrintString("Iterativ(11): ");
    tmpStr = libgogo.IntToString(tofib2);
    libgogo.PrintString(tmpStr);
    libgogo.PrintString("\n");
}

func fib1(i uint64) uint64 {
    var retVal uint64;
    var fib1 uint64;
    var fib2 uint64;
    if i <= 2 {
        retVal = 1;
    } else {
        fib1 = fib1(i-1);
        fib2 = fib1(i-2);
        retVal = fib1 + fib2;
    }
    return retVal;
}

func fib2(n uint64) uint64 {
    var a uint64;
    var b uint64;
    var new uint64;
    var cnt uint64;
    var retValue uint64;

    if n<=2 {
        retValue = 1;
    } else {
        a = 1;
        b = 1;
        for cnt=3;cnt<=n;cnt=cnt+1 {
            new = a+b;
            a = b;
            b = new;
        }
        retValue = new;
    }
    return retValue;
}
