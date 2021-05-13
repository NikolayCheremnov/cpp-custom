int res = 0;

void add(bool d) {
res = res + d;
return;
add(2);
}

void main() {
add(5);
}