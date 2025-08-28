# VM Instruction Spec

### OpDefFunc
함수를 정의합니다. Oprand1에 함수의 이름을 전달합니다.

### OpCall
함수를 호출합니다. Oprand1에 호출할 함수의 이름을 전달합니다.

### OpReturn
함수 호출이 끝났음을 알립니다.

### OpRegSet
레지스터에 값을 씁니다. Oprand1에 레지스터 번호를, Oprand2에 값을 전달합니다.

### OpRegMov
레지스터의 값을 다른 레지스터로 옮깁니다. Oprand1에 대상 레지스터를, Oprand2에 옮길 레지스터를 전달합니다.

### OpMemSet
메모리에 값을 씁니다. Oprand1에 메모리 영역의 이름을, Oprand2에 값을 전달합니다.

### OpMemMov
메모리의 값을 다른 메모리 영역으로 옮깁니다. Oprand1에 대상 메모리 영역을, Oprand2에 옮길 메모리 영역을 전달합니다.

### OpRslSet
결과 레지스터에 값을 씁니다. Oprand1에 값을 전달합니다.

### OpRslMov
결과 레지스터의 값을 다른 레지스터로 옮깁니다. Oprand1에 대상 레지스터를 전달합니다.

### OpLdr
메모리에서 값을 읽어 레지스터에 씁니다. Oprand1에 레지스터 번호를, Oprand2에 메모리 영역의 이름을 전달합니다.

### OpStr
레지스터의 값을 메모리에 씁니다. Oprand1에 메모리 영역의 이름을, Oprand2에 레지스터 번호를 전달합니다.

### OpRslStr
결과 레지스터의 값을 메모리에 씁니다. Oprand1에 메모리 영역의 이름을 전달합니다.

### OpStrReg
Oprand1 레지스터에 지정된 이름의 메모리 영역에, Oprand2 레지스터에 담긴 값을 씁니다.

### OpSyscall
시스템 콜을 호출합니다. Oprand1에 시스템 콜 번호를 전달합니다.

### OpAdd / OpSub / OpMul / OpDiv / OpMod
사칙연산을 수행합니다. Oprand1과 Oprand2의 값을 연산한 뒤 Oprand3에 지정된 레지스터에 값을 씁니다.

### OpAnd / OpOr / OpNot
논리연산을 수행합니다.

### OpCmpEq / OpCmpNeq
두 레지스터의 값이 같은지 다른지 판별한 뒤 Oprand3에 지정된 레지스터에 값을 씁니다.

### OpBrch
Oprand1의 값이 참일 경우 Oprand2의 값을, 거짓일 경우 Oprand3의 값을 결과 레지스터에 씁니다.

### OpClearReg
모든 레지스터의 값을 지웁니다.

### OpHlt
VM을 정지시킵니다.

### OpJmp
Oprand1에 지정된 주소로 점프합니다.

### OpJmpIfFalse
Oprand1 레지스터의 값이 거짓일 경우 Oprand2에 지정된 주소로 점프합니다.

### OpCstInt / OpCstReal / OpCstStr
Oprand1 레지스터의 값을 정수/실수/문자열로 변환한 뒤 결과 레지스터에 값을 씁니다.