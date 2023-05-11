void sumNumbers(int n) {
    int sum = 0;
    for (int i = 0; i < n; i = i + 1) {
        sum = sum + i;
    }
}

void squareSumNumbers(int n) {
    int sum = 0;
    for (int i = 0; i < n; i = i + 1) {
        for (int j = 0; j < n; j = j + 1) {
            sum = sum + j;
        }
    }
}

void factorial(int n) {
    long int res = 1;
    for (int i = 1; i <= n; i = i + 1) {
        res = res * i;
    }
}

void testOperation(int a, long int b, short int c, bool d) {
    int x = a + b * c;
    bool y = b + c * d;
    const long int foo = x + y;
    for (x = 0; x < 2 + 2 * 2; x = x * 2) {
        for (y = 0; y < x; y = y + a) {
            a = a - 1;
            b = b - 1;
            c = c - 1;
            d = d - 1;
        }
    }
    int tmp = x;
    x = y;
    y = tmp;
}

void average(int a, int b) {
    int res = (a + b) / 2;
}

void recursion(int n) {
    for (int i = 0; i < n; i = i + 1) {
        int val = n - 1;
        recursion(val);
    }
}

void main() {
    // 1. test sum numbers
    const int count = 10;
    sumNumbers(count);

    // 2. test square sum numbers
    int countSquare = count * count;
    squareSumNumbers(countSquare);

    // 3. test factorial
    factorial(count);

    // 4. test some operations
    testOperation(count, countSquare, 10, 15);

    // 5. and average
    average(10, 15);

    // 6. recursion test
    recursion(3);

    return;
}
