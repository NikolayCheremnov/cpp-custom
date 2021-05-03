// first good test

const short int a = 1;

void main() {
    /*write code here
        ...
    */
    // assigments
    int variable = 20;
    short int variable2 = 30;
    long int variable3 = 404    // invalid lexeme 'const', expected ';'
    const int c = 1;
    int f = 1;

    // proc call
    proc(variable, variable2);

    // expressions
    variable = c - variable2;
    variable3 = (variable / variable2) % c - 10;
    // cycle
    for(int i = 1; i < 10; i = i + 1) {
        f = f * i;
    }
}

// procedure
void proc(int a, b) {
    int c = a + b;
}

//final comment without \n