# Syscall Spec

### I/O Flush
#### Call Number 2
stdout 메모리 영역을 읽어 메모리에 있는 내용을 I/O에 write하고 flush한다.

### String Length
#### Call Number 3
Register 0에 담긴 string value의 길이를 return한다.

### String Substring
#### Call Number 4
Reigster 0에 담긴 string에 대해 Register 1과 Register 2로 정해진 범위에 대해 잘라낸 뒤 return한다.

### String Matching
#### Call Number 5
Register 0에 담긴 string에 대해 Register 1으로 주어진 값을 찾아 위치를 return한다.

### String Replace
#### Call Number 6
Register 0에 담긴 string에 대해 Register 1으로 주어진 값을 찾아 전부 Register 2의 값으로 교체한 후 return한다.

### String Regexp
##### Call Number 7
Register 0에 담긴 string에 대해 Register 1으로 주어진 Regexp 식을 적용해 return한다.