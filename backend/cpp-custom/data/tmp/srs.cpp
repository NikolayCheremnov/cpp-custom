int res = 0;

void add(bool d) {
res = res + d;
add(2);
}

void main() {
add(5);
}