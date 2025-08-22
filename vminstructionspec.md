# VM Instruction Spec

### OpDefFunc
함수를 정의한다. Oprand1에 함수의 이름을 전달한다.

### OpCall
함수를 호출한다. Oprand1에 호출할 함수의 이름을 전달한다.

### OpReturn
함수 호출이 끝났음을 알린다.

### OpRegSet
레지스터에 값을 쓴다. Oprand1에 레지스터 번호를, Oprand2에 값을 전달한다.

### OpRegMov
레지스터의 값을 다른 레지스터로 옮긴다. Oprand1에 대상 레지스터를, Oprand2에 옮길 레지스터를 전달한다.

### OpMemSet
메모리에 값을 쓴다. Oprand1에 메모리 영역의 이름을, Oprand2에 값을 전달한다.

### OpMemMov
메모리의 값을 다른 메모리 영역으로 옮긴다. Oprand1에 대상 메모리 영역을, Oprand2에 옮길 메모리 영역을 전달한다.

### OpRslSet
결과 레지스터에 값을 쓴다. Oprand1에 값을 전달한다.

### OpRslMov
결과 레지스터의 값을 다른 레지스터로 옮긴다. Oprand1에 대상 레지스터를 전달한다.

### OpLdr
메모리에서 값을 읽어 레지스터에 쓴다. Oprand1에 레지스터 번호를, Oprand2에 메모리 영역의 이름을 전달한다.

### OpStr
레지스터의 값을 메모리에 쓴다. Oprand1에 메모리 영역의 이름을, Oprand2에 레지스터 번호를 전달한다.

### OpRslStr
결과 레지스터의 값을 메모리에 쓴다. Oprand1에 메모리 영역의 이름을 전달한다.

### OpStrReg
레지스터에 지정된 메모리 영역에 레지스터에 지정된 메모리 영역의 값을 쓴다. Oprand1에 대상 메모리 영역이 지정된 레지스터를, Oprand2에 값이 지정된 레지스터를 전달한다.

### OpSyscall
시스템 콜을 호출한다. Oprand1에 시스템 콜 번호를 전달한다.

### OpAdd / OpSub / OpMul / OpDiv / OpMod
사칙연산을 수행한다. Oprand1과 Oprand2의 값을 연산한 뒤 Oprand3에 지정된 레지스터에 값을 쓴다.

### OpAnd / OpOr / OpNot
논리연산을 수행한다.

### OpCmpEq / OpCmpNeq
두 레지스터의 값이 같은지 다른지 판별한 뒤 Oprand3에 지정된 레지스터에 값을 쓴다.

### OpBrch
Oprand1의 값이 참일 경우 Oprand2의 값을, 거짓일 경우 Oprand3의 값을 결과 레지스터에 쓴다.

### OpClearReg
모든 레지스터의 값을 지운다.

### OpHlt
VM을 정지시킨다.