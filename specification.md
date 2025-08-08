# CUTTER : Compact Universial Tokenized Text Excution Rules

## Object
Object는 Cutter의 기본 객체 단위이다. Object는 값을 띌 수 있으며, 호출 가능한 형태다. <br>
Object는 Evaluation Value를 가지며, 모든 Object는 Evaluation Value를 반환하여야 한다. <br> <br>

예를 들어, 단순히 `5`를 저장하는 Object `foo`는 `5`를 Evalutation Value로써 반환한다.
```
@define(foo 5)

@foo
> 5
```

Object는 호출 시 인자를 가질 수 있다.
```
@define(foo fizz buzz sum(fizz ` ` buzz))

@foo(`Hello World`, `from Cutter!`)
> Hello World from Cutter!
``` 

## Noraml Text
Cutter는 모든 텍스트가 일반 출력을 통해 출력된다. @를 접두사로 호출된 Object의 Evaluation Value는 모두 최종적으로 텍스트로 치환되어 출력된다.

## Atom Value
Object가 아닌 제일 기본적인 단위의 Value이다. 모든 Evaluation Value가 해당 형태 중 하나를 도출하여야 한다. 다음과 같은 Value를 가질 수 있다.
```
`John Doe` - str
42         - int
3.141592   - real
!t/!f      - bool
```

## Anonymus Object

익명 Object의 경우 ()로 둘러싸인 모든 내용은 Evaluation Value로 간주된다. 익명 Object로 간주되기 위해서는 몇가지 조건이 존재한다.

- Object의 호출의 인자로 사용될 때에만 인정받을 수 있다.
- Evaluation Value가 무조건 존재해야만 한다.