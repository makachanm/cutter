# Standard Functions

### add / sub / mul / div / mod
Int / Real / Str에 한해 사칙연산을 수행한다.

### set
Object에 선언된 내용을 바꿀 때 사용한다. 새롭게 정의된 값이 기존과 Type이 다를 경우 Object의 Type이 이를 반영해 변화한다.

### make[type]arr / getarr
주어진 Type에 대해 Array를 만들고, 접근한다.
make[type]arr의 첫 인수는 Array의 이름을 입력받는다.
getarr의 첫 인수는 Array의 Index에 해당하는 값이 대입될 Object의 이름을 전달받고, 두번째 인수는 해당하는 Index를 전달받는다.