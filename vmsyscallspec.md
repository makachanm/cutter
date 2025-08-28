# Syscall Spec

### Set Function Return
#### Call Number 1
Register 0에 담긴 함수 이름과 Register 1에 담긴 반환 값으로 함수의 반환값을 설정합니다.

### I/O Flush
#### Call Number 2
stdout 메모리 영역을 읽어 메모리에 있는 내용을 I/O에 write하고 flush합니다.

### String Length
#### Call Number 3
Register 0에 담긴 string value의 길이를 반환합니다.

### String Substring
#### Call Number 4
Reigster 0에 담긴 string에 대해 Register 1과 Register 2로 정해진 범위에 대해 잘라낸 뒤 반환합니다.

### String Matching
#### Call Number 5
Register 0에 담긴 string에 대해 Register 1으로 주어진 값을 찾아 위치를 반환합니다.

### String Replace
#### Call Number 6
Register 0에 담긴 string에 대해 Register 1으로 주어진 값을 찾아 전부 Register 2의 값으로 교체한 후 반환합니다.

### String Regexp
#### Call Number 7
Register 0에 담긴 string에 대해 Register 1으로 주어진 Regexp 식을 적용해 찾은 모든 값들을 공백으로 구분된 문자열로 반환합니다.

### Array Make
#### Call Number 8
Register 0에 담긴 문자열을 이름으로 하는 새로운 배열을 생성합니다.

### Array Push
#### Call Number 9
Register 0에 담긴 이름의 배열에 Register 1의 값을 추가합니다.

### Array Set
#### Call Number 10
Register 0에 담긴 이름의 배열의 Register 1 위치에 Register 2의 값을 설정합니다.

### Array Get
#### Call Number 11
Register 0에 담긴 이름의 배열의 Register 1 위치의 값을 가져와 반환합니다.

### Array Length
#### Call Number 13
Register 0에 담긴 이름의 배열의 길이를 반환합니다.

### Get Environment Variable
#### Call Number 14
Register 0에 담긴 이름의 환경변수의 값을 반환합니다.

### Excute Command
#### Call Number 15
Register 0에 담긴 명령어를 실행합니다.

### Get Operating System Kernel Types
#### Call Number 16
현재 실행중인 OS의 커널 타입을 반환합니다. linux는 1, bsd는 2, darwin은 3, windows는 4를, 그 외엔 5를 반환합니다.