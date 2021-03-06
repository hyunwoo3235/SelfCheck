# SelfCheck
학생 자가진단을 위한 기능들을 제공하는 서버

### 기능

* 학교 코드 검색
* 자가진단 현황 확인
* 자가진단 제출
* 기타 등등

##TODO

* 더 깔끔하게 코드 갈아엎기
* 일정 시간에 자동 자가진단 추가
* 교육부 사이트 수정


```
SelfCheck
│   README.md
│   main.go    메인 파일
│   Dockerfile    도커 빌드용
│   database.db    학교 정보 데이터베이스
│
└───assets
│   │   favicon.svg
│   
└───database
    │   database.go
│   
└───eduro
    │   selfcheck.go
    │   
```